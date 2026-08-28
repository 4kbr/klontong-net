package eventbus

// TODO:
//   type Handler func(ctx context.Context, e Event) error
//   type Publisher interface { Publish(ctx, ...Event) error }
//   type Subscriber interface { Subscribe(eventType string, h Handler) }
//   type Bus interface { Publisher; Subscriber }
//
// Modul mempublish lewat OUTBOX, bukan langsung ke bus. Yang memanggil
// bus.Publish hanyalah outbox relay.
