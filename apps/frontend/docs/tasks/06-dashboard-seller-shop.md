# Fase 06 — Shell dashboard + Seller shop

> Prasyarat: 01, 02 | Ref: kontrak `paths/seller-shop.yaml`;
> `docs/ARCHITECTURE.md` §5 (area `/api/v1/seller/*`, peran + keanggotaan toko),
> ADR-007 (outlet berkoordinat).

## Tujuan

Kerangka aplikasi `dashboard` dan modul toko penjual: onboarding, profil toko,
akun pencairan, dokumen verifikasi, outlet, dan pengelola. **Selesai kalau**
seorang seller bisa mendaftar toko, melengkapi profil + rekening pencairan,
mengunggah dokumen, serta menambah/mengubah outlet berkoordinat dan anggota.

## Aturan khusus fase ini

- Semua rute di balik `RequireRole('seller')`. Peran hanya gerbang kasar —
  tampilkan/aktifkan aksi sesuai **keanggotaan toko** & status verifikasi dari
  `seller/me`; jangan andalkan peran saja.
- Toko `pending`/`rejected`/`suspended` → batasi aksi, tampilkan banner status +
  alasan (dari admin, Fase 09).
- Outlet **wajib koordinat** (lat/long) — dipakai pemilihan outlet & ongkir.
  Sediakan input map/geocode sederhana; validasi rentang koordinat.
- `payout-account` menyimpan data rekening — perlakukan sebagai sensitif, jangan
  log, mask di tampilan.
- Upload dokumen: ikuti pola `upload-url` bila kontrak memberi presigned URL
  (lihat juga Fase 07 untuk gambar produk).

## Urutan kerja

### Tipe & API (packages/api)
- [ ] `src/endpoints/seller-shop.api.ts` — `registerShop` (`/seller/register`),
      `getMyShop`/`updateMyShop` (`/seller/me`), `getPayoutAccount`/`updatePayoutAccount`,
      `listDocuments`/`uploadDocument`, `listOutlets`/`createOutlet`,
      `getOutlet`/`updateOutlet`/`deleteOutlet`, `listMembers`/`inviteMember`,
      `getMember`/`updateMemberRole`/`removeMember`.
- [ ] `src/hooks/` — `useMyShop`, `usePayoutAccount`, `useDocuments`, `useOutlets`, `useOutletMutations`, `useMembers`, `useMemberMutations`.
- [ ] `src/schemas/seller-shop.ts` — Zod: registrasi toko, profil, rekening (nomor/nama/bank), outlet (nama, alamat, `latitude`, `longitude`, metode kirim didukung), undang anggota.

### State (Zustand / Query)
- [ ] `dashboardStore` — konteks toko aktif (bila multi-toko), status verifikasi cache.
- [ ] `queryKeys.seller.*` namespace.

### Komponen (packages/ui / dashboard)
- [ ] `DashboardShell` (sidebar nav seller+admin, topbar, `NotificationBell`),
      `VerificationBanner`, `CoordinateInput`/`MiniMap`, `DocumentUploader`,
      `OutletForm`, `MemberTable` + `InviteMemberDialog`, `MaskedField`.

### Halaman & rute (dashboard, `RequireRole('seller')`)
- [ ] `/seller/onboarding` — `registerShop` (untuk user berperan buyer yang ingin jadi seller: arahkan alur peran sesuai kontrak).
- [ ] `/seller` — ringkasan + `VerificationBanner`.
- [ ] `/seller/toko` — profil toko + `payout-account` + dokumen.
- [ ] `/seller/outlet` — daftar + tambah/edit + hapus (cegah hapus outlet dengan stok/riwayat sesuai respons).
- [ ] `/seller/anggota` — daftar + undang + ubah peran + hapus.

### Wiring
- [ ] `dashboard` app: `router.tsx` route group `(seller)` + `(admin)` (admin diisi Fase 09), layout `DashboardShell`.
- [ ] MSW handlers `seller-shop` stateful; fixture: 1 toko `active`, 1 toko `pending`, 2 outlet berkoordinat, 2 anggota.
- [ ] Login dashboard (dari Fase 02) → deteksi peran → arahkan seller ke `/seller`, admin ke `/admin`.

## Test wajib

- User tanpa peran seller membuka `/seller/*` → ditolak (halaman akses ditolak), bukan crash.
- Toko `pending` → `VerificationBanner` tampil, aksi tertentu (mis. publish produk nanti) dinonaktifkan dengan alasan.
- Outlet tanpa koordinat valid → form menolak submit; dengan koordinat → muncul di daftar.
- `payout-account`: nomor rekening ter-mask di tampilan; tidak muncul utuh di DOM state yang ter-serialize.
- Undang anggota → muncul di `MemberTable`; ubah peran → tercermin; hapus diri sendiri sebagai pemilik → dicegah (ikuti respons server).

## Sengaja TIDAK dikerjakan

- Produk, varian, harga, stok — Fase 07.
- Suborder, shipment, keuangan — Fase 08.
- Panel admin — Fase 09.
