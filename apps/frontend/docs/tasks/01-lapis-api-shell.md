# Fase 01 — Lapis API & kerangka app

> Prasyarat: 00 | Ref: `contracts/openapi/components` (envelope, responses,
> parameters), `TASKS.md` (aturan global: Envelope, Auth, Lapis API, Paginasi)

## Tujuan

Membangun `@klontong/api` (client, error, hooks dasar, skema) dan kerangka
runtime kedua app (Query provider, MSW, routing shell, layout, guard, head
manager). **Selesai kalau** tiap app punya satu halaman *smoke* yang memanggil
endpoint mock lewat hook dari `@klontong/api` dan menampilkan data, plus rute
terproteksi mengarahkan ke login saat tak ada sesi.

## Aturan khusus fase ini

- `client.ts` adalah **satu-satunya** tempat `fetch` dipanggil di seluruh repo.
- `client.ts` meng-*unwrap* `{ data, meta }` dan melempar `ApiError` (berisi
  `code`, `message`, `fields`, `retryable`, `requestId`, `status`) untuk `{ error }`.
- Refresh **single-flight**: banyak 401 bersamaan hanya memicu **satu** panggilan
  `/auth/refresh`; semua request menunggu hasilnya lalu retry **sekali**.
- MSW store **in-memory** — reset antar test; di dev boleh persist ke
  `sessionStorage` untuk kenyamanan.
- Belum ada fitur bisnis; cukup 1–2 endpoint untuk membuktikan alur.

## Urutan kerja

### Tipe & API (packages/api)
- [ ] `src/schema.d.ts` — hasil `pnpm gen:api` (sudah ada dari Fase 00, refresh).
- [ ] `src/http/config.ts` — baca `VITE_API_BASE_URL`; helper `requestId()` (uuid v4).
- [ ] `src/http/errors.ts` — `class ApiError`, `isApiError()`, tipe union `error.code` dari kontrak.
- [ ] `src/http/client.ts` — `request<T>(method, path, { query, body, headers, auth, idempotencyKey })`:
      susun URL + querystring, set `Authorization`, `X-Request-ID`, `Idempotency-Key`,
      `session_token` (keranjang tamu), parse envelope, map error, hook refresh 401.
- [ ] `src/http/auth-token.ts` — penyimpanan access token (memori) + refresh token
      (`localStorage`), `setSession()`, `clearSession()`, `getAccessToken()`.
- [ ] `src/query/queryClient.ts` — `QueryClient` default (retry: 1 non-4xx, `staleTime`),
      `src/query/keys.ts` — factory `queryKeys` terpusat.
- [ ] `src/money.ts` — `type Money = number`, `formatRupiah(v)`, `bpsToPercent(bps)`, `percentLabel(bps)`.
- [ ] `src/schemas/common.ts` — Zod: `metaSchema`, `apiErrorSchema`, helper `paginated(schema)`.
- [ ] `src/endpoints/health.api.ts` + `src/hooks/useHealth.ts` — endpoint smoke (`GET /healthz`, `GET /readyz`).

### State (Zustand / Query) — per app `src/stores/`
- [ ] `authStore.ts` — `{ user, status: 'anon'|'authing'|'auth', signIn(), signOut() }` (non-persist; refresh token di `@klontong/api`).
- [ ] `guestCartStore.ts` (storefront) — `session_token` di `localStorage`.
- [ ] Bootstrap: saat app mount, bila ada refresh token → `authRefresh()` → set sesi; else `anon`.

### Komponen (packages/ui)
- [ ] Inisialisasi shadcn primitives yang dipakai luas: `button`, `input`, `card`,
      `dialog`, `dropdown-menu`, `toast`/`sonner`, `skeleton`, `badge`, `table`, `form`.
- [ ] `Money.tsx` — render `formatRupiah`; prop `estimate?: boolean` → tambah label "perkiraan".
- [ ] `PageHead.tsx` — wrapper head manager (`react-helmet-async` atau `@unhead/react`): `title`, `description`, `og`, `jsonLd`.
- [ ] `EmptyState.tsx`, `ErrorState.tsx` (tampilkan `error.code` + tombol retry + `requestId` kecil), `Pagination`/`LoadMore.tsx`, `DataTable.tsx`, `Guarded` helpers.

### Halaman & rute (kedua app)
- [ ] `src/main.tsx` — `<QueryClientProvider>` + `<HelmetProvider>` + `<RouterProvider>` + MSW start (dev).
- [ ] `src/routes/router.tsx` — `createBrowserRouter`, rute lazy, `errorElement` global.
- [ ] `src/layouts/` — `RootLayout` (header/footer storefront; sidebar dashboard), `AuthLayout`.
- [ ] `src/routes/guards.tsx` — `RequireAuth` (redirect ke `/login?next=`), `RequireRole(role)` (dashboard).
- [ ] Halaman smoke `/` — panggil `useHealth()` / satu list mock, render lewat `<Money/>` bila relevan.

### Wiring
- [ ] `src/test/msw/` — `server.ts` (node, untuk Vitest), `browser.ts` (dev), `handlers/` per area, `db.ts` store in-memory, `fixtures/` diketik dari `schema.d.ts`.
- [ ] `src/test/setup.ts` — `beforeAll` start server, `afterEach` `resetHandlers` + reset `db`.
- [ ] Handler awal: `auth` (login/refresh/logout/me), `health`.

## Test wajib

- `client.ts`: envelope sukses → mengembalikan `data`; envelope error → melempar `ApiError` dengan `code` benar.
- Refresh single-flight: 3 request paralel kena 401 → `/auth/refresh` dipanggil **sekali**, ketiganya retry & sukses.
- Refresh gagal → `authStore` jadi `anon`, diarahkan ke `/login`.
- `formatRupiah(1234500)` → `"Rp1.234.500"`; `<Money estimate>` menampilkan label perkiraan.
- `RequireAuth` tanpa sesi → redirect `/login?next=/target`.
- MSW `db` benar-benar reset antar test (tidak bocor state).

## Sengaja TIDAK dikerjakan

- Endpoint & halaman fitur (auth UI, katalog, dst) — Fase 02+.
- Prerender/SSR — Fase 10.
