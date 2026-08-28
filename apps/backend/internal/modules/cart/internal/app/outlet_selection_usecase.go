package app

// Pemilihan outlet untuk setiap baris keranjang.
//
// Satu penjual bisa punya beberapa outlet dengan stok berbeda. Baris keranjang
// harus terikat ke satu outlet, dan pilihannya menentukan ongkir dan
// ketersediaan.
//
// TODO — SuggestOutlet(ctx, sellerID, variantID uuid.UUID, qty decimal,
//                      buyerLat, buyerLng *float64) (uuid.UUID, error):
//   1. inventory.OutletsWithStock -> outlet mana saja yang stoknya cukup
//   2. saring yang aktif dan mendukung setidaknya satu metode kirim
//   3. kalau pembeli punya koordinat, pilih yang TERDEKAT
//   4. kalau tidak, pilih outlet default penjual
//   5. tidak ada yang memenuhi -> ErrInsufficientStock dengan pesan yang jelas
//
// TODO — ReassignOutlets(ctx, cartID): dipanggil saat alamat pengiriman berubah.
//   Outlet terdekat dari alamat A belum tentu terdekat dari alamat B, dan
//   ongkirnya bisa jauh berbeda.
//
// KEPUTUSAN YANG PERLU DIAMBIL: apakah satu suborder boleh mengambil dari
// BEBERAPA outlet penjual yang sama? Itu memungkinkan pesanan yang tidak muat di
// satu outlet tetap jalan, tapi berarti satu suborder punya beberapa pengiriman.
// Rekomendasi: TIDAK untuk sekarang — satu suborder satu outlet. Catat sebagai
// ADR kalau nanti berubah.
