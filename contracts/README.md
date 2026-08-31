# Kontrak API Klontong Net (OpenAPI)

Sumber kebenaran **bentuk API** untuk backend (`apps/backend`) dan frontend
(`apps/frontend`). Ditulis **contract-first**: kontrak lebih dulu, implementasi
menyusul. Frontend tidak perlu menunggu backend — pakai mock server dari kontrak ini.

- Spesifikasi: **OpenAPI 3.0.3**, multi-file, di-`$ref`, dibundel jadi satu file.
- Envelope, error code, paginasi, enum status, dan format uang ditetapkan di
  [`openapi/components/schemas/_common.yaml`](openapi/components/schemas/_common.yaml).
- Konvensi diturunkan dari `apps/backend/internal/platform/httpx`, `.../errs`,
  `.../pagination`, dan `apps/backend/docs/tasks/*.md`.

## Status

| Milestone | Cakupan | Status |
|---|---|---|
| **M1** | Seluruh `/api/v1` buyer: auth, katalog, cart, checkout, order, payment status, review, notification | selesai |
| **M2** | `/api/v1/seller/*`, `/api/v1/admin/*`, `/webhook/*` | selesai |

Kontrak sekarang mencakup **seluruh API**. `pnpm run lint` + `pnpm run bundle` hijau,
Prism mock terbukti jalan.

## Perintah

```bash
cd contracts
make install     # atau: pnpm install
make lint        # Spectral — 0 error sebelum dianggap valid
make bundle      # -> dist/openapi.yaml (dan dist/openapi.json)
make mock        # Prism mock server di http://localhost:4010
make preview     # dokumentasi Redoc di browser
```

Mock dengan contoh statis (default) atau acak:

```bash
pnpm run mock            # pakai `example`/`examples` dari kontrak
pnpm run mock:dynamic    # generate acak dari schema
```

## Cara frontend memakai

1. `make bundle` menghasilkan `dist/openapi.yaml`.
2. Generate tipe TypeScript, mis. `npx openapi-typescript ../contracts/dist/openapi.yaml
   -o src/api/schema.d.ts` (di repo frontend).
3. Untuk dev tanpa backend: `make mock` lalu arahkan `VITE_API_BASE_URL` ke
   `http://localhost:4010`.
4. Jangan mengetik ulang bentuk response — selalu dari tipe hasil generate.

## Struktur

```
openapi/
├── openapi.yaml                 root: info, servers, security, tags, $ref paths
├── components/
│   ├── securitySchemes.yaml
│   ├── parameters.yaml          limit, cursor, Idempotency-Key, path id
│   ├── headers.yaml             X-Request-ID, Retry-After
│   ├── responses.yaml           Error 400/401/403/404/409/422/429/500/502
│   └── schemas/
│       ├── _common.yaml         Money, Meta, Envelope, Error, enum status
│       └── <modul>.yaml         schema per modul
└── paths/
    └── <area>.yaml              operasi per area/kelompok
```

## Aturan yang dikunci

- **Envelope sukses**: `{ "data": <T>, "meta": <Meta> }`.
- **Envelope error**: `{ "error": { "code", "message", "fields"?, "retryable"? } }`.
  `code` = `lower_snake_case`, stabil, spesifik (`out_of_stock`, `price_changed`, ...).
- **Paginasi**: keyset. Query `?limit`, `?cursor` (opaque base64). Response
  `meta: { "next_cursor": string|null, "has_more": boolean }`. Bukan offset.
- **Uang**: `integer` (int64) rupiah. Selalu angka. Persen = basis poin (`integer`,
  250 = 2,5%).
- **Kuantitas**: `number` (boleh pecahan; satuan non-pecahan ditolak server).
- **Header**: `Authorization: Bearer <jwt>` (akses, TTL 15 menit).
  `Idempotency-Key` **wajib** di `POST /api/v1/checkout`. `X-Request-ID` di echo pada
  semua response. `Retry-After` pada 429.
- **ADR-004**: request checkout **tanpa field harga/total**. Klien kirim pilihan;
  `PlaceOrderRequest.client_grand_total` hanya pembanding → memicu `price_changed`.
- **Order** = order induk + array suborder, tiap suborder punya status & nominal
  sendiri. Jangan dikolapskan jadi satu status.

## Cara mengubah kontrak

Kontrak ini **mengikat kedua sisi** — lihat
[ADR-015](../docs/DECISIONS.md). Backend mengimplementasikan apa yang tertulis di
sini; frontend men-generate tipenya dari sini.

Kalau implementasi perlu menyimpang:

1. Ubah `openapi/**` **di PR yang sama** dengan perubahan implementasinya. Jangan
   tinggalkan penyimpangan yang hidup di kode tapi tidak tercermin di kontrak.
2. `make lint` dan `make bundle` wajib hijau (`make ci`).
3. Frontend menjalankan `pnpm gen:api` dan ikut meng-commit `schema.d.ts`.

Penyimpangan diselesaikan **ke arah kontrak**, bukan ke arah yang lebih mudah
diubah. Menambah field aman; menghapus atau mengubah tipe field tidak — itu
memutus tipe yang sudah dipakai frontend.

Perubahan pada "Aturan yang dikunci" di atas butuh ADR baru di
[`docs/DECISIONS.md`](../docs/DECISIONS.md), bukan sekadar edit YAML.

## Catatan ketidaksesuaian scaffold (ditandai di `description`, diseragamkan saat implementasi)

- ~~`order/.../routes.go` pakai `/api/v1/seller/orders/{suborderID}/...`; task doc 12
  menulis `/api/v1/seller/suborders/{id}/...`.~~ **Selesai** — task doc 12 sudah
  disamakan ke `/api/v1/seller/orders/{suborderID}/...`, cocok dengan kontrak dan
  scaffold `routes.go`.
- `GET /api/v1/products/{slug}` vs `GET /api/v1/products/{productId}/reviews` — campur
  slug/id. Didokumentasikan apa adanya.
- Tidak ada webhook kurir; update tracking lewat worker `SyncTracking` (M2 catatan).
