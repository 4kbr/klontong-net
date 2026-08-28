# Catatan Keputusan (ADR)

Catatan keputusan teknis beserta **alasan** dan **konsekuensinya**. Gunanya:
enam bulan lagi saat kamu bertanya "kenapa dulu begini?", jawabannya ada di sini.

## Cara memakai dokumen ini

Buat entri baru untuk keputusan yang **sulit dibalik** atau yang **kelihatan
aneh dari luar**. Bukan setiap keputusan.

```markdown
## ADR-0NN: Judul singkat

- **Tanggal:** YYYY-MM-DD
- **Status:** Diusulkan | Diterima | Ditolak | Digantikan oleh ADR-0XX

### Konteks
### Keputusan
### Alternatif
### Konsekuensi
```

**ADR tidak pernah dihapus atau diedit isinya.** Kalau keputusan berubah, tulis
yang baru dan ubah status yang lama jadi "Digantikan oleh ADR-0XX".

---

## ADR-001: Modular monolith, bukan microservice

- **Tanggal:** 2026-08-28
- **Status:** Diterima

### Konteks

Checkout menyentuh keranjang, katalog, harga, stok, promosi, ongkir, pesanan,
dan pembayaran — semuanya dalam satu tarikan napas, dan sebagian besar harus
berhasil atau gagal bersama. Tim kecil, belum ada data beban nyata.

### Keputusan

Satu binary, satu database, modul dengan batas eksplisit di dalamnya.

### Alternatif

**Microservice sejak awal.** Ditolak: menahan stok, memotong kuota voucher, dan
membuat pesanan harus atomik. Memisahkannya berarti membangun saga sejak hari
pertama untuk masalah yang belum kita punya.

**Monolith biasa tanpa modul.** Ditolak: tanpa batas yang dipaksakan, dalam enam
bulan handler pesanan akan memanggil repository stok langsung, dan pemecahan
jadi mustahil selamanya.

### Konsekuensi

- Transaksi lintas domain cukup `BEGIN`/`COMMIT`.
- Batas modul harus ditegakkan aktif (folder `internal/` + `depguard`).
- Scaling hanya horizontal untuk seluruh aplikasi. Diterima sampai ada bukti itu
  masalah.

---

## ADR-002: Satu Order induk, banyak Suborder per penjual

- **Tanggal:** 2026-08-28
- **Status:** Diterima

### Konteks

Ini marketplace. Satu keranjang berisi barang dari beberapa penjual, dari outlet
berbeda, dengan cara kirim yang bisa berbeda pula. Pembeli membayar sekali.

### Keputusan

```
Order (milik pembeli, satu pembayaran, satu alamat)
 ├── Suborder  penjual A, outlet A1, kurir ekspedisi
 ├── Suborder  penjual B, outlet B2, antar lokal
 └── Suborder  penjual C, outlet C1, ambil di toko
```

Suborder adalah **unit kerja sesungguhnya**: ongkir, status, pengiriman,
komisi, dan pencairan semuanya di level ini. Status Order induk adalah
**turunan** dari status anak-anaknya, bukan disetel langsung.

### Alternatif

**Satu Order datar dengan kolom seller_id.** Ditolak: tidak bisa menampung
pesanan multi-penjual sama sekali.

**Pecah jadi beberapa Order terpisah saat checkout.** Ditolak: pembeli menganggap
dirinya melakukan satu pembelian dan membayar satu kali. Memecahnya jadi tiga
pesanan berarti tiga tagihan, dan itu bertentangan dengan yang ia alami.

### Konsekuensi

- Ini keputusan yang menyentuh **hampir semua modul**. Setiap kali menulis kode
  yang berhubungan dengan pesanan, tanya dulu: ini level order atau suborder?
- Satu penjual menolak pesanan **tidak membatalkan** pesanan penjual lain.
  Yang terjadi: stok suborder itu dilepas, refund parsial, order induk tetap
  jalan. Menyamakannya dengan pembatalan penuh adalah bug yang akan membuat dua
  penjual lain kehilangan penjualan tanpa sebab.
- Refund hampir selalu parsial. Rancang untuk itu sejak awal.
- Pembeli melihat satu pesanan dengan beberapa kemajuan berbeda. Sampaikan apa
  adanya — "1 dari 3 toko sudah dikirim" lebih jujur dan lebih berguna daripada
  satu status tunggal yang menyederhanakan.
- Penjual mengakses `/suborders/{id}`, bukan `/orders/{id}`. Ia tidak pernah
  melihat bagian penjual lain.

---

## ADR-003: Stok ditahan saat checkout, bukan saat masuk keranjang

- **Tanggal:** 2026-08-28
- **Status:** Diterima

### Konteks

Kapan stok dikurangi menentukan pengalaman belanja dan berapa banyak barang yang
"hilang" dari peredaran.

### Keputusan

Stok ditahan (`reserved`) saat **checkout**, dengan masa berlaku
(`STOCK_RESERVATION_TTL`, default 30 menit). Worker melepas reservasi yang
kedaluwarsa. Stok benar-benar berkurang (`on_hand`) saat suborder **dikirim**.

`STOCK_RESERVATION_TTL` sengaja **lebih pendek** dari `PAYMENT_EXPIRY`.

### Alternatif

**Tahan sejak masuk keranjang.** Ditolak: barang hilang dari peredaran karena ada
yang menimbunnya di keranjang selama seminggu. Untuk barang klontong yang stoknya
tipis, ini mematikan.

**Jangan tahan sama sekali, kurangi saat bayar.** Ditolak: pembeli menyelesaikan
pembayaran lalu diberi tahu barangnya habis. Uang sudah masuk dan harus
dikembalikan — pengalaman terburuk dari semua pilihan.

### Konsekuensi

- Barang di keranjang **bisa habis** diambil orang lain sebelum checkout. Itu
  harus dikomunikasikan jelas di halaman keranjang, bukan jadi kejutan saat
  tombol checkout ditekan.
- **Worker pelepas reservasi wajib jalan.** Kalau tidak, stok habis di layar
  padahal barangnya ada di rak — dan tidak ada satu pun pesan error yang muncul.
  Ini kegagalan senyap, jadi beri metrik dan alarm.
- Urutan TTL itu disengaja: reservasi lepas lebih dulu, baru pesanan kedaluwarsa.
  Kalau terbalik, ada jendela waktu di mana stok tertahan untuk pesanan yang
  sudah mati. Validasi urutan ini di `config.Load()`.
- `on_hand` dan `reserved` disimpan terpisah supaya penjual bisa tahu "ada di rak
  tapi sudah dipesan" — pertanyaan yang pasti muncul.

---

## ADR-004: Harga dan total selalu dihitung ulang di server

- **Tanggal:** 2026-08-28
- **Status:** Diterima

### Konteks

Klien menampilkan harga, subtotal, ongkir, diskon, dan total. Godaannya adalah
menerima angka-angka itu saat checkout supaya tidak perlu menghitung dua kali.

### Keputusan

Klien mengirim **pilihan**, bukan angka: alamat mana, kurir mana, voucher apa,
barang apa dan berapa banyak. Server menghitung ulang semuanya dari nol.

Struct request checkout **tidak punya field harga atau total sama sekali**.

Klien boleh mengirim total yang ia tampilkan sebagai **pembanding**; kalau
berbeda dari hitungan server, checkout ditolak dengan `ErrPriceChanged` beserta
rinciannya — bukan diam-diam menagih jumlah yang berbeda dari yang dilihat
pembeli.

### Alternatif

**Percaya angka dari klien.** Ditolak, dan bukan karena teoretis: siapa pun yang
membuka devtools bisa membeli seharga Rp1. Ini kerentanan e-commerce paling
klasik dan paling sering ditemukan.

**Tandatangani angka dari server lalu terima kembali.** Lebih baik, tapi tetap
membekukan harga selama tanda tangan berlaku dan menambah kompleksitas tanpa
menghilangkan kebutuhan menghitung ulang.

### Konsekuensi

- Perhitungan harga jalur tampilan dan jalur checkout **wajib memakai kode yang
  sama** (`pricing.Resolve`). Kalau keduanya menghitung dengan cara berbeda,
  cepat atau lambat hasilnya berbeda, dan pembeli yang menemukannya.
- Harga bisa berubah antara pembeli melihat dan menekan bayar. Itu kondisi
  normal, bukan error — tampilkan perubahannya dengan jelas dan minta konfirmasi.
- Checkout jadi lebih berat karena menghitung ulang seluruh keranjang. Itu harga
  yang layak dibayar.

---

## ADR-005: Uang sebagai int64 rupiah, persentase sebagai basis poin

- **Tanggal:** 2026-08-28
- **Status:** Diterima

### Konteks

Setiap sistem yang menyentuh uang menghadapi pertanyaan ini, dan yang salah
memilih akan menemukan invoice yang jumlahnya tidak cocok berbulan-bulan
kemudian.

### Keputusan

- Uang disimpan dan dihitung sebagai **`int64` rupiah**. Tidak ada `float`,
  tidak pernah, di mana pun di jalur uang.
- Persentase (komisi, diskon) disimpan sebagai **basis poin** `int`.
  250 = 2,5%. 1250 = 12,5%.
- Pembagian yang menyisakan memakai `money.Distribute`, yang menjamin **jumlah
  hasilnya persis sama dengan total** dan mengalokasikan sisa secara
  deterministik.
- Aturan pembulatan ditulis eksplisit di `money/rounding.go`, bukan dibiarkan ke
  perilaku default.

### Alternatif

**`float64`.** Ditolak. Ini bukan preferensi gaya — `0.1 + 0.2 != 0.3` dan
akumulasinya di ribuan transaksi menghasilkan selisih yang tidak bisa dijelaskan
ke penjual.

**`decimal` untuk semua.** Dipertimbangkan serius. Ditolak sebagai penyimpanan
utama karena rupiah tidak punya satuan pecahan dalam praktik, dan `int64` lebih
cepat, lebih sederhana, serta tidak punya perilaku pembulatan tersembunyi.
`decimal` boleh dipakai untuk **kuantitas** (0,5 kg beras itu sah) dan sebagai
perantara perhitungan persen — tapi hasilnya selalu kembali ke `int64`.

### Konsekuensi

- `money.Distribute` adalah fungsi yang **wajib punya test properti**: untuk
  input acak apa pun, jumlah hasil harus sama dengan total dan tidak ada bagian
  negatif. Ia dipakai di dua tempat kritis — membagi voucher ke suborder, dan
  membagi diskon suborder ke tiap baris barang untuk keperluan refund parsial.
- Sisa pembulatan **tidak boleh dibuang**. Rp10.000 dibagi tiga menyisakan
  1 rupiah, dan rupiah itu harus jatuh ke salah satu bagian secara deterministik.
- Komisi dibulatkan **ke bawah**: marketplace tidak boleh mengambil lebih dari
  haknya karena pembulatan.
- Nilai uang dikirim di JSON sebagai **angka**, bukan string terformat.
  "Rp12.500" adalah urusan frontend.
- `depguard` menolak `math/big` di luar `platform/money` supaya jalur uang tidak
  diam-diam berubah.

---

## ADR-006: Stok disimpan dalam satuan terkecil

- **Tanggal:** 2026-08-28
- **Status:** Diterima

### Konteks

Ini ciri khas dagang klontong: satu barang dijual per pcs, per renceng (isi 10),
dan per dus (isi 40) sekaligus. Ketiganya adalah barang fisik yang sama di rak
yang sama.

### Keputusan

Varian menyimpan `unit_code` (satuan jual), `content_quantity` (isi), dan
`base_unit_code` (satuan terkecil). **Stok disimpan dalam satuan terkecil.**
Penjualan 2 dus mengurangi stok 80 pcs.

Konversi terjadi di `catalog`, lewat `Variant.ToBaseQuantity()`. Modul
`inventory` hanya mengenal satuan dasar dan tidak tahu apa-apa soal dus.

### Alternatif

**Stok per varian, masing-masing berdiri sendiri.** Ditolak: stok tiga varian
yang sebenarnya barang yang sama akan saling berbohong. Penjual menjual 40 pcs,
sistem masih mengira ada 1 dus, dan pembeli berikutnya memesan barang yang tidak
ada.

**Satu varian saja, konversi di frontend.** Ditolak: harga per dus tidak selalu
40× harga per pcs — justru itu inti dagang grosir. Tier harga butuh varian
terpisah.

### Konsekuensi

- Semua kuantitas di `inventory.Port` dalam **satuan dasar**. Pemanggil yang
  mengonversi. Menaruh konversi di `inventory` akan menyebarkan pengetahuan
  tentang satuan ke modul yang tidak perlu tahu.
- `order.Item` menyimpan salinan `content_quantity` supaya `BaseQuantity()` tetap
  benar meski penjual mengubah definisi varian setelahnya.
- `content_quantity` wajib > 0. Kalau `unit_code == base_unit_code`, nilainya
  harus 1.
- Satuan yang tidak mengizinkan pecahan (dus, renceng) harus menolak kuantitas
  pecahan di keranjang. 0,5 kg beras sah; 0,5 dus tidak.

---

## ADR-007: Satu suborder mengambil dari satu outlet

- **Tanggal:** 2026-08-28
- **Status:** Diterima

### Konteks

Satu penjual bisa punya beberapa outlet dengan stok berbeda-beda. Pesanan yang
tidak muat di satu outlet secara teknis bisa dipenuhi dari dua outlet.

### Keputusan

Satu suborder terikat ke **satu outlet**. Kalau stok satu outlet tidak cukup
untuk seluruh baris penjual itu, checkout gagal dengan pesan yang jelas alih-alih
diam-diam memecahnya.

Pemilihan outlet: yang stoknya cukup, aktif, mendukung metode kirim, dan
**terdekat** dengan alamat pembeli.

### Alternatif

**Satu suborder boleh dari beberapa outlet.** Ditolak untuk sekarang: berarti
satu suborder punya beberapa pengiriman dengan ongkir masing-masing, dan
seluruh model fulfillment jadi berlipat rumitnya. Manfaatnya belum terbukti
sebanding.

### Konsekuensi

- Ada pesanan yang gagal padahal barangnya ada, hanya tersebar di dua outlet.
  Frekuensinya perlu dipantau — kalau sering, keputusan ini layak ditinjau ulang.
- Mengubah alamat pengiriman bisa mengubah outlet yang dipilih, dan karenanya
  mengubah ongkir. Panggil ulang pemilihan outlet setiap kali alamat berubah.
- Ongkir dihitung **per suborder** dari koordinat outletnya masing-masing.

---

## ADR-008: Idempotensi wajib di setiap jalur yang menyentuh uang

- **Tanggal:** 2026-08-28
- **Status:** Diterima

### Konteks

Pembeli menekan tombol bayar dua kali. Jaringan timeout lalu aplikasi mengulang.
Gateway mengirim webhook yang sama lima kali. Worker restart di tengah pencairan.
Semuanya normal, semuanya akan terjadi.

### Keputusan

Idempotensi ditegakkan **berlapis**, dan lapisan yang mengikat selalu di database:

| Jalur | Kunci | Penjaga |
|---|---|---|
| Checkout | header `Idempotency-Key` | Redis + unique index nomor pesanan |
| Pembuatan pembayaran | `payment.idempotency_key` | unique index |
| Webhook gateway | `(provider, event_id)` | unique index |
| Pencairan | `payout.idempotency_key` | unique index + idempotency key ke penyedia |
| Handler event | `event_id` | `ON CONFLICT DO NOTHING` |

Selain itu, **jangan pernah mengandalkan webhook sebagai satu-satunya jalur**.
Worker rekonsiliasi menanyakan status ke gateway untuk pembayaran yang
menggantung.

### Alternatif

**Andalkan Redis saja.** Ditolak: Redis boleh hilang, dan saat itu terjadi
perlindungannya hilang tepat ketika paling dibutuhkan.

**Andalkan webhook saja.** Ditolak: webhook hilang. Jaringan putus, kita sedang
deploy, gateway bermasalah. Tanpa rekonsiliasi, ada pesanan yang sudah dibayar
tapi selamanya berstatus menunggu — dan yang menemukannya adalah pembeli yang
marah, bukan monitoring kita.

### Konsekuensi

- `POST /checkout` **menolak** request tanpa `Idempotency-Key`. Jangan diam-diam
  melanjutkan.
- Webhook duplikat dibalas **200**, bukan error. Gateway mengirim ulang sampai
  mendapat 200; membalas error membuat mereka mengirim selamanya.
- Kunci idempotensi dibuat **sekali** saat objek dibuat, bukan setiap percobaan.
  Kunci yang berubah tiap retry menghilangkan seluruh perlindungannya.
- Rekonsiliasi harian membandingkan total kita dengan laporan gateway. Selisih
  harus nol.

---

## ADR-009: Transisi status lewat state machine eksplisit

- **Tanggal:** 2026-08-28
- **Status:** Diterima

### Konteks

Suborder punya sembilan status dan transisinya tidak bebas. Dari `shipped` tidak
boleh kembali ke `confirmed`; dari `cancelled` tidak boleh ke mana-mana.

### Keputusan

Peta transisi ditulis sebagai **data** di `order/internal/domain/state.go`.
Perubahan status hanya lewat `Suborder.Transition(to, now)`, yang memeriksa peta
itu. Menyetel field `Status` langsung dari usecase dilarang.

Setiap transisi menulis baris di `order_status_events`.

### Alternatif

**`if` tersebar di tiap usecase.** Ditolak: dalam tiga bulan tidak ada yang tahu
lagi transisi mana yang sah, dan akan muncul pesanan yang statusnya mustahil.

### Konsekuensi

- Satu tempat untuk diperiksa saat bertanya "kenapa pesanan ini tidak bisa
  dibatalkan".
- `delivered` dan `completed` **dibedakan**. Yang pertama berarti barang sampai;
  yang kedua berarti tidak ada lagi sengketa dan dana boleh cair. Menyamakan
  keduanya berarti membayar penjual sebelum pembeli sempat komplain.
- Status Order induk **tidak pernah disetel langsung** — ia dihitung dari
  anak-anaknya lewat `SyncStatusFromSuborders()`.
- Riwayat transisi jadi bukti saat ada sengketa.

---

## ADR-010: Pesanan menyimpan snapshot, bukan referensi

- **Tanggal:** 2026-08-28
- **Status:** Diterima

### Konteks

Pesanan merujuk ke produk, harga, dan alamat. Semuanya bisa berubah atau dihapus
setelah pesanan dibuat.

### Keputusan

`order_items` menyimpan **salinan** nama produk, nama varian, SKU, foto, satuan,
isi per satuan, harga satuan, dan berat. `orders.shipping_address` menyimpan
salinan alamat sebagai jsonb. `suborders.commission_bps` membekukan komisi saat
pesanan dibuat.

### Alternatif

**Simpan id saja, join saat menampilkan.** Ditolak: membuka pesanan tahun lalu
akan menampilkan harga hari ini. Produk yang dihapus penjual membuat riwayat
pesanan kosong. Invoice berubah setiap kali penjual mengedit nama produk.

### Konsekuensi

- Tabel pesanan lebih besar. Diterima — ini bukan tempat berhemat.
- Menampilkan pesanan **tidak perlu join** ke katalog, dan itu justru lebih cepat.
- Perubahan komisi marketplace **tidak mengubah** bagi hasil pesanan yang sudah
  jalan. Ini penting secara hukum, bukan hanya teknis.
- Foto produk yang dihapus dari storage akan menghasilkan tautan mati di pesanan
  lama. Pertimbangkan menyimpan salinan foto untuk pesanan, atau jangan pernah
  menghapus foto yang pernah dipesan.

---

## ADR-011: Buku besar double entry, bukan kolom saldo

- **Tanggal:** 2026-08-28
- **Status:** Diterima

### Konteks

Marketplace menahan uang pembeli, memotong komisi, lalu mencairkan sisanya ke
penjual. Ada saldo tertahan, saldo tersedia, komisi, refund, dan pencairan.

### Keputusan

Buku besar double entry di `settlement_entries`. **Tidak ada kolom saldo yang
di-UPDATE.** Saldo adalah hasil penjumlahan entri. Setiap jurnal wajib seimbang:
jumlah debit sama dengan jumlah kredit.

Kalau penjumlahan jadi lambat, buat **snapshot saldo** berkala — snapshot tetap
bisa dihitung ulang dari entri.

### Alternatif

**Kolom `balance` yang di-UPDATE.** Ditolak. Ia akan menyimpang cepat atau
lambat — satu bug, satu race, satu proses yang mati di tengah. Dan saat itu
terjadi, kamu tidak punya cara tahu angka mana yang benar. Dengan buku besar,
saldo selalu bisa dihitung ulang dan setiap perubahan punya jejaknya.

### Konsekuensi

- Setiap pergerakan uang menulis jurnal seimbang. Validasi keseimbangannya
  sebelum commit dan **gagalkan transaksi** kalau tidak seimbang — jangan
  dipaksakan.
- Pertimbangkan mencabut hak `UPDATE` dan `DELETE` pada tabel entri di level role
  database. Disiplin kode saja tidak cukup untuk tabel yang gunanya jadi bukti.
- Refund setelah dana cair menghasilkan **saldo negatif atau piutang**, bukan
  error. Itu kasus nyata dan sistem harus bisa menanganinya.
- Laporan rekonsiliasi harian membandingkan pembayaran masuk, pendapatan penjual,
  komisi, dan pencairan. Selisihnya harus nol; kalau tidak, selidiki hari itu juga.

---

## ADR-012: COD punya alur uang yang terbalik

- **Tanggal:** 2026-08-28
- **Status:** Diterima

### Konteks

Di segmen ini COD hampir selalu ada. Tapi arah arus kasnya berbeda total dari
pembayaran gateway.

### Keputusan

Untuk pembayaran gateway: uang masuk ke kita dulu, ditahan, lalu diteruskan ke
penjual setelah masa tahan.

Untuk COD: uang **tidak pernah lewat kita** di awal. Kurir menerimanya dari
pembeli, lalu menyetorkan ke kita belakangan. Sampai setoran itu masuk, yang ada
di pembukuan adalah **piutang ke kurir** (`courier/cod_receivable`).

Pesanan COD sah tanpa pembayaran di muka: `payment.status` tetap `pending` dan
baru `settled` saat setoran diterima.

### Alternatif

**Perlakukan COD sama seperti gateway.** Ditolak: berarti mencatat kas masuk yang
belum benar-benar ada, dan pembukuan tidak akan pernah cocok dengan rekening.

**Tidak menyediakan COD.** Ditolak berdasarkan kebutuhan pasar.

### Konsekuensi

- Alur `settlement` bercabang berdasarkan metode pembayaran. Tulis jurnalnya
  eksplisit dan berbeda; jangan memaksakan satu jalur.
- Ada **risiko gagal bayar** yang tidak ada di gateway: barang diantar, pembeli
  menolak membayar. Putuskan siapa yang menanggung — penjual, kurir, atau
  marketplace — dan catat sebagai ADR tersendiri sebelum COD dinyalakan.
- Rekonsiliasi setoran kurir adalah pekerjaan rutin tersendiri. Sediakan
  laporannya sejak awal.
- Pertimbangkan batas nilai untuk COD. Pesanan COD Rp5 juta dari akun baru adalah
  risiko yang tidak sebanding.

---

## ADR-013: Worker sebagai proses terpisah

- **Tanggal:** 2026-08-28
- **Status:** Diterima

### Konteks

Ada delapan pekerjaan latar: relay outbox, pelepas reservasi, pembatalan pesanan
kedaluwarsa, rekonsiliasi pembayaran, pematangan dana, pencairan, sinkronisasi
tracking, dan penyegaran agregat katalog.

### Keputusan

Proses terpisah, `cmd/worker`.

### Alternatif

**Goroutine di dalam API.** Ditolak: setiap deploy API memutus pekerjaan yang
sedang berjalan, beban keduanya tidak bisa diskalakan sendiri-sendiri, dan
pencairan dana yang terpotong deploy adalah masalah yang mahal.

### Konsekuensi

- Dua binary untuk di-deploy.
- Semua pekerjaan yang berebut baris memakai `FOR UPDATE SKIP LOCKED`.
- **Pekerjaan yang menyentuh uang butuh kunci terdistribusi** di Redis, dan
  idempotency key di database sebagai pertahanan berikutnya. Redis bisa hilang;
  unique index tidak.
- Saat development, `make dev` hanya menjalankan API. Jalankan `make worker` di
  terminal lain. Kalau lupa: reservasi stok tidak pernah lepas, pesanan tidak
  pernah kedaluwarsa, dana tidak pernah cair. Ketiganya bergejala sama — terlihat
  normal padahal tidak ada yang terjadi.

---

## ADR-014: Panggil REST penyedia langsung, bukan SDK

- **Tanggal:** 2026-08-28
- **Status:** Diusulkan

### Konteks

Payment gateway dan agregator ongkir punya SDK Go resmi. Pemakaian kita sempit:
buat transaksi, cek status, refund, verifikasi webhook, hitung tarif, buat
booking, tarik tracking.

### Keputusan

Panggil REST langsung dengan `net/http`, dibungkus di balik interface domain
(`payment.Gateway`, `fulfillment.CourierProvider`).

### Alternatif

**Pakai SDK resmi.** Untuk permukaan sesempit ini, membaca kode yang jelas
memanggil endpointnya lebih mudah di-debug daripada menelusuri lapisan
abstraksi. SDK juga mengikat bentuk kita ke bentuk mereka, dan mengganti penyedia
jadi jauh lebih mahal.

### Konsekuensi

- Kendali penuh atas timeout, retry, dan logging — semuanya penting di jalur
  checkout.
- Kita menulis sendiri struct request/response untuk endpoint yang dipakai.
  Jumlahnya sedikit.
- Perubahan API penyedia tidak otomatis terbawa. Imbangi dengan integration test
  terhadap lingkungan sandbox mereka.
- Interface di domain memungkinkan **implementasi palsu** untuk development.
  Tanpa itu, menguji alur pesanan di lokal butuh tunnel dan akun sandbox setiap
  kali — dan itu cukup menyebalkan untuk membuat orang berhenti mengetesnya.
- **Status masih Diusulkan.** Ubah jadi Diterima setelah integrasi pertama
  selesai dan terbukti nyaman.

---

## ADR-015: (template — salin untuk keputusan berikutnya)

- **Tanggal:**
- **Status:** Diusulkan

### Konteks

### Keputusan

### Alternatif

### Konsekuensi

---

## Yang masih terbuka

Catat di sini supaya tidak hilang, lalu ubah jadi ADR saat diputuskan:

- **Komisi dihitung dari nilai barang saja atau termasuk ongkir?** Rekomendasi:
  barang saja — mengambil komisi dari ongkir yang diteruskan ke kurir berarti
  penjual rugi.
- **Siapa menanggung gagal bayar COD?** Harus diputuskan sebelum COD dinyalakan.
- **Jarak antar lokal: Haversine dengan faktor pengali, atau API routing?**
  Haversine menghitung garis lurus; jarak tempuh nyata bisa 1,3–1,6 kali lipat.
- **Batas waktu konfirmasi penjual berapa lama?** Setelah itu suborder ditolak
  otomatis.
- **Masa komplain setelah barang diterima berapa lama?** Ini yang menentukan
  kapan dana penjual matang.
- **Katalog master bersama.** Saat ini tiap penjual punya baris produknya
  sendiri. Kalau nanti ingin menyatukan halaman produk lintas penjual, itu
  perubahan besar yang butuh ADR sendiri.
