# Frontend

Belum dikerjakan. Folder ini sengaja dibiarkan kosong sampai kontrak API backend
cukup stabil.

Saat mulai, yang perlu diputuskan lebih dulu:

- **Berapa aplikasi.** Setidaknya ada tiga permukaan yang berbeda jauh: toko untuk
  pembeli, dasbor untuk penjual, dan panel admin marketplace. Menyatukan ketiganya
  dalam satu aplikasi terdengar hemat tapi biasanya berakhir dengan bundle besar
  yang penuh percabangan peran.
- **Rendering.** Halaman produk butuh SEO, dasbor penjual tidak. Ini argumen kuat
  untuk memisahkan storefront dari dasbor.
- **Tipe API.** Jangan mengetik ulang bentuk response. Generate dari spesifikasi
  OpenAPI backend.
- **Harga dan total.** Frontend **tidak pernah** menghitung total yang mengikat.
  Ia boleh menampilkan perkiraan, tapi angka yang dibayar selalu datang dari
  server. Baca ADR-004.
