package eventbus

// TODO:
//   type Event struct {
//       ID uuid.UUID
//       Type string              // "<modul>.<agregat>.<aksi lampau>"
//       AggregateType string
//       AggregateID uuid.UUID
//       Payload json.RawMessage
//       OccurredAt time.Time
//   }
//   func New(eventType, aggType string, aggID uuid.UUID, payload any) (Event, error)
//   func Decode[T any](e Event) (T, error)
//
// Payload event adalah KONTRAK PUBLIK. Menambah field aman; menghapus atau
// mengubah tipe merusak consumer diam-diam.
