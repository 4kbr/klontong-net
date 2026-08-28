package domain

// TODO:
//   type Price struct { ID; VariantID; OutletID *uuid.UUID;
//                       Amount money.Amount
//                       CompareAtAmount money.Amount
//                       CostAmount money.Amount
//                       Currency string
//                       StartsAt, EndsAt *time.Time
//                       IsActive bool
//                       Tiers []QuantityTier
//                       CreatedAt; UpdatedAt }
//   func NewPrice(variantID uuid.UUID, amount money.Amount) (*Price, error)
//         Amount harus > 0
//         CompareAtAmount kalau diisi harus > Amount, kalau tidak itu
//         harga coret palsu dan itu menyesatkan pembeli
//   func (p *Price) IsEffectiveAt(t time.Time) bool
//
// OutletID nil berarti berlaku untuk semua outlet. Sebagian besar penjual
// memakai satu harga; yang punya outlet di lokasi berbeda kadang perlu berbeda.
//
// CostAmount (modal) hanya untuk laporan margin penjual. TIDAK BOLEH pernah
// keluar di response yang bisa dilihat pembeli.
