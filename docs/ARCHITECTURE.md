# Arsitektur

Bagaimana sistem disusun. Untuk alasan di balik pilihannya baca
[DECISIONS.md](DECISIONS.md); untuk cara mengerjakannya baca
[GUIDES.md](../apps/backend/docs/GUIDES.md).

---

## 1. Bentuk yang paling menentukan

Sebelum apa pun yang lain, pahami ini. Klontong Net adalah **marketplace**, dan
satu keranjang bisa berisi barang dari beberapa penjual.

```
Order  (milik pembeli, satu pembayaran, satu alamat)
 ├── Suborder  penjual A, outlet A1, kurir ekspedisi   → ongkir, status, payout sendiri
 ├── Suborder  penjual B, outlet B2, antar lokal       → ongkir, status, payout sendiri
 └── Suborder  penjual C, outlet C1, ambil di toko     → ongkir, status, payout sendiri
```

**Suborder adalah unit kerja sesungguhnya.** Pengiriman, pembatalan, komisi, dan
pencairan dana semuanya terjadi di level ini. Status Order induk adalah
**turunan** dari status anak-anaknya, tidak pernah disetel langsung.

Konsekuensi yang paling sering salah dipahami: satu penjual menolak pesanan
**tidak membatalkan** pesanan penjual lain. Yang terjadi adalah stok suborder itu
dilepas, refund parsial dikeluarkan, dan order induk tetap berjalan.

Setiap kali menulis kode yang berhubungan dengan pesanan, tanya dulu: ini level
order atau suborder? ADR-002.

---

## 2. Gambaran besar

```
                    ┌──────────────────────────────────────┐
  pembeli ─────────▶│  chi router + middleware global      │
  penjual ─────────▶│  (requestid, log, recover, cors)     │
  admin   ─────────▶└──────────────┬───────────────────────┘
  gateway ─────────▶               │
                                   │
   ┌──────────┬──────────┬─────────┼─────────┬──────────┬──────────┐
   ▼          ▼          ▼         ▼         ▼          ▼          ▼
identity  customer   seller    catalog  inventory  pricing   promotion
                        │         │         │         │          │
                        └────┬────┴────┬────┴────┬────┘          │
                             ▼         ▼         ▼               │
                            cart ────────────▶ order ◀───────────┘
                                                 │
                        ┌────────────────────────┼──────────────┐
                        ▼                        ▼              ▼
                     payment              fulfillment      settlement
                        │                        │              │
                        ▼                        ▼              ▼
                   gateway pembayaran      agregator kurir   penyedia transfer

                        outbox ──▶ event bus ──▶ audit · notification · review
```

`order` bergantung ke paling banyak modul, dan itu wajar — checkout memang titik
temu seluruh sistem. Yang harus dijaga: **modul lain tidak boleh bergantung balik
ke `order`** kecuali lewat event.

---

## 3. Daftar modul

| Modul | Milik | Bergantung ke |
|---|---|---|
| `identity` | akun, sesi, peran | — |
| `customer` | profil pembeli, alamat | `identity` |
| `seller` | penjual, outlet, pengelola, verifikasi | `identity` |
| `catalog` | produk, varian, kategori, satuan | `seller` |
| `inventory` | stok per outlet, reservasi, mutasi | `catalog` |
| `pricing` | harga, tier grosir, harga per outlet | `catalog` |
| `promotion` | voucher, diskon, kuota | `seller` |
| `cart` | keranjang, pengelompokan per penjual | `catalog`, `pricing`, `inventory`, `seller` |
| `order` | checkout, pesanan, suborder, siklus hidup | banyak |
| `payment` | gateway, COD, webhook, refund | — |
| `fulfillment` | tiga metode kirim, ongkir, tracking | `seller` |
| `settlement` | komisi, buku besar, pencairan | `seller` |
| `review` | ulasan produk dan penjual | `order`, `identity` |
| `notification` | inbox, email | konsumen event |
| `audit` | catatan aksi | konsumen event |

Dua cara komunikasi antar modul, dan hanya dua:

| | Kapan | Caranya |
|---|---|---|
| Sinkron | butuh jawaban sekarang untuk melanjutkan | `port.go` milik modul tujuan |
| Asinkron | efek samping, boleh terlambat | outbox → event bus |

---

## 4. Empat lapisan di dalam modul

```
modules/order/
├── module.go     ← pintu masuk publik
├── port.go       ← kontrak untuk modul lain
├── events.go     ← nama & payload event publik
└── internal/     ← TERKUNCI oleh compiler Go
    ├── domain/      entity, aturan, state machine, interface repository
    ├── app/         usecase, transaksi, izin
    ├── infra/       implementasi repository, klien penyedia luar
    └── transport/rest/
```

Arah dependensi: `transport → app → domain ← infra`.

`domain` tidak mengimpor `net/http`, `pgx`, atau `chi`. Ia mendefinisikan
interface repository; `infra` yang mengimplementasikannya.

Package di dalam `internal/` hanya bisa diimpor dari dalam direktori induknya,
jadi `modules/cart` **tidak bisa dikompilasi** kalau mencoba mengimpor
`modules/order/internal/domain`. Batas dijaga compiler, bukan disiplin tim.

---

## 5. Empat area HTTP

| Area | Pemanggil | Akses |
|---|---|---|
| `/api/v1/*` | pembeli | sebagian publik, sebagian butuh login |
| `/api/v1/seller/*` | dasbor penjual | login + peran seller + keanggotaan toko |
| `/api/v1/admin/*` | panel marketplace | peran admin |
| `/webhook/*` | gateway & kurir | verifikasi signature, tanpa sesi, tanpa CORS |

Area pembeli **tidak boleh** dipasangi middleware auth secara menyeluruh.
Halaman produk dan pencarian harus bisa dibuka tanpa login — itu penting untuk
SEO dan konversi. Pakai `OptionalAuth` di sana.

Peran hanyalah gerbang kasar. Bahwa seseorang berperan `seller` tidak berarti ia
berhak mengubah toko mana pun; kepemilikan diperiksa di usecase. Tanpa itu,
mengganti id di URL cukup untuk mengubah produk penjual lain.

---

## 6. Alur checkout

Alur paling kompleks di sistem. Dua tahap, dan tahap pertama tidak mengikat
apa pun.

```
POST /api/v1/checkout/preview          (boleh dipanggil berkali-kali)
  ├─ ambil keranjang
  ├─ kelompokkan per penjual → calon suborder
  ├─ pilih outlet per kelompok (stok cukup + terdekat)
  ├─ pricing.ResolveMany     BATCH
  ├─ inventory.AvailableMany BATCH
  ├─ fulfillment.Quote       per suborder
  ├─ promotion.Validate      → bagi diskon ke suborder (money.Distribute)
  └─ kembalikan rincian lengkap; TIDAK menahan stok

POST /api/v1/checkout                  (Idempotency-Key WAJIB)
  ├─ HITUNG ULANG semuanya dari nol
  ├─ bandingkan dengan angka yang ditampilkan klien
  │    berbeda → ErrPriceChanged beserta rinciannya
  └─ WithinTx:
       ├─ inventory.Reserve  (kunci baris, urutan konsisten)
       ├─ promotion.Redeem   (kuota dipotong di sini)
       ├─ buat Order + Suborder + Item, semuanya snapshot
       ├─ bekukan commission_bps tiap suborder
       ├─ cart.MarkConverted
       └─ outbox: EventOrderPlaced
     COMMIT
  └─ payment.CreatePayment (SETELAH commit) → instruksi bayar
```

Tiga hal yang tidak boleh dilanggar di sini:

**Tidak ada panggilan jaringan di dalam transaksi.** Memanggil gateway di dalam
transaksi berarti menahan koneksi database selama pihak lain berpikir, dan
rollback tidak membatalkan transaksi yang sudah terlanjur dibuat di sana.

**Kunci baris stok dalam urutan konsisten** (urutkan `outlet_id, variant_id`).
Dua pembeli yang checkout barang sama persis akan bertemu di sini, dan urutan
kunci yang berbeda antar transaksi adalah cara paling umum menciptakan deadlock.

**Klien mengirim pilihan, bukan angka.** Struct request checkout tidak punya
field harga atau total sama sekali. ADR-004.

---

## 7. Model stok

Stok disimpan **per outlet**, dalam **satuan terkecil**.

```
Varian "Indomie Goreng 1 dus"
  unit_code        = dus
  content_quantity = 40
  base_unit_code   = pcs

Pembeli memesan 2 dus  →  stok berkurang 80 pcs
```

Tanpa ini, stok tiga varian yang sebenarnya barang yang sama akan saling
berbohong. ADR-006.

Dua angka disimpan terpisah:

```
quantity_on_hand   fisik ada di rak
quantity_reserved  sudah di-checkout, belum diambil
tersedia = on_hand - reserved
```

Siklus hidupnya:

```
checkout       → reserved += qty, buat Reservation dengan expires_at
suborder kirim → on_hand -= qty, reserved -= qty, tulis Movement kind=sale
pesanan batal  → reserved -= qty, TIDAK menulis Movement (barang tak pernah keluar)
TTL habis      → worker melepas reservasi
```

Stok ditahan saat **checkout**, bukan saat masuk keranjang. Menahan sejak
keranjang berarti barang hilang dari peredaran karena ada yang menimbunnya
berminggu-minggu. Konsekuensinya barang di keranjang bisa habis diambil orang
lain, dan itu harus dikomunikasikan jelas — bukan jadi kejutan. ADR-003.

`inventory_movements` adalah **buku besar stok**: setiap perubahan tercatat
beserta sebab dan saldo sesudahnya. `quantity_on_hand` adalah ringkasan yang bisa
dihitung ulang dari sini. Kalau keduanya tidak cocok, movements yang benar.

---

## 8. Model uang

Aturan tunggal: **`int64` rupiah, tidak ada float, tidak pernah.**

Persentase disimpan sebagai **basis poin**: 250 = 2,5%. Bilangan bulat, tidak ada
pecahan yang perlu disimpan.

Fungsi terpenting adalah `money.Distribute`, yang membagi satu nilai ke beberapa
bagian dengan jaminan **jumlah hasilnya persis sama dengan total**:

```
Voucher Rp10.000 untuk keranjang 3 penjual
  → dibagi proporsional ke 3 suborder
  → sisa pembulatan dialokasikan deterministik, tidak dibuang
  → jumlah 3 bagian = Rp10.000 persis
```

Dipakai di dua tempat kritis: membagi voucher ke suborder, dan membagi diskon
suborder ke tiap baris barang (yang dibutuhkan untuk refund parsial per barang).

Ini fungsi yang **wajib punya test properti**. ADR-005.

---

## 9. Alur uang

```
Pembayaran gateway:
  pembeli bayar  →  kas kita naik, kewajiban ke penjual naik (pending)
  suborder selesai  →  komisi dipotong, masuk revenue kita
  lewat masa tahan  →  pending penjual → available
  pencairan  →  available turun, kas kita turun

COD (arah terbalik):
  kurir terima uang  →  piutang ke kurir naik
  kurir setor        →  piutang turun, kas naik
```

Dicatat dalam **buku besar double entry**. Tidak ada kolom saldo yang di-UPDATE;
saldo adalah hasil penjumlahan entri, dan setiap jurnal wajib seimbang.

Kolom saldo yang diperbarui langsung akan menyimpang cepat atau lambat, dan saat
itu terjadi kamu tidak punya cara tahu angka mana yang benar. ADR-011.

COD punya alur yang berbeda total: uang tidak pernah lewat kita di awal. Sampai
setoran kurir masuk, yang ada di pembukuan adalah piutang. ADR-012.

Laporan rekonsiliasi harian membandingkan pembayaran masuk, pendapatan penjual,
komisi, dan pencairan. **Selisihnya harus nol.**

---

## 10. Tiga metode pengiriman

| Metode | Ongkir | Syarat |
|---|---|---|
| `pickup` | nol | outlet mendukung |
| `local_delivery` | tarif dasar + per km | outlet mendukung, alamat punya koordinat, dalam radius, di atas minimum |
| `courier` | dari agregator | outlet mendukung, berat varian terisi |

Ongkir dihitung **per suborder**, dari koordinat outletnya masing-masing. Satu
pesanan berisi tiga penjual menghasilkan tiga perhitungan berbeda.

Quote punya **masa berlaku**. Tarif kurir berubah; checkout dengan quote
kedaluwarsa harus menghitung ulang, bukan memakai angka basi — kalau tidak, ada
selisih ongkir yang harus ditanggung seseorang.

Kegagalan agregator kurir **tidak boleh menggagalkan seluruh checkout**.
Tampilkan opsi lain dan catat masalahnya. Kehilangan satu opsi lebih baik
daripada kehilangan pesanan.

---

## 11. Event dan outbox

Event ditulis sebagai baris database dalam transaksi yang sama dengan perubahan
datanya:

```
BEGIN
  INSERT INTO order_orders ...
  INSERT INTO inventory_reservations ...
  INSERT INTO outbox_events ...
COMMIT                        ← semuanya, atau tidak sama sekali
```

Relay memungutnya dengan `FOR UPDATE SKIP LOCKED`, jadi aman dijalankan banyak
instance.

Jaminannya **at-least-once**. Semua consumer wajib idempoten — dan di sistem ini
itu bukan soal kerapian: event yang diproses dua kali di `settlement` berarti
penjual dibayar dua kali.

Peristiwa penting terjadi di level **suborder**, bukan order. Consumer yang salah
mendengarkan level akan bereaksi terlalu dini:

| Consumer | Mendengarkan | Bukan |
|---|---|---|
| `inventory` commit stok | `SuborderShipped` | `OrderPaid` |
| `settlement` matangkan dana | `SuborderCompleted` | `OrderPaid` |
| `notification` ke penjual | `OrderPaid` | — |

---

## 12. Pekerjaan worker

`cmd/worker` menjalankan delapan pekerjaan:

| Pekerjaan | Kalau tidak jalan |
|---|---|
| Relay outbox | tidak ada event yang sampai ke consumer |
| Pelepas reservasi kedaluwarsa | stok habis di layar padahal ada di rak |
| Pembatalan pesanan tak dibayar | pesanan menggantung selamanya |
| Rekonsiliasi pembayaran | pesanan yang sudah dibayar tetap berstatus menunggu |
| Pematangan dana | dana penjual tidak pernah bisa dicairkan |
| Pemrosesan pencairan | uang tidak pernah sampai ke penjual |
| Sinkronisasi tracking | status pengiriman tidak pernah berubah |
| Penyegaran agregat katalog | rating dan jumlah terjual basi |

Semuanya bergejala sama: sistem terlihat normal padahal tidak ada yang terjadi.
Beri metrik dan alarm untuk masing-masing.

Pekerjaan yang menyentuh uang butuh **kunci terdistribusi** di Redis, dan
idempotency key di database sebagai pertahanan berikutnya. Redis bisa hilang;
unique index tidak. ADR-013.

---

## 13. Testing

| Level | Cakupan | Butuh |
|---|---|---|
| Unit | `domain/*`, `money` | tidak ada |
| Usecase | `app/*` dengan port palsu | tidak ada |
| Repository | `infra/*` | Postgres asli |
| Integrasi | checkout end to end | Postgres + gateway palsu |

Yang paling berharga di sistem ini:

**Test properti `money.Distribute`.** Untuk input acak apa pun, jumlah hasil
harus sama dengan total dan tidak ada bagian negatif. Bug di sini muncul sebagai
selisih satu rupiah yang tidak bisa dijelaskan ke penjual.

**Checkout bersamaan untuk stok terakhir.** Dua goroutine membeli barang yang
stoknya tinggal satu. Tepat satu harus berhasil.

**Idempotensi.** Checkout dengan kunci sama dua kali → satu pesanan. Webhook sama
dua kali → satu perubahan status. Event settlement dua kali → satu jurnal.

**Refund parsial.** Satu penjual menolak, dua lainnya jalan terus, jumlah refund
tepat, order induk tidak batal.

**Rekonsiliasi buku besar.** Setelah rangkaian transaksi, jumlah debit sama
dengan jumlah kredit dan saldo tiap akun sesuai harapan.

Lima hal itu tidak akan ketahuan dari klik manual, dan semuanya adalah jenis bug
yang muncul di produksi pada saat paling tidak nyaman.

---

## 14. Kalau nanti dipecah

Sinyal yang layak dipertimbangkan: satu modul punya profil beban yang jauh
berbeda, atau tim yang mengerjakannya sudah saling menghalangi.

Kandidat termudah: `notification`, `audit`, `review` — ketiganya konsumen murni,
jadi yang berubah hanya transport event-nya.

Kandidat berikutnya: `catalog` untuk sisi baca, kalau trafik browse jauh
melampaui trafik transaksi.

Yang **tidak** layak dipecah: `cart`, `order`, `inventory`, `promotion`,
`payment`. Kelimanya bertemu dalam satu transaksi di checkout, dan memisahkannya
berarti membangun saga untuk masalah yang bisa diselesaikan `BEGIN`/`COMMIT`.
