package order

// Modul order: checkout, pesanan, dan siklus hidupnya.
//
// Titik temu seluruh sistem. Ia bergantung ke paling banyak modul, dan itu
// wajar — checkout memang tempat semuanya bertemu. Yang harus dijaga: modul lain
// TIDAK boleh bergantung balik ke `order` kecuali lewat event.
//
// TODO:
//   type Config struct {
//       Pool; Tx; Outbox
//       Cart cart.Port; Catalog catalog.Port; Pricing pricing.Port
//       Inventory inventory.Port; Promotion promotion.Port
//       Fulfillment fulfillment.Port; Payment payment.Port
//       Customers customer.Port; Sellers seller.Port
//       Clock; Logger
//   }
//   func New(cfg Config) *Module
//   func (m *Module) RegisterRoutes(r chi.Router)    // pembeli, penjual, admin
//   func (m *Module) Port() Port
//   func (m *Module) RegisterSubscriptions(bus eventbus.Subscriber)
//   func (m *Module) RegisterWorkers(runner *app.Runner)   // pembatalan kedaluwarsa
