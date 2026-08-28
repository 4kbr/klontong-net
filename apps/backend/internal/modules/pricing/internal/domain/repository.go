package domain

// TODO:
//   type PriceRepository interface {
//       Upsert/FindByVariant/FindByVariantAndOutlet/
//       FindManyEffective(ctx, pairs []VariantOutlet, at time.Time) (map[...]*Price, error)/
//       Deactivate
//   }
//   type TierRepository interface { ReplaceForPrice/ListByPrice/ListByPrices }
//
// FindManyEffective harus mengambil harga khusus outlet DAN harga umum dalam
// satu query, lalu resolusinya dilakukan di layer app. Dua query terpisah
// berarti dua kali perjalanan ke database untuk setiap pembukaan keranjang.
