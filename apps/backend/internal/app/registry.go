package app

// Composition root. SATU-SATUNYA file yang mengimpor semua modul sekaligus.
//
// TODO — func NewRegistry(deps Deps) (*Registry, error)
//
//   Urutan pembangunan mengikuti arah dependensi:
//
//     identityMod   := identity.New(...)
//     customerMod   := customer.New(customer.Config{ Users: identityMod.Port() })
//     sellerMod     := seller.New(seller.Config{ Users: identityMod.Port() })
//     catalogMod    := catalog.New(catalog.Config{ Sellers: sellerMod.Port() })
//     inventoryMod  := inventory.New(inventory.Config{ Catalog: catalogMod.Port() })
//     pricingMod    := pricing.New(pricing.Config{ Catalog: catalogMod.Port() })
//     promotionMod  := promotion.New(...)
//     fulfillmentMod:= fulfillment.New(fulfillment.Config{ Sellers: sellerMod.Port() })
//     cartMod       := cart.New(cart.Config{
//                          Catalog: catalogMod.Port(), Pricing: pricingMod.Port(),
//                          Inventory: inventoryMod.Port(), Sellers: sellerMod.Port(),
//                      })
//     paymentMod    := payment.New(...)
//     orderMod      := order.New(order.Config{
//                          Cart: cartMod.Port(), Catalog: catalogMod.Port(),
//                          Pricing: pricingMod.Port(), Inventory: inventoryMod.Port(),
//                          Promotion: promotionMod.Port(), Fulfillment: fulfillmentMod.Port(),
//                          Payment: paymentMod.Port(), Customers: customerMod.Port(),
//                      })
//     settlementMod := settlement.New(settlement.Config{ Sellers: sellerMod.Port() })
//     reviewMod     := review.New(review.Config{ Orders: orderMod.Port() })
//     notificationMod := notification.New(...)
//     auditMod      := audit.New(...)
//
//   `order` bergantung ke paling banyak modul, dan itu wajar — checkout memang
//   titik temu seluruh sistem. Yang harus dijaga: modul lain TIDAK boleh
//   bergantung balik ke `order` kecuali lewat event. Kalau muncul kebutuhan itu,
//   berhenti dan pikirkan ulang.
//
//   Lalu SEMUA langganan event terkumpul di sini:
//     auditMod.RegisterSubscriptions(deps.Bus)
//     notificationMod.RegisterSubscriptions(deps.Bus)
//     inventoryMod.RegisterSubscriptions(deps.Bus)
//     settlementMod.RegisterSubscriptions(deps.Bus)
//     catalogMod.RegisterSubscriptions(deps.Bus)   // agregat rating & sold_count
//
// TODO — func (r *Registry) MountRoutes(mux *chi.Mux)
//   mux.Route("/api/v1", ...)          rute pembeli
//   mux.Route("/api/v1/seller", ...)   dasbor penjual
//   mux.Route("/api/v1/admin", ...)    panel marketplace
//   mux.Route("/webhook", ...)         payment & kurir
//
// TODO — func (r *Registry) RegisterWorkers(runner *Runner)
// TODO — func (r *Registry) Close() error
//
// Registry merakit, tidak memutuskan.
