package domain

// TODO:
//   type PayoutStatus string   // requested | processing | paid | failed
//   type Payout struct { ID; SellerID; Amount money.Amount; Status
//                        BankCode, AccountNumber, AccountName string
//                        ProviderReference string
//                        IdempotencyKey string
//                        RequestedAt, ProcessedAt *time.Time
//                        FailedReason string
//                        Items []PayoutItem }
//   type PayoutItem struct { PayoutID, SuborderID uuid.UUID; Amount money.Amount }
//   func (p *Payout) MarkPaid(ref string, now time.Time) error
//   func (p *Payout) MarkFailed(reason string) error
//
// IdempotencyKey WAJIB unique di database. Pencairan yang berjalan dua kali
// karena worker restart berarti uang keluar dua kali, dan menariknya kembali
// jauh lebih sulit daripada mencegahnya.
//
// Kunci dibuat SEKALI saat payout dibuat, bukan setiap percobaan kirim.
