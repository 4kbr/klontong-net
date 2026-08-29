# Fase 05 — Pesanan, pembayaran, ulasan, notifikasi (buyer)

> Prasyarat: 04 | Ref: kontrak `paths/orders.yaml`, `paths/reviews.yaml`,
> `paths/notifications.yaml`, `paths/shipping.yaml#/suborderShipment`;
> **ADR-002, ADR-009** (status suborder), ADR-010 (snapshot).

## Tujuan

Pembeli melacak pesanan multi-suborder, membayar, membatalkan bila boleh,
mengonfirmasi terima per suborder, menulis ulasan dari pembelian yang selesai,
dan mengelola notifikasi. **Selesai kalau** pembeli bisa membuka pesanan dengan
beberapa suborder berstatus berbeda, membayar lewat mock, dan status diperbarui.

## Aturan khusus fase ini

- **ADR-002.** Detail pesanan = **order induk + array suborder**. Tampilkan
  progres per suborder ("Toko A: dikirim", "Toko B: menunggu konfirmasi").
  Status induk hanya ringkasan turunan — jangan jadikan satu-satunya info.
- **Pembatalan bersifat parsial.** `POST /orders/{id}/cancel` untuk order yang
  belum ada suborder terkirim; setelah sebagian terkirim, pembatalan/penolakan
  terjadi per suborder di sisi seller (refund parsial). UI menjelaskan ini, tidak
  menjanjikan "batalkan semua".
- `confirm-received` ada **per suborder** (`/suborders/{id}/confirm-received`),
  hanya aktif saat suborder `delivered`. `delivered` ≠ `completed` (masa komplain).
- Ulasan hanya untuk item dari pembelian **selesai/diterima**; `order_item_id`
  berasal dari data pesanan, **tidak** dari input bebas user.
- Instruksi pembayaran & status datang dari server; frontend hanya menampilkan
  dan menyediakan tombol `retry`. Tidak ada polling agresif — gunakan
  `refetchInterval` moderat pada status pembayaran `pending`, berhenti saat final.
- Nominal refund/pembayaran ditampilkan apa adanya (integer rupiah).

## Urutan kerja

### Tipe & API (packages/api)
- [ ] `src/endpoints/order.api.ts` — `listOrders` (filter status + keyset),
      `getOrder` (`/orders/{id}`), `cancelOrder`, `getPayment` (`/orders/{id}/payment`),
      `retryPayment` (`/orders/{id}/payment/retry`), `confirmReceived` (`/suborders/{id}/confirm-received`).
- [ ] `src/endpoints/review.api.ts` (lengkapi) — `createReview` (`/reviews`),
      `getReview`/`updateReview`/... (`/reviews/{id}`), `reportReview` (`/reviews/{id}/report`),
      `listMyReviews` (`/me/reviews`).
- [ ] `src/endpoints/notification.api.ts` — `list`, `unreadCount`, `readAll`, `markRead` (`/notifications/{id}/read`), `getPreferences`/`updatePreferences`.
- [ ] `src/hooks/` — `useOrders`, `useOrder`, `useOrderMutations`, `usePayment`
      (dengan `refetchInterval` kondisional), `useReviewMutations`, `useMyReviews`,
      `useNotifications`, `useUnreadCount` (poll ringan), `useNotificationPrefs`.
- [ ] `src/schemas/review.ts` — Zod form ulasan (rating 1–5, teks, foto opsional).

### State (Zustand / Query)
- [ ] `queryKeys.orders(filter)`, `queryKeys.order(id)`, `queryKeys.payment(orderId)`, `queryKeys.notifications`.
- [ ] Invalidasi: `cancelOrder`/`confirmReceived` → invalidate `order(id)` + `orders`.
- [ ] `markRead`/`readAll` → optimistic update `unreadCount`.

### Komponen (packages/ui)
- [ ] `SuborderTimeline` (peta status ADR-009), `SuborderCard` (item snapshot, ongkir, status, aksi kontekstual),
      `OrderStatusSummary` ("2 dari 3 toko dikirim"), `PaymentInstructionPanel` (VA/QR/COD + countdown + retry),
      `RefundNotice`, `ReviewForm`, `ReviewList`, `NotificationItem`, `NotificationBell` (+ unread badge), `PreferencesForm`.

### Halaman & rute (storefront, `RequireAuth`)
- [ ] `/pesanan` — daftar + tab status.
- [ ] `/pesanan/:orderId` — ringkasan induk + daftar `SuborderCard` + `PaymentInstructionPanel` + tombol `cancel` (bila `CanBeCancelledByBuyer`) + `confirm-received` per suborder `delivered`.
- [ ] `/pesanan/:orderId/bayar` (atau panel) — instruksi + retry.
- [ ] `/ulasan/tulis?orderItemId=` — form; muncul sebagai CTA di suborder `completed`.
- [ ] `/akun/ulasan` — `listMyReviews` + edit/hapus.
- [ ] `/notifikasi` + `/akun/notifikasi` (preferensi). `NotificationBell` di header.

### Wiring
- [ ] MSW handlers **stateful**: `orders` (dibuat oleh checkout Fase 04; transisi
      status suborder bisa dipicu manual/fixture untuk simulasi seller),
      `payment` (pending → paid setelah "bayar" ditekan), `reviews`, `notifications`.
- [ ] Fixture: 1 order dengan suborder `shipped` + `awaiting_confirmation` + `completed`.

## Test wajib

- Detail pesanan menampilkan tiap suborder dengan status & nominal sendiri; tidak ada satu status tunggal yang menutupi perbedaan.
- Tombol `cancel` hilang setelah ada suborder `shipped`; UI menjelaskan pembatalan parsial via seller.
- `confirm-received` hanya aktif pada suborder `delivered`; setelah ditekan → status `completed`, CTA ulasan muncul.
- `retryPayment` memanggil endpoint retry; status pembayaran `pending` → `paid` setelah aksi mock, `refetchInterval` berhenti.
- Ulasan: submit dari `order_item_id` pesanan selesai → sukses; mencoba review item yang belum selesai → ditolak (`code` dari server) dan UI tidak menawarkannya.
- `readAll` → `unreadCount` jadi 0 secara optimistic lalu terkonfirmasi.

## Sengaja TIDAK dikerjakan

- Aksi seller atas suborder (confirm/ship/…) — Fase 08.
- Balasan seller atas ulasan — Fase 08.
- Refund/settlement sisi admin — Fase 09.
