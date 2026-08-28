package settlement

// Modul settlement: komisi, buku besar, dan pencairan ke penjual.
//
// Bagian yang paling tidak boleh salah di seluruh sistem. Uang yang hilang
// tidak bisa ditambal dengan hotfix.
//
// TODO:
//   type Config struct { Pool; Tx; Outbox; Sellers seller.Port;
//                        HoldPeriod time.Duration; PayoutMinimum money.Amount;
//                        Disburser Disburser; Clock; Logger }
//   func New(cfg Config) *Module
//   func (m *Module) RegisterRoutes(r chi.Router)
//   func (m *Module) RegisterSubscriptions(bus eventbus.Subscriber)
//   func (m *Module) RegisterWorkers(runner *app.Runner)
