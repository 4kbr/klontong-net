# Fase 09 — Promotion

> Prasyarat: fase 03. Guide: [`../GUIDES.md`](../GUIDES.md) §2 (tahap 9), §10.
> ADR terkait: ADR-005 (basis poin, `money.Distribute`), ADR-008 (idempotensi kuota).
> **Fase menyentuh uang — tulis test lebih dulu.**
> Kontrak (**mengikat**, ADR-015): [`paths/vouchers.yaml`](../../../../contracts/openapi/paths/vouchers.yaml), [`seller-promotion.yaml`](../../../../contracts/openapi/paths/seller-promotion.yaml), [`admin.yaml`](../../../../contracts/openapi/paths/admin.yaml) (bagian vouchers).

## Tujuan

Voucher + pembagian diskon ke suborder. **Selesai kalau** voucher terbagi ke suborder
dan jumlahnya persis sama dengan diskon.

## Aturan khusus fase ini

- Persentase dalam **basis poin** `int` (12,5% = 1250), bukan float.
  `CalculateDiscount`: percentage = `base * ValueBPS / 10000`, di-cap `MaxDiscountAmount`,
  dibulatkan **ke bawah**; fixed = `ValueAmount` tapi tidak melebihi base.
  `MaxDiscountAmount` **wajib** untuk voucher percentage.
- Voucher **seller** hanya berlaku untuk suborder milik seller itu — bug marketplace
  paling klasik adalah voucher toko A mendiskon barang toko B.
- Split diskon ke suborder pakai **`money.Distribute`** — Σ bagian = diskon persis, sisa
  pembulatan dialokasikan deterministik, **tidak dibuang**.
- `OwnerType` (marketplace vs seller) menentukan **siapa menanggung** diskon — salah →
  pembukuan settlement salah. Jangan percaya `OwnerType` dari body request; tentukan dari
  peran pembuat.
- Kuota: **unique index** `(voucher_id, order_id)` **DAN** cek `UsedCount` di dalam tx
  dengan `SELECT ... FOR UPDATE` pada baris voucher. Kunci idempotensi dibuat sekali.
- `Validate` **tidak mengubah apa pun** — dipanggil berulang saat pembeli coba-coba kode.

## Urutan kerja

### Domain
- [ ] `internal/modules/promotion/internal/domain/voucher.go` — `Kind` = percentage|
      fixed_amount|free_shipping; `OwnerType` = marketplace|seller; `AppliesTo` = items|
      shipping; `Voucher{Code,OwnerType,OwnerSellerID *uuid,Kind,ValueBPS int,ValueAmount,
      MaxDiscountAmount,MinOrderAmount,AppliesTo,Quota,QuotaPerUser,UsedCount,StartsAt/
      EndsAt,IsActive}`; `IsUsableAt`, `HasQuota`, `CalculateDiscount(base)`.
- [ ] `.../domain/rule.go` — `ApplyContext{UserID,Subtotals map[uuid]money(per suborder),
      ShippingAmounts map[uuid]money,SellerIDs,At}`; `Voucher.EligibleFor(c) -> (basis,
      err)`. Aturan: voucher seller → hanya suborder seller itu; `MinOrderAmount` dihitung
      dari basis relevan; free_shipping hanya kurangi ongkir, tidak lebih.
- [ ] `.../domain/redemption.go` — `Redemption{VoucherID,UserID,OrderID,SuborderID,
      DiscountAmount}`.
- [ ] `.../domain/repository.go` — `VoucherRepository{...,FindByCode,FindByCodeForUpdate
      (SELECT..FOR UPDATE),ListActive,IncrementUsed,DecrementUsed}`,
      `RedemptionRepository{Create,CountByVoucherAndUser,ListByOrder,DeleteByOrder}`.
- [ ] `.../domain/errors.go` — `ErrVoucherExpired/NotStarted/Inactive`,
      `ErrQuotaExhausted`, `ErrUserQuotaExhausted`, `ErrMinOrderNotMet` (sebut
      kekurangannya), `ErrVoucherNotApplicable`.

### Migrasi
- [ ] `migrations/000011_create_promotion.up.sql` / `.down.sql` — `promotion_vouchers`
      (`owner_type`, `kind`, `value_bps` int, `value_amount` bigint, kuota),
      `promotion_redemptions` sesuai spec; **unique index `(voucher_id, order_id)`**.

### Infra
- [ ] `.../infra/{voucher_repository,redemption_repository,mapper}.go`.

### App
- [ ] `.../app/apply_usecase.go` — `Validate(code, ApplyContext) -> Applied`: (1) cari
      voucher, cek aktif/periode/kuota global+per-user, (2) tentukan basis (voucher
      seller → subtotal seller itu saja), (3) cek `MinOrderAmount`, (4) hitung diskon,
      (5) **split ke suborder eligible dengan `money.Distribute`** (Σ = diskon persis).
      Tidak mengubah apa pun.
- [ ] `.../app/redeem_usecase.go` — `Redeem(voucherID,userID,orderID,perSuborder)` (di
      dalam tx checkout): `FindByCodeForUpdate` kunci baris → cek ulang kuota →
      `IncrementUsed` → buat `Redemption` per suborder. `ReleaseRedemption(orderID)` (saat
      batal/expiry): `DecrementUsed` + hapus redemption; **idempoten**.
- [ ] `.../app/voucher_usecase.go` — `Create/Update/List/DeactivateVoucher`. Seller hanya
      untuk tokonya; admin untuk voucher marketplace — cek di sini.
- [ ] `.../app/{service,dto,port_adapter}.go`.

### Transport
- [ ] `.../transport/rest/{routes,request,response,voucher_handler,handler}.go` — kelola
      voucher di `/api/v1/seller` + `/api/v1/admin`; `Validate` di `/api/v1`.

### Port
- [ ] `internal/modules/promotion/port.go` — `Validate(code, ApplyContext) -> Applied`
      (`Applied.PerSuborder` sudah di-split), `Redeem(voucherID,userID,orderID,
      perSuborder)`, `ReleaseRedemption(orderID)`.
- [ ] `events.go` — `EventVoucherRedeemed`, `EventVoucherExhausted`,
      `EventVoucherExpired`.
- [ ] `module.go` — `Config{...,Sellers seller.Port}`.

### Wiring
- [ ] `internal/app/registry.go` — build `promotion` setelah `seller`.

## Test wajib

- **Split diskon**: voucher Rp10.000 untuk keranjang 3 suborder → 3 bagian, Σ = 10.000
  persis, sisa 1 rupiah jatuh deterministik (property test).
- Voucher seller A tidak mendiskon suborder seller B.
- **Idempotensi kuota**: `Redeem` kunci sama dua kali → `UsedCount` naik sekali (unique
  index yang menjaga).
- `ReleaseRedemption` dijalankan dua kali → `UsedCount` turun sekali.
- `CalculateDiscount` percentage dibulatkan ke bawah dan di-cap `MaxDiscountAmount`.

## Sengaja TIDAK dikerjakan di fase ini

- Pemanggilan `Validate`/`Redeem`/`ReleaseRedemption` dari checkout & expiry (fase 10 /
  12).
- Efek `OwnerType` di buku besar (fase 14).
- Konsumen `EventVoucherRedeemed` di audit (fase 13).
