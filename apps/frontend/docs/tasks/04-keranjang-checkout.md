# Fase 04 — Keranjang & checkout (buyer)

> Prasyarat: 02, 03 | Ref: kontrak `paths/cart.yaml`, `paths/shipping.yaml`,
> `paths/vouchers.yaml`, `paths/checkout.yaml`; **ADR-002, ADR-003, ADR-004**;
> `docs/ARCHITECTURE.md` §6.

## Tujuan

Keranjang multi-penjual dan alur checkout dua tahap sampai pesanan terbentuk +
instruksi pembayaran. **Selesai kalau** pembeli bisa menyelesaikan checkout
keranjang berisi barang dari beberapa penjual terhadap mock, `price_changed` dan
`out_of_stock` ditangani dengan konfirmasi ulang, dan tidak ada satu pun total
mengikat yang dihitung di klien.

## Aturan khusus fase ini

- **ADR-004.** Request `checkout/preview` dan `checkout` hanya mengirim
  **pilihan**: `address_id`, per-suborder `fulfillment` (metode + `courier_id`
  bila ada), `voucher_code`(s), daftar `{ variant_id, outlet_id?, quantity }`.
  **Tidak ada field harga/total** selain `client_grand_total` (opsional,
  pembanding). Server balas `409 price_changed` dengan `error.details` per baris
  (`old_unit_price`/`new_unit_price`) → UI tampilkan diff, minta konfirmasi,
  panggil `preview` lagi, baru `checkout`.
- **ADR-003.** Halaman keranjang menampilkan peringatan stok dari flag server per
  baris ("stok tinggal 2", "habis") **sebelum** tombol checkout. Baris habis
  tidak boleh ikut ke `preview`.
- **ADR-002.** `preview` mengembalikan rincian **per suborder** (per penjual +
  outlet terpilih + ongkir + diskon). Tampilkan terpisah; jangan gabungkan jadi
  satu baris.
- **Idempotency-Key**: buat `uuid` **sekali** saat pembeli masuk halaman
  konfirmasi checkout; simpan di state; pakai ulang untuk setiap retry `checkout`.
  Baru regen bila pembeli mengubah isi pesanan secara berarti.
- Keranjang **dikelompokkan server**; frontend tidak regroup, tidak recompute.
- `preview` boleh dipanggil berkali-kali (tiap ubah alamat/kurir/voucher/qty).
  `debounce` panggilan, tapi selalu `preview` ulang sebelum `checkout`.
- Angka di UI checkout **semuanya dari respons `preview`** (server), bukan
  akumulasi lokal. `<Money estimate>` hanya untuk perkiraan di keranjang sebelum
  `preview`.

## Urutan kerja

### Tipe & API (packages/api)
- [ ] `src/endpoints/cart.api.ts` — `getCart`, `clearCart`, `addItem`, `updateItem` (qty ≤ 0 hapus), `removeItem`, `mergeGuestCart`.
- [ ] `src/endpoints/shipping.api.ts` — `getQuote` (`/shipping/quote`), `getSuborderShipment` (`/suborders/{id}/shipment`, read-only di storefront).
- [ ] `src/endpoints/voucher.api.ts` — `validateVoucher` (`/vouchers/validate`), `listAvailableVouchers` (`/vouchers/available`).
- [ ] `src/endpoints/checkout.api.ts` — `previewCheckout`, `placeOrder` (terima `idempotencyKey`, teruskan sebagai header).
- [ ] `src/hooks/` — `useCart`, `useCartMutations` (optimistic + rollback + invalidate `queryKeys.cart`), `useShippingQuote`, `useVouchers`, `useCheckoutPreview` (mutation/`keepPreviousData`), `usePlaceOrder`.
- [ ] `src/schemas/checkout.ts` — Zod: bentuk request (tanpa harga), parser `CheckoutPreview` + `DetailedError` (`price_changed`, `out_of_stock` details).

### State (Zustand / Query)
- [ ] `checkoutStore` — `addressId`, per-suborder `fulfillmentChoice`, `voucherCodes`, `idempotencyKey`, `lastPreview`, `priceChangeAck`. Persist ringan ke `sessionStorage`.
- [ ] Aturan invalidasi: mutasi keranjang → invalidate `cart` + reset `lastPreview`.

### Komponen (packages/ui)
- [ ] `CartLineItem` (qty stepper, hapus, badge masalah), `SellerGroupCard`,
      `StockWarning`, `VoucherInput` + `VoucherList`, `AddressSelect`,
      `FulfillmentOptionList` (per suborder: pickup/local/courier + estimasi ongkir server),
      `OrderSummary` (angka dari `preview`), `PriceChangeDialog` (tabel diff per baris),
      `OutOfStockDialog`.

### Halaman & rute (storefront, di balik `RequireAuth` untuk checkout)
- [ ] `/keranjang` — grup per penjual, peringatan stok, ubah/hapus, subtotal **perkiraan** berlabel, tombol "checkout" (nonaktif bila semua baris bermasalah).
- [ ] `/checkout` — pilih alamat → per suborder pilih pengiriman → voucher →
      panel ringkasan dari `preview` (auto refresh saat pilihan berubah) →
      tombol "buat pesanan".
- [ ] Alur `price_changed`: `PriceChangeDialog` → "lihat harga baru" → `preview` ulang → user setuju → `checkout`.
- [ ] Alur `out_of_stock`: `OutOfStockDialog` → sesuaikan qty / hapus → kembali ke `preview`.
- [ ] Sukses `checkout` (201) → redirect `/pesanan/:orderId` (Fase 05) dengan instruksi pembayaran.

### Wiring
- [ ] MSW handlers **stateful**: `cart` (grup per penjual, flag stok, merge),
      `shipping` (quote per metode/jarak dummy), `vouchers`, `checkout`
      (`preview` hitung dari `db`; `place` cek `Idempotency-Key` → 1 order,
      skenario `price_changed`/`out_of_stock` bisa dipicu via fixture flag).
- [ ] Fixture: keranjang 3 penjual (pickup / local / courier), 1 varian yang
      harganya "berubah", 1 varian stok kurang.

## Test wajib

- Request `preview`/`checkout` yang dikirim **tidak mengandung** field harga/total (uji shape).
- `checkout` tanpa `Idempotency-Key` → ditolak (`400`); dengan key sama 2× → **satu** order (uji MSW).
- 3 penjual → `preview` mengembalikan 3 suborder; UI menampilkan 3 blok dengan ongkir masing-masing.
- `price_changed`: klien kirim `client_grand_total` lama → dialog diff muncul, order **tidak** dibuat sampai user konfirmasi & `preview` ulang.
- `out_of_stock`: dialog muncul, baris bisa dikurangi, `preview` ulang bersih, lalu `checkout` sukses.
- Keranjang: baris "habis" tak ikut ke `preview`; peringatan stok tampil di `/keranjang`.
- Ubah alamat → `preview` dipanggil ulang (outlet & ongkir bisa berubah).
- Angka `OrderSummary` sama persis dengan `preview.data` (tidak ada penjumlahan lokal).

## Sengaja TIDAK dikerjakan

- Halaman detail pesanan, status pembayaran, retry, cancel — Fase 05.
- Ulasan, notifikasi — Fase 05.
- Sisi seller dari suborder — Fase 08.
