package domain

// TODO:
//   type Method string   // local_delivery | courier | pickup
//   type Status string   // pending | ready | picked_up | in_transit |
//                        // delivered | failed | returned
//   type Shipment struct {
//       ID; SuborderID; OutletID; Method; Status
//       CourierCode, CourierService, TrackingNumber string
//       DriverName, DriverPhone string; DistanceKm float64
//       PickupCode string; PickupExpiresAt, PickedUpAt *time.Time
//       ShippingAmount money.Amount
//       EstimatedDeliveryAt, ShippedAt, DeliveredAt *time.Time
//       ProofStorageKey string
//       CreatedAt; UpdatedAt
//   }
//   func NewShipment(...) (*Shipment, error)
//   func (s *Shipment) Transition(to Status, now time.Time) error
//   func (s *Shipment) IsPickup() bool
//   func (s *Shipment) GeneratePickupCode() error
//
// Satu tabel untuk tiga metode, dengan kolom yang hanya terisi sesuai metodenya.
// Memecah jadi tiga tabel terdengar lebih bersih tapi membuat "tampilkan semua
// pengiriman penjual ini" jadi tiga query dan satu union.
//
// PickupCode: kode yang ditunjukkan pembeli saat mengambil barang. Beri masa
// berlaku — barang yang tidak diambil harus kembali jadi stok, dan tanpa batas
// waktu tidak ada yang tahu kapan itu terjadi.
//
// ProofStorageKey: foto bukti terima. Untuk antar lokal ini sering jadi
// satu-satunya bukti saat pembeli mengaku barang tidak sampai.
