package app

// TODO — Resolve / ResolveMany:
//   1. ambil harga efektif: cari yang khusus outlet dulu, kalau tidak ada pakai
//      yang umum (OutletID nil)
//   2. saring yang IsEffectiveAt(now) — harga terjadwal belum berlaku
//   3. terapkan tier berdasarkan kuantitas
//   4. kembalikan Resolved
//
// Ini fungsi yang menentukan angka di layar pembeli DAN angka di tagihan.
// Keduanya harus memakai jalur kode yang SAMA — kalau checkout menghitung harga
// dengan cara berbeda dari halaman produk, cepat atau lambat keduanya akan
// berbeda dan pembeli yang menemukannya. Lihat ADR-004.
