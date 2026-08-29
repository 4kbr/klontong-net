# Fase 00 — Fondasi & tooling

> Prasyarat: — | Ref: `contracts/README.md`, `docs/ARCHITECTURE.md`, `TASKS.md`
> (aturan global)

## Tujuan

Menyiapkan monorepo pnpm dengan dua app Vite kosong, paket bersama, tema, dan
pipeline pengecekan. **Selesai kalau** `staging` sudah di-merge, `pnpm install`
jalan, `pnpm gen:api` menghasilkan `packages/api/src/schema.d.ts`, `pnpm mock`
menghidupkan Prism di `:4010`, `pnpm check` hijau, dan kedua app (`storefront`,
`dashboard`) boot menampilkan halaman kosong.

## Aturan khusus fase ini

- **Merge/rebase `staging` dulu.** Tanpa `contracts/` tidak ada spec dan tidak ada
  mock. Commit hasil merge terpisah dari scaffold.
- Belum ada logika bisnis, belum ada panggilan API nyata. Hanya kerangka.
- Semua konfigurasi TS `strict: true`. Tidak ada `any`, tidak ada
  `skipLibCheck: false` yang menutup error asli.
- Tema warna emas didefinisikan **sekali** sebagai token Tailwind di
  `packages/ui`, bukan di tiap app.

## Urutan kerja

### Repo & workspace
- [ ] `git merge staging` (atau rebase) — bawa `contracts/`, `docs/` terbaru, `apps/backend/docs/`.
- [ ] `apps/frontend/pnpm-workspace.yaml` — `packages: [storefront, dashboard, packages/*]`.
- [ ] `apps/frontend/package.json` (root privat) — skrip:
      `gen:api` (`openapi-typescript ../../contracts/dist/openapi.yaml -o packages/api/src/schema.d.ts`, didahului `pnpm --dir ../../contracts run bundle`),
      `mock` (`pnpm --dir ../../contracts run mock`),
      `dev`, `build`, `lint`, `typecheck`, `test`,
      `check` (`typecheck && lint && test && build`).
- [ ] `apps/frontend/tsconfig.base.json` — `strict`, `moduleResolution: bundler`, path alias `@klontong/*`.
- [ ] `.gitignore` tambah `node_modules`, `dist`, `coverage`, `.env`, `*.local`.

### Paket bersama (kerangka)
- [ ] `packages/api/` — `package.json` (`name: @klontong/api`), `tsconfig.json`,
      `src/index.ts` re-export, placeholder `src/schema.d.ts` (diisi `gen:api`).
- [ ] `packages/ui/` — `package.json` (`name: @klontong/ui`), Tailwind preset
      `tailwind.preset.ts` dengan palet **emas** (`primary` = skala amber/gold,
      mis. `50..950` berpusat di `#B8860B`/`#D97706`), `src/globals.css`
      (`@tailwind` + CSS vars shadcn), `components.json` untuk shadcn.
- [ ] `packages/config/` (opsional) — `eslint-preset.js`, `tsconfig.json` dasar.

### App Vite (kerangka, x2: storefront + dashboard)
- [ ] `pnpm create vite@latest` template `react-ts` untuk `storefront/` dan `dashboard/`.
- [ ] Hapus boilerplate; `src/main.tsx` render `<App/>` kosong dengan `globals.css` dari `@klontong/ui`.
- [ ] `tailwind.config.ts` `presets: [uiPreset]`, `content` termasuk `packages/ui`.
- [ ] `vite.config.ts` — alias `@/`, `@klontong/api`, `@klontong/ui`; port berbeda (storefront 5173, dashboard 5174).
- [ ] `.env.example` per app: `VITE_API_BASE_URL=http://localhost:4010`.
- [ ] shadcn init di tiap app (`components/ui/` lokal, tapi token & preset dari `@klontong/ui`).

### Tooling
- [ ] ESLint flat config + `@typescript-eslint` + `eslint-plugin-react-hooks` +
      `eslint-plugin-import` (larang import relatif lintas paket) + Prettier.
- [ ] **Aturan lint kustom**: larang `fetch(`/`axios` di luar `packages/api`
      (mis. `no-restricted-imports` / `no-restricted-globals`).
- [ ] Vitest + `@testing-library/react` + `jsdom` + `@testing-library/jest-dom`;
      `vitest.config.ts` per app + `packages/api`.
- [ ] `pnpm gen:api` dijalankan; `schema.d.ts` di-commit.

## Test wajib

- Smoke: `pnpm --filter storefront build` dan `pnpm --filter dashboard build` sukses.
- `pnpm gen:api` menghasilkan file non-kosong dan `pnpm typecheck` tetap hijau.
- Satu test Vitest trivial per app (mis. render `<App/>` tanpa crash).
- Lint gagal bila ada `fetch(` ditambahkan di `storefront/src` (uji regresi aturan).

## Sengaja TIDAK dikerjakan

- `client.ts`, hooks Query, MSW, routing nyata, tema komponen — Fase 01.
- Halaman fitur apa pun.
