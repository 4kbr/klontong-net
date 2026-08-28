package httpx

// Penerjemah error -> HTTP. Satu-satunya tempat status code ditentukan.
//
// TODO — func Error(w, r, err):
//   ErrInvalidInput    -> 400
//   ErrUnauthorized    -> 401
//   ErrForbidden       -> 403
//   ErrNotFound        -> 404
//   ErrConflict        -> 409
//   ErrTooManyRequests -> 429
//   ErrUpstream        -> 502
//   selain itu         -> 500
//
// Untuk 5xx: log lengkap dengan request_id, kirim pesan generik.
// Jangan pernah meneruskan pesan error gateway pembayaran ke pembeli — kadang
// memuat detail internal yang tidak layak keluar.
