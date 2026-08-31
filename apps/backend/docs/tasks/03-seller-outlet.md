# Fase 03 — Seller + Outlet

> Prasyarat: fase 02. Guide: [`../GUIDES.md`](../GUIDES.md) §2 (tahap 3), §4–§6.
> ADR terkait: ADR-007 (pemilihan outlet), ADR-005 (basis poin komisi).
> Kontrak (**mengikat**, ADR-015): [`paths/seller-shop.yaml`](../../../../contracts/openapi/paths/seller-shop.yaml), [`admin.yaml`](../../../../contracts/openapi/paths/admin.yaml) (bagian sellers).

## Tujuan

Penjual, outlet berkoordinat, pengelola (member), verifikasi. **Selesai kalau** bisa
buka toko dan menambah outlet berkoordinat.

## Aturan khusus fase ini

- `Seller.CommissionBPS` adalah **pointer `*int`**: `nil` = pakai default marketplace;
  membedakan "belum diset" dari "diset 0%".
- Setiap seller butuh **≥ 1 outlet aktif** untuk bisa jualan. Outlet **wajib
  koordinat** yang masuk akal.
- Selalu **≥ 1 owner** — cek `CountByRole` di dalam tx yang sama saat remove/ubah peran.
- Ganti akun payout = **owner saja** (titik penyalahgunaan).
- Peran token `seller` cuma berarti "punya toko", bukan "boleh mengubah toko ini" —
  `authz.go` cek keanggotaan.
- Area HTTP: `/api/v1/seller/*` (login + `RequireRole("seller")` + keanggotaan) dan
  `/api/v1/admin/*` (verifikasi).
- Registrasi seller **menambah peran** `seller` ke akun yang ada via `identity.Port.
  GrantRole` — bukan bikin akun kedua.

## Urutan kerja

### Domain
- [ ] `internal/modules/seller/internal/domain/seller.go` — `Seller{OwnerUserID,Slug,
      Status,CommissionBPS *int,PayoutBank...,VerifiedAt}`; `NewSeller` (slug = nama +
      suffix acak); `Verify` (hanya dari pending), `Suspend`, `CanSell` (hanya verified),
      `SetPayoutAccount`.
- [ ] `.../domain/outlet.go` — `Outlet{lat,lng wajib & sane,IsActive,SupportsPickup/
      LocalDelivery/Courier,OperatingHours}`; `SupportsMethod`, `IsOpenAt`. Outlet tutup
      tetap terima order kurir.
- [ ] `.../domain/member.go` — `MemberRole` = owner|manager|staff; `CanManageProducts`,
      `CanManagePayout` (owner saja), `CanProcessOrders`.
- [ ] `.../domain/document.go` — `Document{Kind,StorageKey,Status,ReviewedBy,
      RejectionReason}` (KTP/NPWP), tidak pernah publik.
- [ ] `.../domain/repository.go` — `SellerRepository` (.../FindBySlug/ListForAdmin/
      UpdateStatus), `OutletRepository` (.../ListActiveBySeller/CountActive),
      `MemberRepository` (.../CountByRole), `DocumentRepository`.
- [ ] `.../domain/errors.go` — `ErrSellerNotVerified/Suspended`, `ErrOutletInactive`,
      `ErrNotSellerMember`, `ErrRoleNotPermitted`, `ErrNoActiveOutlet`, `ErrSlugTaken`,
      `ErrCannotRemoveLastOwner`.

### Migrasi
- [ ] `migrations/000003_create_seller.up.sql` / `.down.sql` — `seller_sellers`
      (status, `commission_bps` nullable, payout bank), `seller_outlets` (lat/lng),
      `seller_members`, `seller_documents` sesuai spec; unique index slug.

### Infra
- [ ] `.../infra/{seller_repository,outlet_repository,member_repository,
      document_repository,mapper}.go`.

### App
- [ ] `.../app/authz.go` — `requireSellerMember`, `requireSellerRole(min)`.
- [ ] `.../app/seller_usecase.go` — `RegisterSeller`: `WithinTx` buat seller (pending) +
      tambah owner sebagai member + `identity.Port.GrantRole("seller")` + outbox
      `EventSellerRegistered`. Plus `GetSeller/GetPublicSeller/UpdateSeller/
      SetPayoutAccount(owner)/ListMySellers`.
- [ ] `.../app/outlet_usecase.go` — `Create/Update/List/DeactivateOutlet`; nonaktifkan
      outlet aktif terakhir yang punya order berjalan → tolak / peringatan dengan jumlah
      produk terdampak.
- [ ] `.../app/member_usecase.go` — `List/Invite/Remove/ChangeMemberRole`; cek
      `CountByRole` di tx yang sama (≥1 owner).
- [ ] `.../app/verification_usecase.go` — `UploadDocument`, `ListDocuments`,
      `ApproveSeller` (`WithinTx` ubah status + outbox `EventSellerVerified`),
      `RejectSeller` (reason wajib).
- [ ] `.../app/{service,dto,port_adapter}.go`.

### Transport
- [ ] `.../transport/rest/{routes,request,response,seller_handler,outlet_handler,
      member_handler,admin_handler,handler}.go` — pisahkan area buyer / seller / admin.

### Port
- [ ] `internal/modules/seller/port.go` — `Get`, `GetMany` (batch), `CanSell`,
      `CommissionBPS` (kembalikan efektif: seller atau default), `GetOutlet`,
      `ListOutlets`, `IsMember`. `Outlet` bawa lat/lng + `Supports*` + `IsActive`.
- [ ] `events.go` — `seller.seller.registered/verified/suspended`,
      `seller.outlet.created/deactivated`.
- [ ] `module.go` — `Config{...,Users identity.Port,Storage,DefaultCommissionBPS}`.

### Wiring
- [ ] `internal/app/registry.go` — build `seller` setelah `identity` (terima
      `identityMod.Port()`).

## Test wajib

- `RegisterSeller` menambah peran ke akun lama, tidak membuat user baru.
- Remove owner terakhir → `ErrCannotRemoveLastOwner`.
- `CommissionBPS` port: seller `nil` → kembalikan default marketplace; seller `0` →
  kembalikan 0.
- Buat outlet tanpa koordinat → ditolak validasi domain.

## Sengaja TIDAK dikerjakan di fase ini

- Produk / varian (fase 4).
- Konsumen `EventSellerSuspended` di catalog/cart (fase 4 / 7).
- Payout / buku besar (fase 14).
