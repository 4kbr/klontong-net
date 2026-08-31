# Fase 11 — Payment + Webhook + Rekonsiliasi

> Prasyarat: fase 01. (Dipakai oleh fase 10 — `CreatePayment` + `gateway_noop` bisa
> disiapkan lebih dulu sebelum sisa fase ini.) Guide: [`../GUIDES.md`](../GUIDES.md) §2
> (tahap 11), §12. ADR: 008 (idempotensi + rekonsiliasi), 012 (COD), 014 (REST langsung).
> **Fase menyentuh uang — hati-hati, `webhook_usecase` jalur paling berisiko.**
> Kontrak (**mengikat**, ADR-015): [`paths/webhook.yaml`](../../../../contracts/openapi/paths/webhook.yaml), [`orders.yaml`](../../../../contracts/openapi/paths/orders.yaml) (payment), [`admin.yaml`](../../../../contracts/openapi/paths/admin.yaml) (payments, refunds).

## Tujuan

Gateway pembayaran + COD + webhook + rekonsiliasi. **Selesai kalau** pesanan jadi `paid`
lewat webhook palsu.

## Aturan khusus fase ini

- `Payment.MarkSettled` **idempoten**: `ProviderReference` sama dua kali = no-op tanpa
  error. Webhook datang berulang dan bisa **tidak berurutan** (`settled` sebelum
  `pending`) — cek transisi, jangan asumsikan urutan.
- Webhook: verifikasi signature **waktu-konstan** dulu; payload tak terverifikasi
  **tidak pernah** diproses. Webhook duplikat dibalas **200** (bukan error — gateway
  kirim ulang sampai 200).
- `ErrAmountMismatch` (gateway laporkan jumlah beda) → **catat insiden, tandai untuk
  review manusia, JANGAN tandai order paid**.
- COD: `CreatePayment` **tanpa panggil gateway**; status mulai `pending`, jadi `settled`
  hanya saat setoran kurir masuk (ADR-012). COD tetap lewat `CreatePayment` supaya
  `order` tidak perlu bercabang.
- Panggilan gateway **tidak di dalam tx checkout** — setelah commit atau dua langkah.
- Refund menunjuk **suborder**, bukan cuma order. Σ semua refund ≤ pembayaran — cek di
  dalam tx.
- REST provider langsung dengan `net/http`, tanpa SDK (ADR-014). `Idempotency-Key` untuk
  operasi yang menciptakan sesuatu.
- Jangan pernah andalkan webhook sebagai satu-satunya jalur → worker rekonsiliasi
  `Inquiry` untuk pembayaran menggantung (ADR-008).

## Urutan kerja

### Domain
- [ ] `internal/modules/payment/internal/domain/gateway.go` — `ChargeRequest
      {OrderNumber,Amount,Channel,Customer,Items,ExpiresAt,IdempotencyKey}`,
      `ChargeResult{ProviderReference,RedirectURL,VANumber,QRString,Status}`,
      `Gateway interface{Charge,Inquiry(providerRef)->Status,Refund(providerRef,amount,
      reason)->string,VerifyWebhook(rawBody,headers)->WebhookPayload}`.
- [ ] `.../domain/payment.go` — `Method` = gateway|cod; `Status` = pending|authorized|
      settled|failed|expired|refunded|partially_refunded; `Payment{OrderID,Method,
      Provider,ProviderReference,Amount,Status,Channel,PaidAt/ExpiredAt,FailedReason,
      IdempotencyKey}`; `MarkSettled(at,providerRef)` **idempoten**, `MarkFailed`,
      `MarkExpired`, `CanRefund`.
- [ ] `.../domain/refund.go` — `Refund{PaymentID,OrderID,SuborderID *uuid,Amount,Reason,
      Status,ProviderReference,RequestedBy,ProcessedAt}`; `IsPartial(paymentAmount)`.
- [ ] `.../domain/repository.go` — `PaymentRepository{Create,Update,FindByID,FindByOrder,
      FindByProviderRef,ListPending(olderThan)}`, `WebhookRepository{Record(*WebhookEvent)
      -> (isDuplicate, err),MarkProcessed,MarkFailed}`, `RefundRepository{Create,Update,
      ListByPayment,SumByPayment}`.
- [ ] `.../domain/errors.go` — `ErrPaymentAlreadySettled`, `ErrPaymentExpired`,
      `ErrInvalidSignature`, `ErrAmountMismatch`, `ErrRefundExceedsPayment`,
      `ErrUnsupportedChannel`, `ErrGatewayUnavailable`.

### Migrasi
- [ ] `migrations/000009_create_payment.up.sql` / `.down.sql` — `payment_payments`
      (`method` gateway|cod, `provider_reference`, status enum, amount `bigint`),
      `payment_refunds`, `payment_webhook_events` sesuai spec; **partial unique index**
      satu payment aktif per order; unique `(provider, event_id)` untuk webhook; unique
      `idempotency_key`.

### Infra
- [ ] `.../infra/gateway_noop.go` — implementasi palsu yang bisa dipicu manual (kerjakan
      **lebih awal**, dipakai fase 10).
- [ ] `.../infra/gateway_midtrans.go` — impl `domain.Gateway`: Charge (Core API/Snap);
      `VerifyWebhook` = SHA512 `order_id + status_code + gross_amount + server_key`,
      compare waktu-konstan; map status (capture/settlement → settled, pending → pending,
      deny/cancel/expire → failed/expired). REST langsung, no SDK.
- [ ] `.../infra/{payment_repository,refund_repository,webhook_repository,mapper}.go` —
      `Record` kembalikan flag duplikat, bukan error.

### App
- [ ] `.../app/payment_usecase.go` — `CreatePayment(orderID,amount,method,channel,
      idempotencyKey)`: (1) cek kunci idempotensi → kembalikan hasil sebelumnya; (2)
      `cod` → buat Payment `pending` **tanpa gateway**, kembalikan instruksi "bayar di
      tempat"; (3) `gateway` → `gateway.Charge`, simpan `ProviderReference`; (4) simpan
      `ExpiredAt = now + Expiry`. Plus `GetByOrder`, `CancelPayment`, `MarkCODCollected`.
- [ ] `.../app/webhook_usecase.go` (**paling berisiko**) — `HandleWebhook(rawBody,
      headers)`: (1) `gateway.VerifyWebhook` — signature buruk → tolak + log, jangan
      proses; (2) `WebhookRepository.Record` — duplikat → balas 200 dan berhenti; (3)
      cari Payment by `ProviderReference`; (4) **cocokkan jumlah** — beda →
      `ErrAmountMismatch`, tandai untuk manusia, JANGAN paid; (5) `WithinTx` update status
      + outbox `EventPaymentSettled`; (6) balas 200 cepat, proses lanjut lewat event.
- [ ] `.../app/refund_usecase.go` — `RequestRefund(RefundInput) -> uuid`: cek Σ refund ≤
      payment via `SumByPayment` di dalam tx; buat Refund `requested`; worker panggil
      `gateway.Refund`. COD: refund = batalkan tagihan (uang tak pernah masuk).
- [ ] `.../app/reconcile_usecase.go` — `ReconcilePending()` worker: pembayaran lama
      menggantung → `gateway.Inquiry` status sebenarnya. Plus rekonsiliasi harian bandingkan
      total kita vs laporan settlement gateway (selisih nol).
- [ ] `.../app/{service,dto,port_adapter}.go`.

### Transport
- [ ] `.../transport/rest/{routes,request,response,payment_handler,webhook_handler,
      admin_handler,handler}.go` — `/webhook/payment` (verifikasi signature, tanpa sesi,
      tanpa CORS); status pembayaran di `/api/v1`; admin di `/api/v1/admin`.

### Port
- [ ] `internal/modules/payment/port.go` — `CreatePayment(orderID,amount,method,channel,
      idempotencyKey) -> Instruction`, `GetByOrder`, `RequestRefund(input) -> uuid`.
- [ ] `events.go` — `EventPaymentCreated/Settled/Failed/Expired/RefundCompleted`.
      `EventPaymentSettled` payload wajib bawa `OrderID` + jumlah yang benar-benar
      diterima.
- [ ] `module.go` — `Config{Pool,Tx,Outbox,Gateway,WebhookSecret,Expiry,Clock,Logger}`;
      worker rekonsiliasi (didaftarkan fase 13).

### Wiring
- [ ] `internal/app/registry.go` — build `payment`; pilih implementasi `Gateway`
      berdasarkan config (`gateway_noop` untuk dev).

## Test wajib

- **Idempotensi webhook**: webhook `settled` sama dua kali → satu perubahan status, dua-
  duanya dibalas 200.
- Webhook signature salah → ditolak, order tidak berubah.
- `ErrAmountMismatch`: gateway laporkan jumlah beda → order **tidak** jadi paid, insiden
  tercatat.
- Webhook tidak berurutan (`settled` sebelum `pending`) → status akhir benar.
- COD `CreatePayment` tidak memanggil gateway; status `pending`.
- `RequestRefund` melebihi sisa pembayaran → `ErrRefundExceedsPayment`.

## Sengaja TIDAK dikerjakan di fase ini

- Jurnal buku besar akibat `EventPaymentSettled` (fase 14).
- Pendaftaran worker rekonsiliasi di `cmd/worker` (fase 13).
- Setoran COD end-to-end ke buku besar (fase 14, ADR-012).
