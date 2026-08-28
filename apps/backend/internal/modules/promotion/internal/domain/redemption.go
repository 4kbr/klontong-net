package domain

// TODO:
//   type Redemption struct { ID; VoucherID; UserID; OrderID; SuborderID;
//                            DiscountAmount money.Amount; CreatedAt }
//
// Kuota ditegakkan dengan unique index (voucher_id, order_id) DAN pengecekan
// UsedCount di dalam transaksi dengan SELECT ... FOR UPDATE pada baris voucher.
// Mengecek kuota sebelum transaksi lalu memakainya sesudahnya adalah balapan
// yang pasti kalah saat flash sale.
