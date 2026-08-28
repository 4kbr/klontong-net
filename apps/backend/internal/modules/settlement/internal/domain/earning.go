package domain

// TODO:
//   type Earning struct { ID; SuborderID; SellerID
//                         GrossAmount, CommissionAmount, NetAmount money.Amount
//                         Status string        // pending | matured | paid | reversed
//                         MaturesAt time.Time
//                         CreatedAt }
//   func NewEarning(suborderID, sellerID uuid.UUID, gross, commission money.Amount,
//                   holdPeriod time.Duration, now time.Time) *Earning
//   func (e *Earning) Mature(now time.Time) error
//   func (e *Earning) Reverse(reason string) error
//
// MaturesAt = waktu suborder selesai + masa tahan. Masa tahan ada supaya ada
// jendela waktu untuk sengketa dan pengembalian sebelum uang keluar.
//
// Reverse dipakai saat refund terjadi SETELAH dana matang atau bahkan sudah
// cair. Itu bukan kasus langka, dan sistem harus bisa menanganinya: hasilnya
// jadi saldo negatif atau piutang ke penjual, bukan error.
