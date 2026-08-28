package rest

// TODO:
//   Publik (pembeli):
//     GET  /api/v1/sellers/{slug}
//     GET  /api/v1/sellers/{slug}/products
//
//   Dasbor penjual:
//     POST /api/v1/seller/register
//     GET  /api/v1/seller/me
//     PATCH /api/v1/seller/me
//     PUT  /api/v1/seller/me/payout-account
//     GET/POST /api/v1/seller/outlets
//     PATCH/DELETE /api/v1/seller/outlets/{outletID}
//     GET/POST /api/v1/seller/members
//     PATCH/DELETE /api/v1/seller/members/{userID}
//     POST /api/v1/seller/documents
//
//   Admin:
//     GET  /api/v1/admin/sellers
//     POST /api/v1/admin/sellers/{sellerID}/approve
//     POST /api/v1/admin/sellers/{sellerID}/reject
//     POST /api/v1/admin/sellers/{sellerID}/suspend
