package inventory

// Modul inventory: stok per outlet, reservasi, mutasi.
//
// Modul yang paling sering jadi sumber keluhan kalau salah: pembeli melihat
// stok ada lalu gagal di checkout, atau penjual kehabisan barang yang menurut
// sistem masih ada.
//
// TODO:
//   type Config struct { Pool; Tx; Outbox; Catalog catalog.Port;
//                        ReservationTTL time.Duration; Clock; Logger }
//   func New(cfg Config) *Module
//   func (m *Module) RegisterRoutes(r chi.Router)
//   func (m *Module) Port() Port
//   func (m *Module) RegisterSubscriptions(bus eventbus.Subscriber)
//   func (m *Module) RegisterWorkers(runner *app.Runner)   // pelepas reservasi kedaluwarsa
