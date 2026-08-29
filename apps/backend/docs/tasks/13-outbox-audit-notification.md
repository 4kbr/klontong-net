# Fase 13 — Outbox Worker + Audit + Notification

> Prasyarat: fase 01 + fase 10 (dan idealnya 11, 12). Guide:
> [`../GUIDES.md`](../GUIDES.md) §2 (tahap 13), §8, §9. ADR: 013 (worker proses
> terpisah). `../ARCHITECTURE.md` §11–§12.

## Tujuan

Rakit `cmd/worker` (8 pekerjaan), plus modul konsumen murni `audit` dan `notification`.
**Selesai kalau** setiap aksi penting tercatat di audit dan notifikasi terkirim.

## Aturan khusus fase ini

- Jaminan event **at-least-once** → **semua handler idempoten**. Kunci: unique index
  pada `event_id`, `INSERT ... ON CONFLICT DO NOTHING`.
- Payload event rusak → **log dan return nil** (jangan buat relay retry selamanya).
- Satu event bisa menghasilkan **dua notifikasi** (mis. `OrderPaid` → pembeli + penjual).
- Email **tidak dikirim di handler** — handler menyimpan notifikasi, worker
  `SendPending` yang mengirim (SMTP lambat tidak boleh memperlambat relay outbox).
- Audit **tidak pernah** menyimpan password hash / token / server key — filter field
  sensitif sebelum menulis. Tabel append-only, tanpa method Update.
- Semua 8 pekerjaan worker **gagal senyap** kalau tidak jalan → beri **metrik + alarm**
  tiap pekerjaan.
- Pekerjaan yang menyentuh uang (pematangan dana, pencairan — fase 14) butuh **kunci
  terdistribusi Redis** + idempotency key DB.

## Urutan kerja

### 1. `cmd/worker` + runner

- [ ] `cmd/worker/main.go` — `config.Load` → `logger.New` → ctx SIGINT/SIGTERM →
      `postgres.NewPool` → `redis.New` → `app.Deps` → `app.NewRegistry(deps)` →
      `Runner` → `registry.RegisterWorkers(runner)` → `runner.Start(ctx)`.
- [ ] `internal/app/registry.go` `RegisterWorkers(runner *Runner)` — kumpulkan job dari
      tiap modul. **8 pekerjaan** (`../ARCHITECTURE.md` §12):
  1. Relay outbox — `outbox.Relay.Run` (`FOR UPDATE SKIP LOCKED`).
  2. Pelepas reservasi kedaluwarsa — `inventory` `ReleaseExpired` (fase 6).
  3. Pembatalan pesanan tak dibayar — `order` `ExpireUnpaidOrders` (fase 12).
  4. Rekonsiliasi pembayaran — `payment` `ReconcilePending` (fase 11).
  5. Pematangan dana — `settlement` `MatureEarnings` (fase 14).
  6. Pemrosesan pencairan — `settlement` `ProcessPayouts` (fase 14).
  7. Sinkronisasi tracking — `fulfillment` `SyncTracking` (fase 8).
  8. Penyegaran agregat katalog — `catalog` refresh rating/sold (di bawah).
  > Job 5 & 6 didaftarkan di fase 14; di fase ini pasang 1–4, 7, 8.
- [ ] Tiap job: `Job{Name,Interval,Run}` dengan metrik pending/last-run/error.

### 2. Agregat katalog (job 8)

- [ ] `catalog/internal/app/` — usecase refresh (`ProductRepository.UpdateAggregates`)
      yang menghitung ulang `RatingAvg`, `RatingCount`, `SoldCount` per batch dengan jeda.
- [ ] `catalog/module.go` `RegisterWorkers` + `RegisterSubscriptions` (dengar
      `EventReviewPublished` dari fase 15, `EventSuborderCompleted` untuk sold count).

### 3. audit

- [ ] `internal/modules/audit/internal/domain/entry.go` — `ActorType` = user|seller|
      admin|system|gateway|courier; `Entry{ActorType,ActorID *uuid,ActorLabel,Action,
      TargetType,TargetID,Before/After map[string]any,IP,UserAgent,CreatedAt}`; append-
      only.
- [ ] `.../domain/repository.go` — `Repository{Create,List(filter,cursor,limit),
      ListByTarget}`.
- [ ] `.../app/handlers.go` — satu handler per event penting: `OnOrderPlaced`,
      `OnOrderCancelled`, `OnSuborderRejected`, `OnPaymentSettled`, `OnRefundCompleted`,
      `OnPayoutCompleted`, `OnStockAdjusted`, `OnPriceChanged`, `OnSellerVerified`,
      `OnSellerSuspended`, `OnVoucherRedeemed`. Idempoten (`event_id` unique, `ON CONFLICT
      DO NOTHING`). Payload rusak → log + return nil.
- [ ] `.../app/{service,query_usecase}.go` — query log admin.
- [ ] `.../infra/audit_repository.go`.
- [ ] `.../transport/rest/{routes,response,handler}.go` — **admin only**.
- [ ] `migrations/000015_create_audit_and_outbox.up.sql` / `.down.sql` — lengkapi bagian
      `audit_events` (bagian `outbox_events` sudah di fase 1).
- [ ] `module.go` — `Config{Pool,Users identity.Port,Clock,Logger}`;
      `RegisterSubscriptions`. `port.go` kosong (konsumen murni). `events.go` kosong.

### 4. notification

- [ ] `internal/modules/notification/internal/domain/notification.go` —
      `Notification{RecipientUserID,Kind,Payload,Channel,ReadAt/SentAt,FailedReason}`;
      `Preference{UserID,EmailEnabled/PushEnabled/WhatsAppEnabled,OrderUpdates/Promotions/
      SellerUpdates}`, `ShouldSend(kind,channel)`.
- [ ] `.../domain/template.go` — `Template{Kind,Channel,Subject,BodyTemplate,Locale,
      IsActive}`. Isi pesan dari template + payload event, **bukan** ditulis di Go.
- [ ] `.../domain/repository.go` — `NotificationRepository{Create,CreateMany,ListByUser,
      MarkRead,MarkAllRead,CountUnread,ListPendingSend}`, `PreferenceRepository{Get,
      Upsert}`, `TemplateRepository{FindByKindAndChannel}`.
- [ ] `.../app/handlers.go` — `OnOrderPlaced` → buyer "selesaikan pembayaran";
      `OnOrderPaid` → buyer "pembayaran diterima" + seller "pesanan baru";
      `OnSuborderConfirmed/Rejected/Shipped/Delivered` → pesan buyer; `OnOrderExpired` →
      buyer; `OnLowStock` → seller; `OnPayoutCompleted` → seller; `OnReviewPublished` →
      seller. Idempoten. **Simpan notifikasi saja**, tidak mengirim email.
- [ ] `.../app/send_usecase.go` — `SendPending()` worker: ambil notifikasi belum terkirim,
      render dari template, kirim via channel sesuai preferensi; gagal → catat + retry
      backoff, menyerah setelah batas.
- [ ] `.../app/query_usecase.go` — inbox listing, mark-read.
- [ ] `.../infra/{notification_repository,preference_repository,template_repository}.go`.
- [ ] `.../transport/rest/{routes,response,handler}.go` — inbox di `/api/v1`.
- [ ] `migrations/000014_create_notification.up.sql` / `.down.sql` —
      `notification_notifications` (channel inapp|email|whatsapp|push),
      `notification_preferences`, `notification_templates` sesuai spec.
- [ ] `module.go` — `Config{Pool,Users identity.Port,Mailer,Clock,Logger}` (tanpa Tx/
      Outbox); `RegisterSubscriptions` + `RegisterWorkers` (`SendPending`).
      `port.go` mungkin `UnreadCount(userID)`. `events.go` kosong.

### Wiring
- [ ] `internal/app/registry.go` — build `notification` + `audit` (paling akhir); panggil
      **semua** `RegisterSubscriptions` (audit, notification, inventory, settlement,
      catalog). Lupa panggil = handler tidak pernah jalan, **tanpa error**.

## Test wajib

- **Idempotensi handler**: event settlement/audit/notification yang sama diproses dua
  kali → satu baris (unique `event_id`).
- `OrderPaid` → tepat dua notifikasi (buyer + seller).
- Payload event rusak → handler return nil, relay tidak retry selamanya.
- Relay outbox: dua instance jalan bersama → tiap event dipublish sekali (`SKIP LOCKED`).
- `SendPending` gagal kirim → `FailedReason` tercatat, retry dengan backoff.

## Sengaja TIDAK dikerjakan di fase ini

- Job 5 & 6 (pematangan dana, pencairan) — didaftarkan di fase 14.
- Logika buku besar settlement (fase 14).
- Channel WhatsApp/push nyata — cukup inapp + email (MailHog).
