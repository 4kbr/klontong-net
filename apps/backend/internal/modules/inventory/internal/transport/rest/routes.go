package rest

// TODO (dasbor penjual saja; pembeli tidak pernah mengakses stok langsung):
//   GET   /api/v1/seller/outlets/{outletID}/stock
//   PATCH /api/v1/seller/outlets/{outletID}/stock/{variantID}
//   POST  /api/v1/seller/stock/transfer
//   GET   /api/v1/seller/stock/movements
//   GET/POST /api/v1/seller/outlets/{outletID}/opnames
//   POST  /api/v1/seller/opnames/{opnameID}/finish
//
// Ketersediaan untuk pembeli disajikan lewat modul catalog, bukan dari sini.
