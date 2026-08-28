package domain

// TODO:
//   type Kind string        // percentage | fixed_amount | free_shipping
//   type OwnerType string   // marketplace | seller
//   type AppliesTo string   // items | shipping
//   type Voucher struct {
//       ID; Code; Name; Description
//       OwnerType; OwnerSellerID *uuid.UUID
//       Kind; ValueBPS int; ValueAmount money.Amount
//       MaxDiscountAmount, MinOrderAmount money.Amount
//       AppliesTo
//       Quota, QuotaPerUser, UsedCount int
//       StartsAt, EndsAt *time.Time
//       IsActive bool
//   }
//   func (v *Voucher) IsUsableAt(t time.Time) bool
//   func (v *Voucher) HasQuota() bool
//   func (v *Voucher) CalculateDiscount(base money.Amount) money.Amount
//         percentage: base * ValueBPS / 10000, dibatasi MaxDiscountAmount,
//                     dibulatkan KE BAWAH
//         fixed:      ValueAmount, tapi tidak boleh melebihi base
//
// Persentase disimpan dalam BASIS POIN, bukan float. Diskon 12,5% adalah 1250.
// Tidak ada pecahan yang perlu disimpan dan tidak ada float di jalur uang.
//
// MaxDiscountAmount wajib untuk voucher persentase. Tanpa batas, "diskon 50%"
// pada pesanan Rp5 juta berarti kerugian Rp2,5 juta dari satu transaksi.
