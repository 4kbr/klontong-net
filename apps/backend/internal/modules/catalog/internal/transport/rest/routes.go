package rest

// TODO:
//   Publik:
//     GET /api/v1/categories
//     GET /api/v1/products                 (browse + filter + urut)
//     GET /api/v1/products/search
//     GET /api/v1/products/{slug}
//
//   Dasbor penjual:
//     GET/POST /api/v1/seller/products
//     GET/PATCH/DELETE /api/v1/seller/products/{productID}
//     POST /api/v1/seller/products/{productID}/publish
//     GET/POST /api/v1/seller/products/{productID}/variants
//     PATCH/DELETE /api/v1/seller/variants/{variantID}
//     POST /api/v1/seller/products/{productID}/images/upload-url
//
//   Admin:
//     GET/POST/PATCH /api/v1/admin/categories
//
// Rute publik memakai OptionalAuth, bukan Authenticate. Halaman produk harus
// bisa dibuka tanpa login.
