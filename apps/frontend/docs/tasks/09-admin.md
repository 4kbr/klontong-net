# Fase 09 — Panel admin

> Prasyarat: 06 | Ref: kontrak `paths/admin.yaml`; **ADR-002** (refund parsial),
> **ADR-008** (rekonsiliasi), **ADR-011** (buku besar), ADR-012 (COD).

## Tujuan

Panel marketplace: verifikasi toko, kategori, voucher marketplace, pengawasan
pesanan/pembayaran, refund, settlement & payout, laporan rekonsiliasi, dan audit.
**Selesai kalau** admin bisa menjalankan seluruh endpoint `/api/v1/admin/*`
terhadap mock: menyetujui/menolak toko, mengelola kategori & voucher, memicu
rekonsiliasi/refund/retry payout, dan membaca laporan + audit.

## Aturan khusus fase ini

- Semua rute di balik `RequireRole('admin')`.
- Aksi berisiko (suspend toko, cancel order, refund, retry payout) → dialog
  konfirmasi + alasan wajib + tampilkan dampak ("membatalkan order ini melepas
  stok & memicu refund").
- **Refund hampir selalu parsial** (per suborder) — form refund berbasis suborder/
  item, bukan "refund seluruh order".
- **Rekonsiliasi**: laporan harian membandingkan pembayaran masuk, pendapatan
  penjual, komisi, pencairan. **Selisih harus nol** — sorot merah bila tidak,
  jangan sembunyikan.
- **Buku besar**: tampilkan angka apa adanya dari entri; jangan ada kolom saldo
  yang "dihitung ulang" di klien.
- Audit bersifat **read-only** dan jadi bukti — jangan sediakan edit/hapus.
- Kategori dimiliki admin (dipakai seller di Fase 07 sebagai read-only).

## Urutan kerja

### Tipe & API (packages/api)
- [ ] `src/endpoints/admin.api.ts` —
      sellers: `list` (`/admin/sellers`), `approve`, `reject`, `suspend`;
      `listCategories`/`createCategory`/... (`/admin/categories`);
      vouchers: `list`/`create`/... (`/admin/vouchers`);
      orders: `list` (`/admin/orders`), `cancel` (`/admin/orders/{id}/cancel`);
      payments: `list` (`/admin/payments`), `reconcile` (`/admin/payments/{id}/reconcile`);
      `listRefunds`/`createRefund` (`/admin/refunds`);
      `listSettlements` (`/admin/settlements`);
      payouts: `list` (`/admin/payouts`), `retry` (`/admin/payouts/{id}/retry`);
      `getReconciliationReport` (`/admin/reports/reconciliation`);
      audit: `list` (`/admin/audit`), `byTarget` (`/admin/audit/{targetType}/{targetId}`).
- [ ] `src/hooks/` — satu hook per resource (`useAdminSellers`, `useAdminCategories`,
      `useAdminVouchers`, `useAdminOrders`, `useAdminPayments`, `useAdminRefunds`,
      `useAdminSettlements`, `useAdminPayouts`, `useReconciliationReport`, `useAudit`)
      + hook mutasi terkait, semua keyset di list.
- [ ] `src/schemas/admin.ts` — Zod: alasan approve/reject/suspend, kategori
      (nama, parent, slug), voucher marketplace (tipe, nominal/bps, kuota, periode,
      `owner_type`), refund (target suborder/item + nominal + alasan), cancel order.

### State (Zustand / Query)
- [ ] `queryKeys.admin.*` namespace.
- [ ] Invalidasi: approve/reject/suspend → invalidate `adminSellers` (+ detail);
      reconcile/refund/retry → invalidate resource terkait + `reconciliationReport`.

### Komponen (packages/ui / dashboard)
- [ ] `AdminTable` (generik: kolom, filter, keyset, aksi baris), `SellerReviewPanel`
      (dokumen + tombol approve/reject/suspend + alasan), `CategoryTreeEditor`,
      `MarketplaceVoucherForm`, `OrderInspector` (induk + suborder, read-only + cancel),
      `PaymentInspector` + `ReconcileButton`, `RefundForm` (berbasis suborder/item),
      `SettlementTable`, `PayoutTable` + `RetryButton`, `ReconciliationReport`
      (baris + total + badge "selisih nol / TIDAK nol"), `AuditLogTable`.

### Halaman & rute (dashboard, `RequireRole('admin')`)
- [ ] `/admin` — ringkasan (toko menunggu verifikasi, payout gagal, selisih rekonsiliasi).
- [ ] `/admin/toko` + `/admin/toko/:id` — verifikasi.
- [ ] `/admin/kategori`, `/admin/voucher`.
- [ ] `/admin/pesanan` + `/admin/pesanan/:id`.
- [ ] `/admin/pembayaran` + `/admin/refund`.
- [ ] `/admin/settlement` + `/admin/payout`.
- [ ] `/admin/laporan/rekonsiliasi`.
- [ ] `/admin/audit` + `/admin/audit/:targetType/:targetId`.

### Wiring
- [ ] MSW handlers stateful `admin`: approve toko → status toko berubah (terlihat di dashboard seller Fase 06);
      reconcile/refund/retry mengubah state pembayaran/payout; `reconciliation` menghitung dari `db` (fixture selisih nol + skenario tak nol).
- [ ] Fixture: 2 toko pending, kategori tree, 1 voucher marketplace, beberapa payment/settlement/payout, entri audit.

## Test wajib

- Non-admin membuka `/admin/*` → ditolak.
- Approve toko → status berubah; membuka dashboard seller (fixture user seller yang sama) menampilkan banner `active`.
- `suspend` tanpa alasan → ditolak; dengan alasan → tercatat & tampil di audit.
- `RefundForm` menargetkan satu suborder → refund parsial; order induk & suborder lain tidak terpengaruh.
- `ReconciliationReport`: fixture seimbang → badge "selisih nol"; fixture timpang → badge merah + baris penyebab disorot.
- `retryPayout` pada payout `failed` → status berubah; daftar ter-invalidate.
- `AuditLogTable` tidak punya aksi edit/hapus.

## Sengaja TIDAK dikerjakan

- Perubahan pada permukaan buyer/seller di luar yang dipicu aksi admin.
- E2E & hardening — Fase 10.
