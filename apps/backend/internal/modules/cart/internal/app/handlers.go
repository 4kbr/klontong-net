package app

// TODO — handler event:
//   OnProductArchived  -> tandai baris keranjang yang memuatnya
//   OnSellerSuspended  -> tandai seluruh baris dari penjual itu
//   OnPriceChanged     -> kosongkan cache harga baris terkait
//
// Menandai, BUKAN menghapus. Menghapus barang dari keranjang orang diam-diam
// membuat pembeli bingung mencari barang yang tadi ada. Tandai, tampilkan
// alasannya, biarkan pembeli yang memutuskan.
