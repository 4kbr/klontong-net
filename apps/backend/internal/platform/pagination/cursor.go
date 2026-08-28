package pagination

// TODO: cursor keyset.
//   type Cursor struct { SortValue string; ID uuid.UUID }
//   Encode/Decode base64 JSON
//   type Page[T any] struct { Items []T; NextCursor string; HasMore bool }
//
// Katalog perlu diurutkan bermacam-macam (terbaru, termurah, terlaris, rating).
// Cursor harus membawa nilai kolom pengurut, bukan hanya waktu. Ambil n+1 baris
// untuk tahu HasMore tanpa COUNT.
