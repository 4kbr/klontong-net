package mailer

// TODO:
//   type Mail struct { To []string; Subject, Template string; Data map[string]any }
//   type Mailer interface { Send(ctx, Mail) error }
//   func NewSMTP(cfg config.SMTPConfig) Mailer
//
// Di dev arahkan ke MailHog :1025.
// Pengiriman email TIDAK dilakukan di jalur request — modul notification yang
// menanganinya lewat event, supaya checkout tidak melambat karena SMTP lambat.
