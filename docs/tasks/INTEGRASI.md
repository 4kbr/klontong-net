# Fase Integrasi — Lepas Mock, Sambung ke Backend Asli

> Lintas-app. Prasyarat: backend fase 12 **dan** frontend fase 05 selesai.
> Kontrak: [`contracts/openapi`](../../contracts/openapi) — mengikat (ADR-015).

## Sifat fase ini: tidak memblokir

Backend dan frontend dikerjakan paralel dan **tidak saling menunggu**. Fase ini
dijalankan saat kedua sisi kebetulan sudah siap, bukan sebagai gerbang.

Kalau frontend sampai fase 05 sementara backend belum sampai fase 12: **frontend
lanjut ke fase 06 di atas MSW.** Jangan berhenti, jangan menunggu, jangan
menurunkan cakupan. Fase ini menunggu backend, bukan sebaliknya.

Kalau backend sampai fase 12 lebih dulu: backend lanjut ke fase 13. Jalur buyer
tetap bisa diuji sendiri lewat `curl` + `gateway_noop`.

## Tujuan

Membuktikan bahwa yang dibangun dua sisi di atas kontrak yang sama benar-benar
menyambung. **Selesai kalau** satu pesanan multi-penjual bisa dibuat dari
storefront asli, dibayar lewat gateway palsu, dikirim dari dashboard seller, dan
dikonfirmasi diterima pembeli — tanpa satu pun MSW handler aktif.

## Aturan khusus fase ini

- **Setiap beda diselesaikan ke arah kontrak** (ADR-015), bukan ke arah yang lebih
  mudah diubah. Kalau kontrak sendiri yang salah, ubah kontrak lebih dulu, lalu
  kedua sisi menyesuaikan.
- **Jangan menambal di frontend.** Field yang bentuknya tidak sesuai kontrak
  diperbaiki di backend, bukan dinormalisasi diam-diam di `client.ts`. Tambalan
  seperti itu membuat kontrak jadi fiksi.
- **Jangan melonggarkan kontrak** supaya implementasi lolos. Mengubah field jadi
  `nullable` karena backend belum mengisinya adalah cara kontrak kehilangan
  gunanya.
- Perbedaan yang ditemukan dicatat semua dulu, baru diperbaiki. Memperbaiki satu
  per satu sambil jalan menyembunyikan pola.

## Urutan kerja

### Persiapan
- [ ] Backend jalan: `make up`, `make migrate`, `make dev`, dan `make worker` di
      terminal lain. **Worker wajib** — tanpa itu reservasi stok tidak pernah lepas
      dan pesanan tidak pernah kedaluwarsa, dan gejalanya terlihat seperti bug
      frontend.
- [ ] `PAYMENT_GATEWAY=noop` supaya pembayaran bisa dipicu manual tanpa tunnel.
- [ ] CORS backend mengizinkan origin dev frontend (storefront dan dashboard punya
      port sendiri). Kredensial dan header `Idempotency-Key`, `X-Request-ID` ikut
      diizinkan.
- [ ] `VITE_API_BASE_URL` kedua app diarahkan ke API asli (`http://localhost:8080`),
      bukan Prism `:4010`.
- [ ] MSW **dimatikan** di mode ini. Sediakan satu saklar env
      (mis. `VITE_USE_MSW=false`) supaya bisa bolak-balik tanpa mengedit kode.

### Data uji
- [ ] `make seed` menghasilkan data yang cukup untuk jalur multi-penjual: 1 pembeli
      beralamat berkoordinat, **3 penjual** dengan outlet berkoordinat berbeda,
      produk dengan varian dus/renceng/pcs (`content_quantity` benar), stok di tiap
      outlet, dan 1 voucher marketplace + 1 voucher penjual.
- [ ] Tiga metode kirim terwakili: satu outlet `pickup`, satu `local_delivery`
      (punya delivery zone), satu `courier`.

### Smoke jalur buyer (manual, lewat UI asli)
- [ ] Daftar → verifikasi → login. Refresh token berotasi, 401 memicu
      single-flight refresh, bukan logout langsung.
- [ ] Browse katalog tanpa login (`OptionalAuth` benar-benar opsional), cari,
      buka detail produk lewat `slug`.
- [ ] Isi keranjang dari **3 penjual**, pastikan pengelompokan datang dari server.
- [ ] `POST /checkout/preview` → tiga suborder, tiga ongkir, diskon voucher terbagi
      dan **jumlahnya persis** sama dengan nilai voucher.
- [ ] `POST /checkout` dengan `Idempotency-Key`. Kirim dua kali dengan kunci yang
      sama → **satu** pesanan.
- [ ] Picu webhook gateway palsu → pesanan jadi `paid`. Kirim webhook yang sama dua
      kali → balasan 200, status tidak berubah dua kali.
- [ ] Dashboard seller: confirm → pack → ship pada satu suborder. Stok `on_hand`
      berkurang **saat kirim**, bukan saat bayar, dan Movement tertulis.
- [ ] Satu penjual lain **menolak** suborder-nya. Pesanan induk **tidak batal**,
      dua suborder lain tetap jalan, stok suborder yang ditolak dilepas.
- [ ] Pembeli `confirm-received` pada suborder yang sudah `delivered`.
- [ ] Ulasan hanya bisa ditulis untuk pembelian yang `completed`.

### Conformance kontrak
- [ ] `make -C contracts bundle` lalu bandingkan respons server asli dengan
      `contracts/dist/openapi.yaml` untuk setiap endpoint yang dilalui smoke di atas.
      Yang diperiksa: nama field (`snake_case`), tipe uang (`integer`), format waktu
      (RFC3339 UTC), nilai enum status, bentuk envelope, bentuk error, dan
      `meta.next_cursor` / `meta.has_more` pada endpoint berdaftar.
- [ ] Header: `X-Request-ID` di echo pada semua respons; `Retry-After` muncul pada
      429; `POST /checkout` **menolak** request tanpa `Idempotency-Key`.
- [ ] Jalur error yang dipakai UI benar-benar keluar dari server dengan `code` yang
      dijanjikan: `price_changed` dan `out_of_stock` beserta `error.details`
      (`LineIssue`), `invalid_credential`, `user_suspended`.

## Test wajib

- **Idempotensi checkout end-to-end.** Dua request dengan `Idempotency-Key` sama →
  satu pesanan di database.
- **Webhook duplikat.** Dua webhook `(provider, event_id)` sama → satu perubahan
  status, keduanya dibalas 200.
- **Penolakan parsial.** Satu penjual menolak → dua suborder lain hidup, order induk
  tidak batal, stok yang ditolak kembali tersedia.
- **Konversi satuan lintas-app.** Beli 2 dus isi 40 dari storefront → stok di
  dashboard seller berkurang 80 satuan dasar.
- **E2E jalur buyer** (Playwright) dijalankan **dua kali**: sekali di atas MSW,
  sekali di atas backend asli. Test yang sama, dua target. Beda hasil = beda
  kontrak, dan itu temuan.

## Definition of Done

- [ ] Seluruh smoke jalur buyer lewat tanpa MSW aktif
- [ ] Nol beda tersisa antara respons server asli dan `contracts/dist/openapi.yaml`
      untuk endpoint yang dilalui
- [ ] `make -C contracts ci` hijau
- [ ] `make check` (backend) dan `pnpm check` (frontend) hijau
- [ ] Setiap beda yang ditemukan punya PR perbaikan di sisi yang menyimpang, bukan
      tambalan di sisi lain
- [ ] Laporan: daftar beda yang ditemukan, di sisi mana diperbaiki, dan apa yang
      sengaja dibiarkan beserta alasannya

## Sengaja tidak dikerjakan di fase ini

- Gateway pembayaran sungguhan dan agregator kurir sungguhan — tetap palsu.
  Integrasi penyedia asli punya fase sendiri dan butuh `make tunnel`.
- Performa, SEO, dan a11y — itu frontend fase 10.
- Settlement dan payout end-to-end (backend fase 14) — kalau backend sudah sampai
  sana, tambahkan rekonsiliasi ke smoke; kalau belum, jangan menunggu.
