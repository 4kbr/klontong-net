package id

// TODO:
//   func New() uuid.UUID       // uuid v7 bila tersedia — terurut waktu, ramah index
//   func Parse(s string) (uuid.UUID, error)
//
// TODO — nomor yang dilihat manusia, ini berbeda dari id internal:
//   func OrderNumber(t time.Time, seq int64) string      // "KN-20260827-000123"
//   func SuborderNumber(orderNumber string, idx int) string  // "KN-20260827-000123-1"
//
// Nomor pesanan dipakai di percakapan customer service dan di transfer bank.
// Ia harus pendek, mudah dibacakan lewat telepon, dan tidak ambigu — hindari
// karakter yang mirip (0/O, 1/I/l). UUID tidak memenuhi syarat itu.
