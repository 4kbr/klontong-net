# Fase 10 — Pengerasan & rilis

> Prasyarat: 02–09 | Ref: `docs/ARCHITECTURE.md` §13 (testing berharga),
> `README` frontend (SEO, ADR-004), `contracts/README.md` (pindah ke API asli).

## Tujuan

Menstabilkan kedua app untuk rilis: E2E alur kritis, aksesibilitas, konsistensi
state loading/empty/error, performa, keputusan SEO, dan checklist beralih dari
mock ke API asli. **Selesai kalau** suite E2E hijau di CI, baseline a11y
terpenuhi, dan ada dokumen "cara pindah ke backend asli".

## Aturan khusus fase ini

- E2E memakai **Playwright** terhadap app yang dijalankan dengan **MSW** (bukan
  Prism) supaya alur stateful deterministik.
- Jangan menambah fitur baru. Hanya perbaikan, uji, dan dokumentasi.
- Setiap perubahan tetap lewat `pnpm check`.

## Urutan kerja

### E2E (Playwright)
- [ ] Setup `@playwright/test`, config per app, fixture MSW aktif, data seed deterministik.
- [ ] Storefront happy path: buka produk (tamu) → daftar/login → tambah 2 barang
      beda penjual → `/keranjang` → `/checkout` (pilih alamat + kurir + voucher)
      → tangani `price_changed` sekali → buat pesanan → lihat instruksi bayar →
      "bayar" → status `paid` → `confirm-received` → tulis ulasan.
- [ ] Seller path: login seller → buat produk 3 varian → set harga + stok →
      publish → terima suborder → `pack` → `ship` → cek saldo → ajukan pencairan.
- [ ] Admin path: login admin → approve toko pending → refund parsial satu
      suborder → buka laporan rekonsiliasi (selisih nol).
- [ ] Negatif: checkout tanpa `Idempotency-Key` tak mungkin dari UI; retry
      jaringan tidak menggandakan pesanan (dobel-klik "buat pesanan").

### Aksesibilitas
- [ ] Integrasi `axe` (mis. `@axe-core/playwright`) di E2E untuk halaman utama.
- [ ] Audit manual: fokus keyboard di dialog checkout, label form, kontras tema
      emas (WCAG AA), `aria-live` untuk error checkout & toast.
- [ ] Perbaiki temuan; jadikan cek axe blocking di CI untuk halaman inti.

### Konsistensi state
- [ ] Audit tiap rute: ada state `loading` (skeleton), `empty` (`EmptyState`),
      `error` (`ErrorState` + `requestId` + retry). Tambal yang kurang.
- [ ] Standarkan penanganan `429` (`Retry-After`) dan `5xx`/`502` (pesan + retry).
- [ ] Pastikan `X-Request-ID` tercatat di semua `ErrorState`.

### Performa
- [ ] Code-split per route group; ukur bundle (`rollup-plugin-visualizer`).
- [ ] Lazy-load rute berat (checkout, dashboard tabel besar).
- [ ] Gambar produk: `loading="lazy"`, `srcset`/ukuran, placeholder.
- [ ] `staleTime`/`gcTime` Query ditinjau per resource; prefetch di hover kartu produk (opsional).

### SEO — keputusan
- [ ] Ukur kebutuhan nyata: apakah crawler target menjalankan JS? Bila tidak →
      pilih salah satu dan catat sebagai ADR frontend:
      (a) prerender build-time untuk `/produk/:slug`, `/produk`, `/toko/:slug`
      (mis. `vite-plugin-prerender`/`@prerenderer`), atau
      (b) migrasi storefront ke SSR (Vike). Dashboard tetap CSR.
- [ ] Implementasikan opsi terpilih **hanya untuk storefront**; verifikasi
      `<title>`, meta, JSON-LD ada di HTML awal.

### Pindah ke API asli
- [ ] Dokumen `apps/frontend/docs/GOING-LIVE.md`: set `VITE_API_BASE_URL` ke
      backend, matikan MSW (`import.meta.env` guard), jalankan `pnpm gen:api`
      terhadap `contracts` terbaru, jalankan E2E terhadap staging backend,
      daftar perbedaan perilaku Prism/MSW vs backend yang perlu dicek manual
      (idempotensi checkout, rotasi refresh, paginasi cursor nyata).
- [ ] CI: job `check` + E2E (MSW) wajib hijau sebelum merge.

## Test wajib

- Ketiga alur E2E (buyer, seller, admin) hijau di CI headless.
- Dobel-klik "buat pesanan" → tepat satu pesanan (Idempotency-Key dipakai ulang).
- axe: 0 pelanggaran serius di beranda, detail produk, checkout, dashboard pesanan seller, panel admin.
- Build produksi kedua app sukses; tidak ada chunk > ambang yang disepakati tanpa alasan.
- Halaman produk hasil prerender/SSR (opsi terpilih) memuat metadata di HTML awal (uji via `curl`/Playwright `route`).

## Sengaja TIDAK dikerjakan

- Fitur produk baru.
- Integrasi analytics/monitoring pihak ketiga (di luar cakupan; catat sebagai backlog bila perlu).
