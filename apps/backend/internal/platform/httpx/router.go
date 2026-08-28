package httpx

// Router root dan pembagian area.
//
// TODO — func NewRouter(deps RouterDeps) *chi.Mux
//
// Middleware global: RequestID, RealIP, Recoverer, Logger, Timeout, CORS.
//
// EMPAT AREA dengan model akses berbeda:
//
//   /api/v1/*           Pembeli. Sebagian PUBLIK (katalog, pencarian, detail
//                       produk) dan sebagian butuh login (keranjang, checkout,
//                       pesanan). Jangan memasang middleware auth untuk seluruh
//                       area ini — halaman produk harus bisa dibuka tanpa login,
//                       dan itu penting untuk SEO dan konversi.
//
//   /api/v1/seller/*    Dasbor penjual. Wajib login + peran seller + keanggotaan
//                       pada toko yang diakses.
//
//   /api/v1/admin/*     Panel marketplace. Wajib peran admin.
//
//   /webhook/*          Payment gateway dan kurir. Verifikasi SIGNATURE,
//                       TANPA sesi, TANPA CORS. Rate limit longgar tapi ada.
//
// Plus /healthz dan /readyz.
//
// Kesalahan yang mahal: memasang satu middleware auth lalu bercabang di
// dalamnya. Pisahkan sejak router.
