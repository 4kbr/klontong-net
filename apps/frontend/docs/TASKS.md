# TASKS — Urutan Implementasi Frontend

Frontend Klontong Net **belum ada**: `apps/frontend/` hanya berisi `README.md`.
Dokumen ini menerjemahkan rancangan arsitektur jadi **checklist per fase** yang
bisa langsung dikerjakan, meniru format `apps/backend/docs/TASKS.md`.

Sumber kebenaran bentuk API adalah **`contracts/`** (OpenAPI 3.0.3, contract-first,
M1 buyer + M2 seller/admin/webhook lengkap). Untuk *alasan* di balik aturan, baca
`docs/DECISIONS.md` (ADR) dan `docs/ARCHITECTURE.md`.

## Cara pakai

1. Kerjakan **berurutan**. Tiap file fase = satu atau beberapa PR/commit fokus.
2. **Jangan lompat.** Fase 04 (checkout) butuh 02 (auth) dan 03 (katalog).
   Dashboard (06–09) butuh 01 + 02.
3. Tiap fase punya file di [`tasks/`](tasks/) dengan template sama: Tujuan →
   Aturan khusus → Urutan kerja (Tipe & API → State → Komponen → Halaman & rute →
   Wiring) → Test wajib → Sengaja tidak dikerjakan.
4. Selesai satu fase, jalankan **Definition of Done** di bawah sebelum lanjut.

## Prasyarat repo

Branch `feature/frontend/master` ada di commit awal `edff74f` dan **belum punya
`contracts/`** (ada di `staging`, `dd646b6`). **Langkah pertama Fase 00**: merge/
rebase `staging` ke branch ini supaya `contracts/`, `docs/` terbaru, dan doc task
backend ikut. Tanpa `contracts/` tidak ada spec untuk generate tipe dan tidak ada
`make mock`.

## Keputusan arsitektur (dikunci)

| Area | Keputusan |
|---|---|
| Bundler | React + TypeScript + **Vite** |
| Struktur | **Dua app** dalam **pnpm workspace**: `storefront/` (guest + buyer), `dashboard/` (seller + admin) |
| State server | **TanStack Query** membungkus fungsi `*.api.ts` |
| State klien/UI | **Zustand** (sesi auth, token keranjang tamu, filter, UI) |
| Form + validasi | React Hook Form + `@hookform/resolvers` + **Zod** |
| UI | **shadcn/ui + Tailwind**, warna utama **emas/amber (kuning tua)** |
| Lapis API | Paket bersama **`@klontong/api`** (`packages/api`): generated types + fetch client + `endpoints/*.api.ts` + hooks Query + skema Zod. Komponen **tak pernah** `fetch`/`axios` langsung. |
| Mock dev | **MSW** (stateful, store in-memory) untuk alur interaktif; **Prism** (`contracts/`, `:4010`) untuk cek bentuk kontrak & fallback |
| Routing | **React Router v7** (SPA), `createBrowserRouter`, rute lazy per fitur, data lewat Query (tanpa loader) |
| SEO | CSR + head manager per rute (`<title>`/`<meta>`/JSON-LD); prerender/SSR ditinjau di Fase 10 |

## Struktur folder target

```
apps/frontend/
├── docs/{TASKS.md, tasks/00..10-*.md}
├── pnpm-workspace.yaml            storefront, dashboard, packages/*
├── package.json                  skrip root: gen:api, mock, check
├── tsconfig.base.json            strict
├── packages/
│   ├── api/  (@klontong/api)
│   │   ├── src/schema.d.ts       openapi-typescript dari contracts/dist/openapi.yaml
│   │   ├── src/client.ts         base URL, inject Bearer + X-Request-ID,
│   │   │                         refresh single-flight saat 401, unwrap {data,meta},
│   │   │                         map {error} -> ApiError bertipe (pakai error.code)
│   │   ├── src/endpoints/*.api.ts
│   │   ├── src/hooks/*.ts        TanStack Query hooks + queryKeys
│   │   ├── src/schemas/*.ts      Zod (form + validasi boundary respons kritis)
│   │   └── src/money.ts          type Money = number; formatRupiah(); bpsToPercent()
│   ├── ui/   (@klontong/ui)      shadcn + tema emas + primitives (Money, StatusBadge,
│   │                             EmptyState, ErrorState, Pagination, DataTable,
│   │                             FormField, PageHead)
│   └── config/                   eslint / tailwind preset / tsconfig (boleh inline)
├── storefront/src/{routes,features,layouts,stores,lib,test}
└── dashboard/src/{routes,features,layouts,stores,lib,test}
```

## 11 Fase

| # | Fase | App | Selesai kalau | Butuh |
|---|---|---|---|---|
| [00](tasks/00-fondasi.md) | Fondasi & tooling | root | merge `staging`; pnpm workspace + 2 app + `packages/*`; Tailwind + tema emas; shadcn init; ESLint/Prettier/TS strict; Vitest; `pnpm gen:api` → `schema.d.ts`; `pnpm mock` jalan; `pnpm check` hijau; kedua app boot | — |
| [01](tasks/01-lapis-api-shell.md) | Lapis API & kerangka app | `packages/api`, kedua app | `client.ts` (auth+refresh+error map+X-Request-ID), `QueryClient`, MSW server + store in-memory, skema Zod dasar, `money.ts`, shell routing + layout + guard + `PageHead`; 1 halaman smoke per app fetch dari MSW | 00 |
| [02](tasks/02-auth-akun.md) | Auth & akun | storefront (+ login dashboard) | register, login (email/telp), refresh, logout, verify-email, request/verify OTP, forgot/reset password; `/me`, ganti password, profile, alamat CRUD + set default; rute terproteksi redirect; merge keranjang tamu saat login | 01 |
| [03](tasks/03-storefront-katalog.md) | Storefront katalog (guest) | storefront | home, pohon kategori, telusur produk + pencarian + filter + paginasi keyset, detail produk (varian, tier grosir, tampilan harga, indikator stok), halaman toko publik; head/JSON-LD per rute | 01 |
| [04](tasks/04-keranjang-checkout.md) | Keranjang & checkout (buyer) | storefront | keranjang (grup per penjual dari server, peringatan stok ADR-003, ubah/hapus baris, merge), `shipping/quote`, voucher `validate`/`available`, `checkout/preview` (berkali-kali), `checkout` + `Idempotency-Key`, tangani `price_changed`/`out_of_stock` per baris + konfirmasi ulang (ADR-004) | 02, 03 |
| [05](tasks/05-pesanan-pembayaran.md) | Pesanan, pembayaran, ulasan, notifikasi (buyer) | storefront | daftar pesanan, detail pesanan (induk + array suborder, status per-suborder ADR-002), instruksi pembayaran + retry, batalkan pesanan, `confirm-received` per suborder, ulasan (buat/daftar/report/my-reviews), notifikasi (inbox/unread-count/read-all/read/preferences) | 04 |
| [06](tasks/06-dashboard-seller-shop.md) | Shell dashboard + Seller shop | dashboard | shell app + nav + guard peran, `seller/register`, `seller/me`, `payout-account`, `documents`, outlet CRUD (koordinat), members | 01, 02 |
| [07](tasks/07-seller-katalog-inventori.md) | Seller katalog & inventori | dashboard | produk/varian CRUD, alur `images/upload-url`, harga + tier, publish; stok per outlet, movements, transfer, opname + finish | 06 |
| [08](tasks/08-seller-pesanan-keuangan.md) | Seller pesanan, fulfillment & keuangan | dashboard | siklus suborder (confirm/reject/pack/ship/ready-for-pickup), shipment (item, confirm-pickup, proof), delivery-zone, voucher toko CRUD, saldo/earnings/payouts, balas ulasan | 06 |
| [09](tasks/09-admin.md) | Panel admin | dashboard | sellers + approve/reject/suspend, kategori, voucher marketplace, orders + cancel, payments + reconcile, refunds, settlements, payouts + retry, laporan rekonsiliasi, audit + audit by target | 06 |
| [10](tasks/10-pengerasan.md) | Pengerasan & rilis | kedua app | Playwright E2E (happy path buyer + alur fulfillment seller), pass aksesibilitas, audit state loading/empty/error, perf (code-split, gambar), keputusan prerender/SSR SEO, checklist pindah ke API asli, README/docs | 02–09 |

## Aturan global (berlaku di semua fase)

Ringkas dari `../../../AGENTS.md`, `docs/ARCHITECTURE.md`, `docs/DECISIONS.md`, dan
`contracts/README.md`. Kalau ragu, baca ulang di sana.

- **ADR-004 — klien tidak pernah menghitung total yang mengikat.** Request
  checkout hanya berisi *pilihan* (alamat, kurir, voucher, barang+qty), tanpa
  field harga/total. `client_grand_total` hanya pembanding → `409 price_changed`
  dengan `error.details` per baris; UI tampilkan lalu minta konfirmasi ulang.
  Perkiraan boleh ditampilkan **berlabel "perkiraan"**; angka yang dibayar selalu
  dari server. Jalur tampilan dan checkout memakai angka dari server yang sama.
- **ADR-002 — Order = induk + array Suborder.** Tiap suborder punya status &
  nominal sendiri. **Jangan kolapskan jadi satu status** — tampilkan "1 dari 3
  toko sudah dikirim". Rancang untuk refund/pembatalan **parsial**.
- **ADR-003 — stok ditahan saat checkout, bukan di keranjang.** Barang di
  keranjang bisa habis diambil orang lain. Komunikasikan **di halaman keranjang**,
  bukan kejutan saat tombol checkout ditekan.
- **Checkout 2 tahap.** `POST /api/v1/checkout/preview` idempoten & boleh
  berkali-kali (tak menahan apa pun). `POST /api/v1/checkout` **wajib** header
  `Idempotency-Key`; buat **sekali** (uuid) per percobaan checkout, pakai ulang
  saat retry — bukan tiap klik.
- **Uang** = `integer` rupiah (int64), selalu angka. Persen = basis poin
  (250 = 2,5%). `"Rp12.500"` murni frontend: `formatRupiah()` dari
  `@klontong/api/money`, **display-only**. **Tidak ada aritmatika uang di
  komponen.**
- **Lapis API.** Komponen/hook halaman **tak pernah** memanggil `fetch`/`axios`.
  Selalu lewat `endpoints/*.api.ts` di `@klontong/api`, dibungkus hook TanStack
  Query. Satu file `*.api.ts` per kelompok (mis. `product.api.ts`, `cart.api.ts`,
  `seller-catalog.api.ts`).
- **Bentuk response hanya dari tipe hasil generate** (`openapi-typescript` →
  `packages/api/src/schema.d.ts`). Jangan mengetik ulang bentuk response.
- **Keranjang dikelompokkan server per penjual** dan baris bermasalah ditandai
  server. Frontend **tidak mengelompokkan ulang** dan tidak menghitung ulang
  ketersediaan.
- **Auth.** `Authorization: Bearer <jwt>` (akses TTL 15 menit). Access token di
  **memori** (store Zustand non-persist). Refresh token dari **body respons**,
  **dirotasi tiap panggilan** → simpan di `localStorage`, ganti tiap refresh;
  pakai ulang token lama mencabut seluruh sesi. Saat 401: **single-flight**
  refresh lalu retry sekali; gagal → logout.
- **Keranjang tamu.** `session_token` di `localStorage`, dikirim sebagai header.
  Setelah login → `POST /api/v1/cart/merge` lalu hapus token tamu.
- **Envelope.** Sukses `{ data, meta }` (unwrap di `client.ts`). Error
  `{ error: { code, message, fields?, retryable? } }` → `ApiError` bertipe;
  cabang UI pakai `error.code` (`lower_snake_case`), bukan `message`. `fields`
  dipetakan ke error field form.
- **Paginasi keyset.** `?limit` + `?cursor` opaque; `meta.next_cursor` +
  `meta.has_more`. Pakai `useInfiniteQuery`. Bukan offset/halaman.
- **Empat audiens.** guest buyer & buyer di `storefront/`; seller & admin di
  `dashboard/`. Pembedaan lewat route group + layout + guard (`RequireAuth`,
  `RequireRole`).
- **State.** Server → TanStack Query (+ `queryKeys` terpusat + invalidation
  eksplisit). Klien/UI → Zustand. Form → RHF + Zod resolver.
- **Header.** `X-Request-ID` di-generate klien per request dan dicatat saat error
  (memudahkan korelasi dengan log backend). `Retry-After` dihormati pada 429.
- **Konvensi.** TypeScript `strict`, tanpa `any`. Komponen kecil & fokus; file
  yang membengkak = sinyal terlalu banyak tanggung jawab. Struktur per fitur
  (`features/<nama>/`). **Komentar Bahasa Indonesia**, ikut gaya repo. Tema warna
  utama emas didefinisikan sekali di `packages/ui`.

## Definition of Done per fase

- [ ] `pnpm typecheck` bersih
- [ ] `pnpm lint` bersih (tanpa `eslint-disable` baru)
- [ ] `pnpm test` hijau
- [ ] `pnpm build` hijau untuk app yang tersentuh
- [ ] MSW handler + fixture (diketik dari `schema.d.ts`) untuk endpoint baru fase ini ada
- [ ] Test wajib fase (lihat bagian "Test wajib" di file fase) ada dan lulus
- [ ] `pnpm gen:api` dijalankan bila kontrak berubah; `schema.d.ts` ikut di-commit
- [ ] Laporan singkat: file yang dibuat/diubah, keputusan yang diambil, yang sengaja ditunda

## Catatan kontrak yang perlu diperhatikan

- Campur **slug vs id**: `GET /api/v1/products/{slug}` tapi
  `GET /api/v1/products/{productId}/reviews`. Ikuti kontrak apa adanya; simpan
  keduanya (`slug` + `id`) di tipe produk hasil fetch.
- `order` routes seller memakai `/api/v1/seller/orders/{suborderId}` (bukan
  `/suborders/{id}`). Kontrak mengikuti `routes.go` backend — pakai path kontrak.
- Tidak ada webhook kurir; update tracking dari sisi frontend hanya **membaca**
  status (worker backend yang menyinkronkan).
- Prism bersifat **stateless**. Alur stateful (keranjang → checkout → pesanan,
  sesi auth) dikembangkan di atas **MSW** dengan store in-memory; Prism untuk
  cek cepat bentuk response saja.
