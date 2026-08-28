package audit

// Modul audit: catatan setiap aksi yang menyentuh uang, stok, atau status
// pesanan. Konsumen event murni.
//
// TODO:
//   type Config struct { Pool; Users identity.Port; Clock; Logger }
//   func New(cfg Config) *Module
//   func (m *Module) RegisterRoutes(r chi.Router)     // hanya admin
//   func (m *Module) RegisterSubscriptions(bus eventbus.Subscriber)
