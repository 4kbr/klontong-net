package auth

// TODO:
//   type Claims struct { UserID uuid.UUID; Roles []string; jwt.RegisteredClaims }
//   type TokenIssuer interface { IssueAccess(userID uuid.UUID, roles []string) (string, time.Time, error) }
//   type TokenVerifier interface { Verify(token string) (Claims, error) }
//   func NewJWT(cfg config.JWTConfig) *JWT
//
// Access token pendek (15m), tidak bisa dicabut. Kemampuan mencabut ada di
// refresh token yang tersimpan di database.
//
// Verify wajib memeriksa metode signing secara eksplisit (cegah "alg: none").
//
// Peran dititipkan di klaim supaya middleware tidak perlu query database di
// setiap request. Konsekuensinya: pencabutan peran baru berlaku setelah token
// kedaluwarsa. Untuk pencabutan yang harus segera (penjual di-suspend),
// periksa status di usecase, bukan mengandalkan token.
