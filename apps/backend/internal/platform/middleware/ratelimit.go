package middleware

// TODO: token bucket. Kunci dan batas berbeda per area:
//   login & registrasi   ketat, per IP — ini yang rawan brute force
//   checkout             ketat, per user — mencegah spam pesanan
//   katalog & pencarian  longgar, per IP
//   webhook              longgar tapi ada; gateway pun bisa salah dan membanjiri
// Mulai in-memory, pindah ke Redis saat sudah multi-instance.
// Balas 429 dengan header Retry-After.
