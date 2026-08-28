package domain

// TODO:
//   type CategoryRepository interface { Create/Update/FindByID/FindBySlug/ListTree/ListActive }
//   type ProductRepository interface {
//       Create/Update/FindByID/FindBySlug/ListBySeller/
//       Search(ctx, SearchQuery) ([]*Product, string, error)
//       UpdateAggregates(ctx, productID uuid.UUID, ratingAvg float64, ratingCount, soldCount int) error
//   }
//   type VariantRepository interface {
//       Create/Update/FindByID/FindManyByID/ListByProduct/CountActive
//   }
//   type ImageRepository interface { Create/ListByProduct/Delete/Reorder }
//
// Search adalah query terberat di sistem. Mulai dengan ILIKE + index trigram;
// pindah ke tsvector kalau sudah tidak cukup. Jangan langsung memasang mesin
// pencari terpisah — itu menambah satu sistem yang harus disinkronkan.
