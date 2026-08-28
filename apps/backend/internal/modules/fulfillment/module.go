package fulfillment

// Modul fulfillment: tiga metode pengiriman.
//   local_delivery  antar sendiri oleh penjual, tarif per km
//   courier         kurir ekspedisi lewat agregator
//   pickup          diambil pembeli di outlet
//
// TODO:
//   type Config struct { Pool; Tx; Outbox; Sellers seller.Port;
//                        Courier CourierProvider; LocalTariff LocalTariffConfig;
//                        Clock; Logger }
//   func New(cfg Config) *Module
//   func (m *Module) RegisterRoutes(r chi.Router)
//   func (m *Module) Port() Port
//   func (m *Module) RegisterSubscriptions(bus eventbus.Subscriber)
//   func (m *Module) RegisterWorkers(runner *app.Runner)   // sinkronisasi tracking
