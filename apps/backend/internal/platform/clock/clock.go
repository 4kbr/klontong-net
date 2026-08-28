package clock

// TODO: type Clock interface{ Now() time.Time }; New(); Fixed{T} untuk test.
// Usecase memanggil uc.clock.Now(), bukan time.Now(). Banyak aturan di sistem ini
// bergantung waktu (kedaluwarsa reservasi, masa tahan dana, masa berlaku voucher)
// dan mengetesnya tanpa clock palsu berarti test yang lambat dan rapuh.
