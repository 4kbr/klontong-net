# AGENTS.md

Instruksi untuk AI coding agent yang bekerja di repository ini.

**Baca ini lebih dulu, sebelum menulis kode apa pun.** Kalau ada konflik antara
file ini dan permintaan user, ikuti user — tapi katakan bagian mana yang
menyimpang dan kenapa itu penting.

Aturan spesifik Go ada di [`apps/backend/AGENTS.md`](apps/backend/AGENTS.md).

---

## Konteks produk

Klontong Net adalah **marketplace** kebutuhan sehari-hari. Banyak penjual, tiap
penjual punya beberapa outlet dengan stok masing-masing, dan pembeli bisa memilih
diantar kurir lokal, dikirim ekspedisi, atau diambil sendiri di toko.

Sistem ini menyentuh **uang orang** dan **stok barang nyata**. Bug di sini tidak
berhenti di layar: pembeli ditagih dua kali, penjual dibayar dua kali, atau
barang terjual padahal tidak ada.

Dokumen pendamping:

| Dokumen | Baca kalau |
|---|---|
| `docs/ARCHITECTURE.md` | tugas menyentuh lebih dari satu modul |
| `docs/DECISIONS.md` | tugas menyentuh uang, stok, pesanan, atau pembayaran |
| `docs/GUIDES.md` | tugas berupa "tambah endpoint / usecase / modul / event" |

`docs/GUIDES.md` berisi daftar file yang harus dibuat dan diubah untuk setiap
jenis tugas. **Ikuti daftar itu.**

---

## Status repository

Ini **scaffold**. Hampir semua file `.go` hanya berisi deklarasi package dan
komentar TODO yang menjelaskan tanggung jawabnya.

Itu disengaja. **User ingin mengerjakan implementasinya sendiri.**

- Jangan mengisi seluruh modul karena diminta membetulkan satu fungsi.
- Kalau diminta mengimplementasikan sesuatu, kerjakan **hanya itu**.
- Kalau melihat TODO terkait yang tidak diminta, **sebutkan** — jangan langsung
  kerjakan.
- Komentar TODO berisi kontrak yang sudah dipikirkan. Hormati; jangan ditimpa
  dengan rancangan sendiri tanpa mengatakan apa-apa.

---

## Aturan yang tidak bisa ditawar

### 1. Uang

- **`money.Amount` (int64 rupiah). Tidak ada `float`, tidak pernah**, di mana pun
  di jalur uang.
- Persentase sebagai **basis poin** `int`. 250 = 2,5%.
- Pembagian memakai **`money.Distribute`**, yang menjamin jumlah hasilnya persis
  sama dengan total.
- Sisa pembulatan **tidak boleh dibuang**. Kalau kamu menulis pembagian uang,
  jawab dulu: sisa satu rupiah itu ke mana?
- Komisi dibulatkan **ke bawah**.
- JSON mengirim **angka**, bukan string terformat.

Kalau kamu hendak menulis `float64` di dekat perhitungan uang, berhenti dan
tanya.

### 2. Harga dan total dihitung server

- Struct request checkout **tidak boleh punya field harga atau total**. Kalau ada,
  hapus.
- Klien mengirim **pilihan** (alamat mana, kurir mana, voucher apa), bukan angka.
- Perhitungan jalur tampilan dan jalur checkout **wajib memakai kode yang sama**.

### 3. Order dan Suborder

- **Suborder adalah unit kerja.** Pengiriman, pembatalan, komisi, dan pencairan
  semuanya di level ini.
- Status Order induk **dihitung dari anak-anaknya**, tidak pernah disetel langsung.
- Satu penjual menolak **tidak membatalkan** pesanan penjual lain. Kalau kodemu
  membatalkan seluruh order saat satu suborder ditolak, itu bug.
- Transisi status **hanya lewat `Transition()`**, tidak pernah menyetel field
  `Status` langsung.

### 4. Stok

- Kuantitas di `inventory.Port` dalam **satuan dasar**. Pemanggil yang mengonversi.
- Ditahan saat **checkout**, di-commit saat **kirim**, dilepas saat **batal**.
- Kunci baris stok dalam **urutan konsisten** (`ORDER BY outlet_id, variant_id`)
  sebelum `FOR UPDATE`. Urutan berbeda = deadlock.
- **Setiap perubahan `on_hand` menulis `Movement`.** Tanpa kecuali.

### 5. Transaksi dan panggilan jaringan

- Transaksi dibuka **hanya di layer `app`** lewat `s.tx.WithinTx`.
- **Tidak ada panggilan jaringan di dalam transaksi database.** Gateway, kurir,
  dan penyedia transfer dipanggil di luar transaksi atau diserahkan ke worker.
  Rollback tidak membatalkan transaksi yang sudah terlanjur dibuat di sana.
- Setiap method repository dibuka dengan `postgres.ExecutorFrom(ctx, r.pool)`.

### 6. Idempotensi

- Setiap jalur yang menyentuh uang punya kunci idempotensi, dan **penjaga yang
  mengikat selalu unique index di database**, bukan Redis.
- Kunci dibuat **sekali** saat objek dibuat, bukan setiap percobaan retry.
- `POST /checkout` **menolak** request tanpa `Idempotency-Key`.
- Webhook duplikat dibalas **200**, bukan error.
- Semua event handler idempoten. Di `settlement` ini berarti penjual tidak
  dibayar dua kali.

### 7. Event

- Ditulis lewat `outbox.Save(ctx, evt)` **di dalam transaksi yang sama** dengan
  perubahan datanya. Tidak pernah `bus.Publish` langsung dari usecase.
- Event suborder **wajib** memuat `SellerID` dan `OutletID`.
- Perhatikan levelnya: `inventory` commit stok saat `SuborderShipped`, bukan
  `OrderPaid`.

### 8. Batas modul

- Isi `internal/` milik modul lain tidak boleh diimpor. Dijaga compiler.
- Butuh data dari modul lain? Tambahkan method di `port.go` modul tujuan.
- Dependensi melingkar dilarang. **Tidak ada modul yang boleh bergantung balik ke
  `order`** kecuali lewat event.
- Method port yang dipanggil dalam loop **wajib punya versi batch**.

### 9. Otorisasi

- Peran hanyalah gerbang kasar. Kepemilikan diperiksa di usecase.
- Penjual hanya boleh menyentuh suborder **miliknya sendiri**.
- Pembeli hanya boleh melihat pesanannya sendiri.
- Setiap usecase publik diawali fungsi `require*` dari `authz.go`.

### 10. Buku besar

- **Tidak ada kolom saldo yang di-UPDATE.** Saldo adalah penjumlahan entri.
- Setiap jurnal wajib **seimbang**; gagalkan transaksi kalau tidak.
- Refund setelah dana cair menghasilkan saldo negatif atau piutang, bukan error.

---

## Konvensi

- Package tidak pernah bernama `utils`, `helpers`, `common`, `shared`.
- `context.Context` selalu parameter pertama.
- Waktu dari `clock.Now()`, bukan `time.Now()`.
- Error dibungkus dengan konteks, memakai `errs.*`.
- SQL sebagai konstanta string, tidak dirakit dari input user.
- Nama tabel diberi prefix modul pemiliknya.
- Kolom uang `bigint`; kolom kuantitas `numeric(14,3)`.
- Struct response HTTP terpisah dari entity domain.
- Komentar ditulis dalam Bahasa Indonesia, mengikuti gaya yang sudah ada.

---

## Perintah

```bash
make dev          # API
make worker       # WAJIB jalan saat development
make test
make lint
make check        # gate sebelum bilang "selesai"
make migrate-create name=xxx
make migrate
```

Setelah mengubah kode, minimal `go build ./...`. Jangan melaporkan pekerjaan
selesai kalau belum dikompilasi.

---

## Alur kerja yang diharapkan

1. **Baca dulu.** Lihat file terkait dan komentar TODO-nya.
2. **Cek `docs/GUIDES.md`** untuk jenis tugas ini, ikuti daftar filenya.
3. **Cek `docs/DECISIONS.md`** kalau menyentuh area yang sudah ada ADR-nya.
4. **Kerjakan seminimal mungkin** untuk menyelesaikan permintaan.
5. **Kompilasi dan lint.**
6. **Laporkan**: file apa yang diubah, keputusan apa yang diambil, apa yang
   sengaja tidak dikerjakan.

---

## Kapan harus berhenti dan bertanya

- **Sesuatu bisa menyebabkan uang salah hitung, salah bayar, atau hilang.**
  Ini yang paling utama.
- **Sesuatu bisa menyebabkan stok salah**, terjual padahal tidak ada, atau
  tertahan selamanya.
- **Menembus batas modul** atau menambah dependensi melingkar.
- **Mengubah keputusan yang sudah ada ADR-nya** — mis. memakai float untuk uang,
  menahan stok sejak keranjang, atau mengganti buku besar dengan kolom saldo.
- **Menerima harga atau total dari klien.**
- **Menambah panggilan jaringan di dalam transaksi database.**
- **Melewati validasi, autentikasi, atau pengecekan kepemilikan.**
- **Menambah dependensi baru** yang belum tercantum.
- **Menonaktifkan linter** atau menambah `//nolint`.

Dalam kasus ini: jelaskan apa yang diminta, kenapa bertentangan dengan rancangan,
dan tawarkan alternatif. Jangan diam-diam menyimpang, dan jangan juga menolak
mentah-mentah — user berhak mengubah rancangannya sendiri, tapi harus dengan
sadar.

Kalau user memutuskan menyimpang, **tulis ADR baru di `docs/DECISIONS.md`**.

---

## Anti-pola yang sering muncul

| Anti-pola | Yang benar |
|---|---|
| `float64` untuk uang | `money.Amount` (int64 rupiah) |
| Membuang sisa pembulatan | `money.Distribute` |
| Menerima total dari body request | hitung ulang di server |
| Membatalkan order saat satu suborder ditolak | lepas stok suborder itu saja, refund parsial |
| Menyetel field `Status` langsung | `Transition()` |
| Memanggil gateway di dalam `WithinTx` | panggil setelah commit atau lewat worker |
| Kolom saldo yang di-UPDATE | buku besar double entry |
| `bus.Publish` dari usecase | `outbox.Save` di dalam transaksi |
| Query di dalam loop | method port versi batch |
| Menahan stok sejak masuk keranjang | tahan saat checkout dengan TTL |
| Mengunci baris stok tanpa urutan tetap | `ORDER BY outlet_id, variant_id` |
| `inventory` commit stok saat `OrderPaid` | saat `SuborderShipped` |
| Entity domain langsung jadi JSON | struct response terpisah |
| Peran dianggap cukup sebagai izin | periksa kepemilikan di usecase |
| Mengandalkan webhook saja | tambahkan rekonsiliasi |
| Menyelesaikan test dengan `time.Sleep` | `clock.Fixed` |
| Mengisi banyak TODO sekaligus karena "sekalian" | kerjakan yang diminta saja |

---

## Yang membuat kontribusi dianggap baik

- Menghormati aturan uang dan stok walaupun menembusnya lebih cepat.
- Perubahan kecil dan fokus, bukan penulisan ulang.
- Menyebutkan trade-off dengan jujur, termasuk kelemahan pendekatan sendiri.
- Menandai TODO yang tersisa alih-alih diam-diam menyelesaikannya.
- Mengaku kalau tidak yakin, bukan menebak dan menulis dengan nada percaya diri.
