# Frontend

Monorepo pnpm untuk frontend Klontong Net. Dua aplikasi terpisah + paket bersama.

```
apps/frontend/
├── storefront/     Vite + React — pembeli (guest & login).   Port dev 5173
├── dashboard/      Vite + React — penjual & admin.            Port dev 5174
├── packages/
│   ├── api/        @klontong/api — tipe hasil generate dari kontrak OpenAPI,
│   │              (Fase 01+) client fetch, endpoint *.api.ts, hooks TanStack Query, skema Zod
│   └── ui/         @klontong/ui — shadcn/ui + Tailwind v4, tema emas, primitives bersama
└── docs/           TASKS.md (indeks fase) + tasks/00..10-*.md (checklist per fase)
```

Stack: React 19 + TypeScript (strict) + Vite 7 · Tailwind CSS v4 (CSS-first) +
shadcn/ui · TanStack Query (server state) + Zustand (state UI) · React Hook Form +
Zod · Vitest + Testing Library · MSW (mock dev stateful) / Prism (`contracts/`).

## Mulai

```bash
cd apps/frontend
pnpm install
pnpm gen:api        # bundle contracts + generate packages/api/src/schema.d.ts
pnpm mock           # Prism mock server di :4010 (terminal lain, opsional)
pnpm dev            # jalankan kedua app (storefront :5173, dashboard :5174)
```

Salin `.env.example` → `.env` di tiap app bila perlu mengubah `VITE_API_BASE_URL`.

## Perintah

| Perintah                                                    | Guna                                                            |
| ----------------------------------------------------------- | --------------------------------------------------------------- |
| `pnpm gen:api`                                              | Rakit kontrak → generate tipe TypeScript ke `packages/api`      |
| `pnpm mock`                                                 | Prism mock server dari `contracts/` di `:4010`                  |
| `pnpm dev`                                                  | Dev server kedua app (paralel)                                  |
| `pnpm typecheck` / `pnpm lint` / `pnpm test` / `pnpm build` | Per aspek                                                       |
| `pnpm check`                                                | Gate: typecheck + lint + test + build (jalankan sebelum commit) |

## Aturan yang mengikat

- **Frontend tidak pernah menghitung total yang mengikat** (ADR-004). Klien kirim
  pilihan; angka yang dibayar dari server.
- **Bentuk response hanya dari tipe hasil `pnpm gen:api`** — jangan diketik ulang.
- **Komponen tak pernah memanggil `fetch`/`axios`** — selalu lewat fungsi di
  `@klontong/api` (ditegakkan ESLint).
- Selengkapnya di [`docs/TASKS.md`](docs/TASKS.md) bagian "Aturan global".
