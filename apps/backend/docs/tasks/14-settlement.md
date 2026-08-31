# Fase 14 — Settlement

> Prasyarat: fase 10, 11, 13. Guide: [`../GUIDES.md`](../GUIDES.md) §2 (tahap 14), §8,
> §9. ADR: 011 (double entry, bukan kolom saldo), 012 (COD arah terbalik), 008
> (idempotensi + rekonsiliasi), 013 (kunci terdistribusi worker uang).
> **"Bagian yang paling tidak boleh salah." Tulis test lebih dulu.**
> Kontrak (**mengikat**, ADR-015): [`paths/seller-finance.yaml`](../../../../contracts/openapi/paths/seller-finance.yaml), [`admin.yaml`](../../../../contracts/openapi/paths/admin.yaml) (settlements, payouts, reports).

## Tujuan

Komisi, buku besar double entry, pematangan dana, pencairan ke penjual, laporan
rekonsiliasi. **Selesai kalau** `ReconciliationReport` selisihnya **nol**.

## Aturan khusus fase ini

- **Tidak ada kolom saldo yang di-UPDATE** (ADR-011). Saldo = penjumlahan entri. Kalau
  lambat → snapshot berkala yang tetap bisa dihitung ulang.
- Setiap jurnal **wajib seimbang**: Σ debit = Σ kredit. Validasi sebelum simpan;
  `ErrUnbalancedJournal` **tidak boleh muncul di production** — kalau muncul, itu bug
  kita, **gagalkan transaksi**.
- Semua handler event **idempoten**: unique index `event_id` pada tabel jurnal,
  `INSERT ... ON CONFLICT DO NOTHING`. Diproses dua kali = penjual dibayar dua kali.
- `Payout.IdempotencyKey` **DB-unique**, dibuat **sekali** saat payout dibuat (bukan per
  percobaan kirim). Idempotency key ke provider = pertahanan **terakhir**, bukan pertama.
- Worker uang (`MatureEarnings`, `ProcessPayouts`) pakai **kunci terdistribusi Redis**
  supaya hanya satu instance yang jalan.
- Refund setelah dana cair → **saldo negatif / piutang**, bukan error. Jangan pernah
  menolak refund karena uang sudah keluar.
- COD (ADR-012): uang tidak pernah lewat kita di awal → `courier/cod_receivable` sampai
  setoran kurir masuk. Alur `settlement` **bercabang** per metode pembayaran — tulis
  jurnalnya eksplisit dan berbeda.
- Ledger usecase dipanggil **dari event handler**, bukan dari HTTP.
- Pertimbangkan cabut GRANT `UPDATE`/`DELETE` pada tabel entri di level role DB.

## Urutan kerja

### Domain
- [ ] `internal/modules/settlement/internal/domain/entry.go` — `Direction` = debit|credit;
      `Entry{JournalID,AccountID,Direction,Amount,ReferenceType,ReferenceID,Description,
      OccurredAt}`; `Journal{ID,Entries,Description}`, `IsBalanced()`, `Validate()`.
      Append-only.
- [ ] `.../domain/account.go` — `OwnerType` = seller|marketplace|buyer|gateway|courier;
      `AccountKind` = pending|available|payable|revenue|clearing;
      `Account{OwnerType,OwnerID *uuid,Kind,Currency}`. Akun minimal:
      `seller/<id>/pending`, `seller/<id>/available`, `marketplace/revenue`,
      `marketplace/clearing`, `courier/cod_receivable`.
- [ ] `.../domain/earning.go` — `Earning{SuborderID,SellerID,GrossAmount,
      CommissionAmount,NetAmount,Status(pending|matured|paid|reversed),MaturesAt}`;
      `NewEarning(...,holdPeriod,now)`, `Mature(now)`, `Reverse(reason)`.
      `MaturesAt = waktu suborder completed + hold period`.
- [ ] `.../domain/payout.go` — `PayoutStatus` = requested|processing|paid|failed;
      `Payout{SellerID,Amount,Status,Bank*,ProviderReference,IdempotencyKey,timestamps,
      FailedReason,Items []PayoutItem}`, `PayoutItem{PayoutID,SuborderID,Amount}`;
      `MarkPaid(ref,now)`, `MarkFailed(reason)`.
- [ ] `.../domain/disburser.go` — `Disburser interface{Transfer(TransferRequest) ->
      (reference, err),Inquiry(reference) -> (status, err),ValidateAccount(bankCode,
      accountNumber) -> (accountName, err)}`. `ValidateAccount` saat seller set rekening,
      bukan saat payout. `Transfer` bawa idempotency key; timeout ≠ gagal (tak diketahui).
- [ ] `.../domain/repository.go` — `AccountRepository{EnsureExists,FindByOwner,Balance
      (accountID)}`, `EntryRepository{CreateJournal(*Journal),ListByAccount,SumByAccount}`,
      `EarningRepository{Create,FindBySuborder,ListMatured(before,limit),UpdateStatus,
      ListBySeller}`, `PayoutRepository{Create,Update,FindByID,ListPending,ListBySeller,
      ExistsInProgress}`.
- [ ] `.../domain/errors.go` — `ErrAccountNotFound`, `ErrUnbalancedJournal`,
      `ErrInsufficientBalance`, `ErrBelowMinimumPayout`, `ErrPayoutInProgress`,
      `ErrNoPayoutAccount`, `ErrEarningNotMatured`, `ErrDisburserUnavailable`.

### Migrasi
- [ ] `migrations/000012_create_settlement.up.sql` / `.down.sql` — `settlement_accounts`
      (`owner_type`, `kind`), `settlement_entries` (append-only, amount `bigint`,
      `journal_id`), `settlement_earnings`, `settlement_payouts` + items sesuai spec.
      Unique `event_id` di tabel jurnal; unique `payout.idempotency_key`. Pertimbangkan
      deferred constraint keseimbangan.

### Infra
- [ ] `.../infra/{account_repository,entry_repository,earning_repository,
      payout_repository,mapper}.go` — `Balance`/`SumByAccount` = agregasi entri, bukan
      kolom.
- [ ] `.../infra/disburser.go` — impl `Disburser` (REST langsung) + versi palsu untuk dev.

### App
- [ ] `.../app/ledger_usecase.go` — tiap fungsi membuat **satu jurnal seimbang**,
      divalidasi sebelum simpan:
  - `RecordOrderPaid(orderID)` — Debit `marketplace/clearing` (total); Credit
    `seller/<id>/pending` (earning tiap seller); Credit `marketplace/revenue` (komisi).
  - `RecordEarningMatured(earningID)` — Debit `seller/pending` → Credit `seller/available`.
  - `RecordPayout(payoutID)` — Debit `seller/available` → Credit `marketplace/clearing`.
  - `RecordRefund(refundID)` — arah dibalik; sudah cair → saldo negatif / piutang,
    jangan tolak.
  - `RecordCODCollected(orderID)` — alur COD berbeda (ADR-012).
- [ ] `.../app/earning_usecase.go` — `MatureEarnings()` worker: `ListMatured` +
      `FOR UPDATE SKIP LOCKED`; per baris `WithinTx` `earning.Mature` +
      `RecordEarningMatured`; idempoten (skip yang sudah matured). Plus
      `ListSellerEarnings`, `GetSellerBalance`.
- [ ] `.../app/payout_usecase.go` — `RequestPayout(sellerID,amount)`: (1) rekening diset
      & tervalidasi, (2) saldo available cukup (dihitung dari ledger, di dalam tx), (3)
      `ExistsInProgress` false, (4) di atas `PayoutMinimum`, (5) buat Payout `requested`
      + `IdempotencyKey`, (6) `RecordPayout`. `ProcessPayouts()` worker: kunci Redis; per
      payout `Disburser.Transfer` dengan idempotency key payout itu; sukses → `MarkPaid`;
      gagal → `MarkFailed` + **reverse jurnal**. `HandleDisburseCallback` /
      `ReconcilePayouts` via `Inquiry`.
- [ ] `.../app/report_usecase.go` — `SellerEarningReport`, `MarketplaceRevenueReport`,
      `ReconciliationReport` (bandingkan Σ pembayaran masuk vs Σ earning seller vs Σ
      komisi vs Σ payout — **selisih harus nol**; jalan harian, alarm bila tidak).
- [ ] `.../app/handlers.go` — semua idempoten: `OnOrderPaid` → buat Earning per suborder
      + `RecordOrderPaid`; `OnSuborderCompleted` → set `MaturesAt`; `OnSuborderRejected`
      → batalkan earning suborder itu; `OnRefundCompleted` → `RecordRefund` + reverse
      earning bila perlu; `OnCODCollected` → `RecordCODCollected`.
- [ ] `.../app/{service,dto,port_adapter}.go`.

### Transport
- [ ] `.../transport/rest/{routes,request,response,balance_handler,payout_handler,
      admin_handler,handler}.go` — saldo + request payout di `/api/v1/seller`; laporan
      rekonsiliasi di `/api/v1/admin`.

### Port
- [ ] `internal/modules/settlement/port.go` — **hanya** `SellerBalance(sellerID) ->
      Balance{Pending,Available}`. Tidak ada modul lain menulis ke ledger.
- [ ] `events.go` — `EventEarningRecorded/EarningMatured/PayoutRequested/PayoutCompleted/
      PayoutFailed`.
- [ ] `module.go` — `Config{Pool,Tx,Outbox,Sellers seller.Port,HoldPeriod,
      PayoutMinimum money.Amount,Disburser,Clock,Logger}`; `RegisterSubscriptions` +
      `RegisterWorkers` (`MatureEarnings`, `ProcessPayouts` — job 5 & 6 di `cmd/worker`).

### Wiring
- [ ] `internal/app/registry.go` — build `settlement` setelah `order` (terima
      `sellerMod.Port()`); panggil `RegisterSubscriptions` + tambahkan job 5 & 6 ke
      `RegisterWorkers`.

## Test wajib (`../GUIDES.md` §15)

- **Buku besar seimbang**: setelah rangkaian transaksi, Σ debit = Σ kredit, saldo tiap
  akun sesuai harapan.
- **Idempotensi settlement**: `EventPaymentSettled` / `OnOrderPaid` diproses dua kali →
  satu jurnal, penjual **tidak** dibayar dua kali.
- **Refund parsial** setelah maturasi → saldo seller negatif / piutang, bukan error.
- `RequestPayout` di bawah `PayoutMinimum` → `ErrBelowMinimumPayout`; saat ada payout
  in-progress → `ErrPayoutInProgress`.
- `ProcessPayouts` gagal transfer → `MarkFailed` + jurnal payout ter-reverse (saldo
  kembali).
- `ReconciliationReport` menghasilkan selisih nol untuk skenario end-to-end.
- COD: `RecordCODCollected` memakai akun `courier/cod_receivable`, bukan `clearing`.

## Sengaja TIDAK dikerjakan di fase ini

- Keputusan "siapa menanggung gagal bayar COD" — **butuh ADR**, minta keputusan user.
- Snapshot saldo periodik — hanya bila penjumlahan terbukti lambat.
- UI laporan (frontend).
