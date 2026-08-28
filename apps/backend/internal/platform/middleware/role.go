package middleware

// TODO:
//   func RequireRole(role string) func(http.Handler) http.Handler
//     Untuk /api/v1/seller/* dan /api/v1/admin/*.
//
// Peran hanyalah gerbang kasar. Bahwa seseorang berperan `seller` tidak berarti
// ia berhak mengubah toko MANA PUN — kepemilikan toko diperiksa di usecase.
// Middleware ini hanya menyaring yang jelas tidak berkepentingan.
