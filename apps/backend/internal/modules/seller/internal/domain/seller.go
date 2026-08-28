package domain

// TODO:
//   type Seller struct { ID; OwnerUserID; Name; Slug; Description; LogoURL;
//                        BannerURL; Status; CommissionBPS *int;
//                        PayoutBankCode, PayoutAccountNumber, PayoutAccountName string;
//                        VerifiedAt; CreatedAt; UpdatedAt }
//   func NewSeller(ownerUserID uuid.UUID, name string) (*Seller, error)
//         slug dari nama + suffix acak agar unik
//   func (s *Seller) Verify(now time.Time) error   // hanya dari pending
//   func (s *Seller) Suspend(reason string) error
//   func (s *Seller) CanSell() bool                // hanya verified
//   func (s *Seller) SetPayoutAccount(...) error
//
// CommissionBPS bertipe pointer: nil berarti pakai default marketplace.
// Membedakan "belum diatur" dari "diatur nol" penting — komisi nol persen itu
// keputusan bisnis yang sah dan tidak boleh tertukar dengan belum diisi.
