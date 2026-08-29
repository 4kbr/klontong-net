# Fase 07 — Seller katalog & inventori

> Prasyarat: 06 | Ref: kontrak `paths/seller-catalog.yaml`,
> `paths/seller-inventory.yaml`; **ADR-005** (uang int64 / bps),
> **ADR-006** (satuan terkecil), ADR-003 (on_hand vs reserved).

## Tujuan

Penjual mengelola produk, varian, gambar, harga + tier grosir, penerbitan, serta
stok per outlet (mutasi, transfer, stok opname). **Selesai kalau** penjual bisa
membuat produk dengan varian dus/renceng/pcs, menetapkan harga per satuan + tier,
menerbitkan, dan menyesuaikan stok per outlet dengan jejak mutasi.

## Aturan khusus fase ini

- **ADR-006.** Varian menyimpan `unit_code`, `content_quantity` (> 0; = 1 bila
  `unit_code == base_unit_code`), `base_unit_code`. Form menjelaskan
  "1 dus = 40 pcs"; stok ditampilkan/diedit dalam **satuan dasar** (atau ikuti
  bentuk kontrak) — beri konversi jelas di UI.
- **ADR-005.** Semua input harga = **integer rupiah**; komponen input menolak
  desimal/format. Diskon/persen tier = **basis poin** (tampilkan sebagai %,
  simpan bps). Tidak ada `float` di jalur harga.
- Harga & tier lewat endpoint terpisah (`/seller/variants/{id}/price`,
  `/seller/variants/{id}/tiers`, `/seller/prices` batch). Tier tidak boleh
  tumpang tindih; validasi `min_qty` menaik.
- Publikasi (`/seller/products/{id}/publish`) butuh kelengkapan (varian, harga,
  minimal 1 gambar, toko terverifikasi) — cek & tampilkan checklist sebelum
  tombol aktif.
- Gambar: `POST /seller/products/{id}/images/upload-url` → PUT ke storage →
  simpan referensi. Tangani kegagalan upload tanpa kehilangan draft.
- Stok: **setiap perubahan `on_hand` menghasilkan Movement** (backend); UI tampil
  `on_hand`, `reserved`, `tersedia = on_hand - reserved` terpisah. Jangan hitung
  ulang reserved sendiri.
- Stok opname: alur draft → isi hitungan fisik → `finish` (selisih jadi Movement).

## Urutan kerja

### Tipe & API (packages/api)
- [ ] `src/endpoints/seller-catalog.api.ts` — produk (`list/create/get/update`),
      `publishProduct`, varian (`list/create` `/products/{id}/variants`,
      `get/update` `/variants/{id}`), `setVariantPrice`, `setVariantTiers`,
      `batchPrices` (`/seller/prices`), `createImageUploadUrl`.
- [ ] `src/endpoints/seller-inventory.api.ts` — `getOutletStock`
      (`/seller/outlets/{id}/stock`), `getOutletStockItem`
      (`/seller/outlets/{id}/stock/{variantId}` — set/adjust), `transferStock`
      (`/seller/stock/transfer`), `listMovements` (`/seller/stock/movements`),
      `listOpnames`/`createOpname` (`/seller/outlets/{id}/opnames`),
      `finishOpname` (`/seller/opnames/{id}/finish`).
- [ ] `src/hooks/` — `useSellerProducts`, `useProductMutations`, `useVariantMutations`,
      `usePriceMutations`, `useOutletStock`, `useStockMutations`, `useMovements` (keyset),
      `useOpnames`.
- [ ] `src/schemas/seller-catalog.ts`, `src/schemas/seller-inventory.ts` — Zod:
      produk, varian (satuan + content_quantity), harga (int, > 0), tier (array
      `min_qty` menaik, harga/bps), penyesuaian stok, baris opname.

### State (Zustand / Query)
- [ ] Draft produk multi-langkah di `sessionStorage` (tahan reload saat upload gambar).
- [ ] Invalidasi: mutasi harga/varian → invalidate produk & varian terkait; mutasi stok → invalidate `outletStock` + `movements`.

### Komponen (packages/ui / dashboard)
- [ ] `ProductForm` (wizard: info → varian → gambar → harga → tinjau),
      `VariantEditor` (satuan, content_quantity, konversi preview),
      `RupiahInput` (integer-only), `TierTable` (validasi tumpang tindih),
      `PublishChecklist`, `ImageUploader` (upload-url + progress + retry),
      `StockTable` (on_hand / reserved / tersedia), `StockAdjustDialog`,
      `TransferForm`, `MovementList`, `OpnameSheet`.

### Halaman & rute (dashboard, `RequireRole('seller')`)
- [ ] `/seller/produk` — daftar + status (draft/published) + filter.
- [ ] `/seller/produk/baru` & `/seller/produk/:id` — `ProductForm`.
- [ ] `/seller/produk/:id/harga` — harga & tier per varian (atau inline di form).
- [ ] `/seller/stok` — pilih outlet → `StockTable` + adjust + transfer.
- [ ] `/seller/stok/mutasi` — `MovementList` (keyset, filter varian/tanggal).
- [ ] `/seller/stok/opname` — daftar + buat + isi + `finish`.

### Wiring
- [ ] MSW handlers stateful: `seller-catalog` (produk/varian/harga/tier/publish/upload-url dummy),
      `seller-inventory` (stok per outlet, movement tiap perubahan, transfer, opname).
- [ ] Fixture: 1 produk 3 varian (pcs/renceng/dus), stok di 2 outlet, beberapa movement.

## Test wajib

- `RupiahInput` menolak `12.500,50` dan input non-digit; menyimpan `12500`.
- Varian dus dengan `content_quantity=40` → preview "1 dus = 40 pcs"; `content_quantity` wajib > 0.
- `TierTable`: tier dengan `min_qty` tumpang tindih / tidak menaik → validasi gagal sebelum submit.
- `PublishChecklist`: produk tanpa gambar/harga → tombol publish nonaktif dengan alasan; lengkap → publish sukses.
- Upload gambar gagal di tengah → draft produk tetap ada, bisa retry.
- Adjust stok +10 → `on_hand` naik, muncul entri di `MovementList`; `tersedia` = `on_hand - reserved`.
- Transfer antar outlet → dua movement (keluar/masuk), total `on_hand` toko tetap.
- `finishOpname` dengan selisih → movement penyesuaian tercatat.

## Sengaja TIDAK dikerjakan

- Suborder / pengiriman / keuangan seller — Fase 08.
- Voucher toko — Fase 08.
- Kategori (dimiliki admin) — Fase 09.
