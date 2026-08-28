package notification

// Modul notification: inbox, email, dan pesan lain. Konsumen event murni.
//
// TODO:
//   type Config struct { Pool; Users identity.Port; Mailer mailer.Mailer;
//                        Clock; Logger }
//   func New(cfg Config) *Module
//   func (m *Module) RegisterRoutes(r chi.Router)
//   func (m *Module) RegisterSubscriptions(bus eventbus.Subscriber)
//   func (m *Module) RegisterWorkers(runner *app.Runner)   // pengiriman tertunda
