# Fase 08 — Seller pesanan, fulfillment & keuangan

> Prasyarat: 06 (butuh juga 07 untuk data produk/stok yang layak) | Ref: kontrak
> `paths/seller-orders.yaml`, `paths/seller-fulfillment.yaml`,
> `paths/seller-promotion.yaml`, `paths/seller-finance.yaml`,
> `paths/seller-reviews.yaml`; **ADR-002, ADR-009** (state machine suborder),
> **ADR-011** (buku besar), ADR-012 (COD).

## Tujuan

Penjual menjalankan siklus suborder end-to-end, mengelola pengiriman & zona
antar, voucher toko, serta melihat saldo/pendapatan/pencairan dan membalas
ulasan. **Selesai kalau** penjual bisa memproses satu suborder dari
`awaiting_confirmation` sampai `shipped`/`ready_for_pickup`, dan melihat saldo +
mengajukan pencairan.

## Aturan khusus fase ini

- **ADR-002.** Penjual hanya melihat/menyentuh **suborder miliknya**
  (`/api/v1/seller/orders/{suborderId}` — perhatikan: path `orders` tapi id-nya
  **suborder**, sesuai kontrak). Tidak pernah melihat bagian penjual lain.
- **ADR-009.** Transisi status hanya lewat aksi yang diizinkan state machine.
  UI hanya menampilkan aksi yang valid dari status sekarang
  (`confirm`/`reject` dari `awaiting_confirmation`; `pack` → `ship` /
  `ready-for-pickup`; tak ada cancel setelah `shipped`). Sisipkan alasan wajib
  pada `reject`.
- **Menolak satu suborder tidak membatalkan order induk** — komunikasikan
  "ini hanya pesanan toko Anda; refund parsial diproses otomatis".
- Ongkir & metode kirim per suborder mengikuti pilihan pembeli; penjual memenuhi
  metode itu (pickup: siapkan; local/courier: buat/booking shipment sesuai
  endpoint, unggah bukti).
- **ADR-011.** Saldo = penjumlahan entri buku besar; tampilkan `pending` vs
  `available` terpisah, jangan hitung sendiri. Pencairan mengurangi `available`.
- **ADR-012.** Pesanan COD: `payment.status` tetap `pending` sampai setoran;
  tampilkan label COD dan jangan sajikan seolah sudah dibayar.
- Balas ulasan (`/seller/reviews/{id}/reply`) satu kali per ulasan (ikuti respons).

## Urutan kerja

### Tipe & API (packages/api)
- [ ] `src/endpoints/seller-order.api.ts` — `list` (`/seller/orders`, filter status + keyset),
      `get` (`/seller/orders/{suborderId}`), `confirm`, `reject` (alasan), `pack`, `ship` (no resi/kurir), `readyForPickup`.
- [ ] `src/endpoints/seller-fulfillment.api.ts` — `getSuborderShipment`
      (`/seller/orders/{suborderId}/shipment`), `getShipment`/`updateShipment`
      (`/seller/shipments/{id}`), `confirmPickup` (`/seller/shipments/{id}/confirm-pickup`),
      `uploadProof` (`/seller/shipments/{id}/proof`), `getDeliveryZone`/`setDeliveryZone`
      (`/seller/outlets/{id}/delivery-zone`).
- [ ] `src/endpoints/seller-promotion.api.ts` — voucher toko `list/create/get/update/delete`.
- [ ] `src/endpoints/seller-finance.api.ts` — `getBalance`, `listEarnings`, `listPayouts`/`requestPayout`, `getPayout` (`/seller/payouts/{id}`).
- [ ] `src/endpoints/seller-review.api.ts` — `listReviews`, `replyReview`.
- [ ] `src/hooks/` — `useSellerOrders`, `useSellerSuborder`, `useSuborderActions`,
      `useShipment`, `useDeliveryZone`, `useSellerVouchers`, `useBalance`,
      `useEarnings`, `usePayouts`, `usePayoutMutations`, `useSellerReviews`.
- [ ] `src/schemas/` — Zod: alasan reject, form ship (resi/kurir), bukti kirim, zona antar (radius/tarif dasar/per km/minimum), voucher toko, pengajuan pencairan.

### State (Zustand / Query)
- [ ] `queryKeys.seller.orders/*`, `queryKeys.seller.finance/*`.
- [ ] Invalidasi: aksi suborder → invalidate `sellerSuborder(id)` + `sellerOrders` + `balance`/`earnings` bila relevan.
- [ ] Aksi suborder pakai mutation dengan konfirmasi (dialog) + optimistic status opsional.

### Komponen (packages/ui / dashboard)
- [ ] `SuborderQueue` (kolom per status / tab), `SuborderDetail` (item snapshot,
      alamat kirim, metode, ongkir, komisi, `SellerEarning`), `SuborderActionBar`
      (aksi valid dari state machine), `RejectDialog` (alasan wajib + catatan
      "tidak membatalkan order induk"), `ShipForm`, `ProofUploader`,
      `PickupConfirm`, `DeliveryZoneForm`, `VoucherForm`, `BalanceCards`
      (`pending`/`available`), `EarningsTable`, `PayoutList` + `RequestPayoutDialog`,
      `ReviewReplyForm`.

### Halaman & rute (dashboard, `RequireRole('seller')`)
- [ ] `/seller/pesanan` — `SuborderQueue` + filter.
- [ ] `/seller/pesanan/:suborderId` — `SuborderDetail` + `SuborderActionBar` + shipment panel.
- [ ] `/seller/pengiriman/:shipmentId` — detail + confirm-pickup + proof.
- [ ] `/seller/outlet/:outletId/zona-antar` — `DeliveryZoneForm`.
- [ ] `/seller/voucher` — daftar + form.
- [ ] `/seller/keuangan` — `BalanceCards` + `EarningsTable` + `PayoutList` + ajukan pencairan.
- [ ] `/seller/ulasan` — daftar + balas.

### Wiring
- [ ] MSW handlers stateful: `seller-orders` (transisi status tervalidasi state machine; reject → refund parsial dicatat di order pembeli; order induk tetap jalan),
      `seller-fulfillment`, `seller-promotion`, `seller-finance` (saldo dari entri dummy; pencairan mengurangi `available`), `seller-reviews`.
- [ ] Fixture: 3 suborder (pickup / local / courier) di berbagai status; entri earning + 1 payout.

## Test wajib

- `SuborderActionBar` hanya menampilkan aksi valid: dari `shipped` tidak ada tombol cancel/reject.
- `reject` tanpa alasan → ditolak; dengan alasan → suborder `rejected`, dan MSW menunjukkan order induk **tidak** batal + refund parsial tercatat.
- `ship` mengisi resi → status `shipped`; `ready-for-pickup` hanya untuk suborder metode pickup.
- Penjual A membuka suborder milik penjual B → `403`/tidak ditemukan; UI menampilkan akses ditolak.
- `BalanceCards`: `available` turun tepat sebesar pencairan setelah `requestPayout`; tidak ada penjumlahan saldo di klien.
- Order COD ditandai COD dan tidak ditampilkan sebagai "sudah dibayar".
- Balas ulasan → muncul di bawah ulasan; balas kedua ditolak sesuai respons.

## Sengaja TIDAK dikerjakan

- Verifikasi toko, kategori, voucher marketplace, rekonsiliasi, refund admin — Fase 09.
- E2E lintas peran — Fase 10.
