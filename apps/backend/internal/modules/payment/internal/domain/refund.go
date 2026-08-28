package domain

// TODO:
//   type Refund struct { ID; PaymentID; OrderID; SuborderID *uuid.UUID;
//                        Amount money.Amount; Reason string; Status;
//                        ProviderReference string; RequestedBy uuid.UUID;
//                        ProcessedAt *time.Time; CreatedAt }
//   func (r *Refund) IsPartial(paymentAmount money.Amount) bool
//
// Refund menunjuk ke SUBORDER, bukan hanya order. Satu penjual menolak,
// dua penjual lain tetap berjalan — yang dikembalikan hanya bagian penjual
// yang menolak, ditambah ongkirnya.
//
// Jumlah seluruh refund tidak boleh melebihi jumlah pembayaran. Periksa itu
// di dalam transaksi, bukan dengan penjumlahan yang dibaca sebelumnya.
