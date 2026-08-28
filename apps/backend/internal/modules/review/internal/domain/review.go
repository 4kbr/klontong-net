package domain

// TODO:
//   type Review struct { ID; OrderItemID; ProductID; VariantID; SellerID; UserID
//                        Rating int; Comment string; Images []string
//                        IsAnonymous bool; Status string
//                        CreatedAt; UpdatedAt }
//   func NewReview(orderItemID uuid.UUID, rating int, comment string) (*Review, error)
//         rating 1..5, komentar boleh kosong
//   func (r *Review) Edit(rating int, comment string, within time.Duration, now time.Time) error
//
// Terikat ke ORDER ITEM dan unique. Ini yang mencegah ulasan palsu: hanya yang
// benar-benar membeli DAN menerima yang bisa menulis, dan hanya sekali per
// barang yang dibeli.
//
// Verifikasi itu terjadi lewat order.Port.HasPurchased, bukan dengan mempercayai
// klien mengirim order_item_id yang benar.
