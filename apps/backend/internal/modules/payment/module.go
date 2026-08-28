package payment

// Modul payment: pembayaran lewat gateway dan COD.
//
// TODO:
//   type Config struct { Pool; Tx; Outbox; Gateway Gateway;
//                        WebhookSecret string; Expiry time.Duration; Clock; Logger }
//   func New(cfg Config) *Module
//   func (m *Module) RegisterRoutes(r chi.Router)   // termasuk /webhook/payment
//   func (m *Module) Port() Port
//   func (m *Module) RegisterWorkers(runner *app.Runner)   // rekonsiliasi
