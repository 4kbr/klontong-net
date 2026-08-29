# Fase 08 — Fulfillment

> Prasyarat: fase 03. Guide: [`../GUIDES.md`](../GUIDES.md) §2 (tahap 8), §12.
> ADR terkait: ADR-014 (REST langsung, bukan SDK). Faktor jarak lokal → **butuh ADR
> baru** (lihat "Yang masih terbuka" di `../DECISIONS.md`).

## Tujuan

Tiga metode kirim (`pickup`, `local_delivery`, `courier`), ongkir, tracking. **Selesai
kalau** ongkir per suborder tampil dan `pickup` gratis.

## Aturan khusus fase ini

- Ongkir dihitung **per suborder**, dari koordinat outlet masing-masing. Satu pesanan 3
  penjual = 3 perhitungan.
- `Quote` punya **masa berlaku** (~30 menit) `ExpiresAt`. Checkout dengan quote
  kedaluwarsa → hitung ulang, bukan pakai angka basi.
- Kegagalan agregator kurir **tidak boleh menggagalkan seluruh quote** — tampilkan opsi
  lain (pickup / local), catat masalahnya. Bungkus panggilan `Rates` dengan timeout
  ketat (2–3 dtk).
- `local_delivery` butuh: outlet mendukung + alamat punya koordinat + dalam radius +
  di atas `MinOrderAmount`. Alamat tanpa pin → `ErrNoCoordinates` dengan pesan yang
  mengarahkan pembeli menaruh pin.
- Panggil REST provider langsung dengan `net/http`, **tanpa SDK** (ADR-014), di balik
  interface `CourierProvider`. Sediakan implementasi palsu.
- Verifikasi signature webhook (kalau ada) dengan perbandingan **waktu-konstan**.

## Urutan kerja

### Domain
- [ ] `internal/modules/fulfillment/internal/domain/quote.go` — `Quote{CartID/SuborderID
      *uuid,SellerID,OutletID,Method,CourierCode/Service,Amount,ETDMinDay/MaxDay,
      DistanceKm,ExpiresAt}`; `IsValid(now)`.
- [ ] `.../domain/shipment.go` — `Method` = local_delivery|courier|pickup; `Status` =
      pending|ready|picked_up|in_transit|delivered|failed|returned;
      `Shipment{SuborderID,OutletID,Method,Status,Courier*,TrackingNumber,Driver*,
      DistanceKm,PickupCode,PickupExpiresAt/PickedUpAt,ShippingAmount,
      EstimatedDeliveryAt/ShippedAt/DeliveredAt,ProofStorageKey}`; `NewShipment`,
      `Transition(to,now)`, `IsPickup`, `GeneratePickupCode`.
- [ ] `.../domain/courier.go` — `RateRequest{origin/dest postal+lat/lng,WeightGram,
      ItemValue}`, `Rate{CourierCode,ServiceName,Amount,ETDMin/MaxDay}`,
      `CourierProvider interface{Rates,CreateBooking,Track}`. Cache Redis keyed (origin,
      dest, rounded weight).
- [ ] `.../domain/local_tariff.go` — `LocalTariff{BaseFare,PerKmFare,MaxDistanceKm,
      MinOrderAmount}`; `Calculate(distanceKm)` (di atas `MaxDistanceKm` → error, opsi
      tidak ditawarkan; bulatkan jarak ke atas per km, konsisten); `Haversine(...)`
      garis lurus + faktor pengali (catat sebagai ADR).
- [ ] `.../domain/repository.go` — `ShipmentRepository{Create,Update,FindByID,
      FindBySuborder,FindByTracking,ListActive,ListBySeller}`,
      `QuoteRepository{CreateMany,FindByID,DeleteExpired}`,
      `TrackingRepository{CreateMany,ListByShipment}`, `ZoneRepository{FindByOutlet,
      Upsert}`.
- [ ] `.../domain/errors.go` — `ErrQuoteExpired`, `ErrMethodNotSupported`,
      `ErrOutOfDeliveryRange`, `ErrNoCoordinates`, `ErrBelowMinimumOrder`,
      `ErrCourierUnavailable`, `ErrInvalidPickupCode`, `ErrPickupExpired`.

### Migrasi
- [ ] `migrations/000010_create_fulfillment.up.sql` / `.down.sql` —
      `fulfillment_shipments` (satu tabel tiga metode, status enum, field kurir),
      `fulfillment_quotes` (`expires_at`), `fulfillment_tracking_events`,
      `fulfillment_zones` sesuai spec. Ongkir `bigint`, jarak `numeric`.

### Infra
- [ ] `.../infra/courier_provider.go` — impl `CourierProvider` (REST langsung) + versi
      palsu untuk dev.
- [ ] `.../infra/{quote_repository,shipment_repository,tracking_repository,mapper}.go`.

### App
- [ ] `.../app/quote_usecase.go` — `Quote(QuoteRequest) -> []Option`: (1) outlet via
      `seller.Port` (koordinat asal + metode didukung), (2) kumpulkan opsi: pickup (biaya
      nol), local_delivery (bila didukung + dest berkoordinat + dalam radius + di atas
      min), courier (panggil `CourierProvider.Rates`), (3) simpan semua sebagai `Quote`
      dengan `ExpiresAt`, (4) hasil kosong = error yang bisa dijelaskan. `ValidateQuote
      (quoteID,at)` → expired → `ErrQuoteExpired`.
- [ ] `.../app/shipment_usecase.go` — `CreateShipment` (courier →
      `CourierProvider.CreateBooking`, simpan AWB; local → catat nama/telp kurir; pickup
      → terbitkan `PickupCode` + expiry), `UpdateShipmentStatus`, `ConfirmPickup
      (shipmentID,code)` (compare waktu-konstan, rate-limit per shipment),
      `UploadDeliveryProof`.
- [ ] `.../app/tracking_usecase.go` — `SyncTracking()` worker: tiap shipment kurir aktif
      panggil `Track`, simpan event baru; saat delivered publish
      `EventShipmentDelivered`.
- [ ] `.../app/{service,dto,port_adapter,handlers}.go`.

### Transport
- [ ] `.../transport/rest/{routes,request,response,quote_handler,shipment_handler,
      handler}.go` — quote di `/api/v1` (checkout), aksi shipment di `/api/v1/seller`.

### Port
- [ ] `internal/modules/fulfillment/port.go` — `Quote(QuoteRequest) -> []Option`
      (**per suborder**), `ValidateQuote(quoteID,at) -> Option`,
      `CreateShipment(input) -> uuid`.
- [ ] `events.go` — `EventShipmentCreated/PickedUp/InTransit/Delivered/Failed/Returned`.
- [ ] `module.go` — `Config{...,Sellers seller.Port,Courier CourierProvider,LocalTariff
      LocalTariffConfig}`; subscription + worker tracking (didaftarkan fase 13).

### Wiring
- [ ] `internal/app/registry.go` — build `fulfillment` setelah `seller`.

## Test wajib

- `pickup` selalu ongkir nol.
- `local_delivery` di luar `MaxDistanceKm` → opsi tidak muncul; di bawah `MinOrderAmount`
  → `ErrBelowMinimumOrder`.
- `CourierProvider.Rates` gagal/timeout → quote tetap kembali dengan opsi pickup/local,
  bukan error total.
- `Quote` kedaluwarsa → `ValidateQuote` → `ErrQuoteExpired`.
- `ConfirmPickup` dengan kode salah → `ErrInvalidPickupCode`; compare waktu-konstan.

## Sengaja TIDAK dikerjakan di fase ini

- Pemanggilan `Quote` dari checkout preview & `CreateShipment` dari `MarkShipped`
  (fase 10 / 12).
- Pendaftaran worker tracking di `cmd/worker` (fase 13).
- ADR faktor jarak lokal — **tandai, minta keputusan user**, jangan pilih sendiri.
