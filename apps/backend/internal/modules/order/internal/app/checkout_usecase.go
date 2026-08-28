package app

// Usecase paling kompleks di seluruh project. Kerjakan setelah semua modul
// pendukungnya jalan, dan tulis testnya lebih dulu.
//
// ALUR DUA TAHAP. Tahap pertama tidak mengikat apa pun:
//
// TODO — PreviewCheckout(ctx, PreviewInput) (CheckoutPreview, error):
//   1. ambil keranjang lewat cart.Port
//   2. ambil alamat lewat customer.Port
//   3. kelompokkan per penjual -> calon suborder
//   4. untuk tiap kelompok: tentukan outlet, hitung harga terkini
//      (pricing.ResolveMany), periksa ketersediaan (inventory.AvailableMany)
//   5. minta opsi ongkir per kelompok lewat fulfillment.Port
//   6. terapkan voucher lewat promotion.Port, BAGI diskonnya ke tiap suborder
//      dengan money.Distribute
//   7. hitung total dan kembalikan rincian LENGKAP per suborder
//
//   Preview boleh dipanggil berkali-kali. Ia tidak menahan stok dan tidak
//   membuat apa pun.
//
// TODO — PlaceOrder(ctx, PlaceOrderInput) (*domain.Order, error):
//   Wajib membawa Idempotency-Key. Lihat ADR-008.
//
//   1. HITUNG ULANG SEMUANYA dari awal, jangan percaya angka dari klien.
//      Klien mengirim pilihan (alamat mana, kurir mana, voucher apa) — bukan
//      harga dan bukan total. Lihat ADR-004.
//   2. Bandingkan dengan angka yang ditampilkan klien; kalau berbeda,
//      kembalikan ErrPriceChanged dengan rinciannya alih-alih diam-diam
//      menagih jumlah yang berbeda dari yang dilihat pembeli.
//   3. WithinTx:
//        a. inventory.Reserve untuk SEMUA item, semua suborder
//           Gagal di sini -> seluruh checkout batal. Ini alasan utama kenapa
//           ia harus di dalam transaksi yang sama.
//        b. promotion.Redeem — kuota voucher dipotong di sini, bukan sebelumnya
//        c. buat Order + Suborder + Item, semuanya snapshot
//        d. bekukan CommissionBPS tiap suborder dari seller.Port
//        e. payment.CreatePayment (gateway) atau tandai COD
//        f. cart.MarkConverted
//        g. outbox: EventOrderPlaced
//      COMMIT
//   4. kembalikan pesanan beserta instruksi pembayaran
//
//   TIDAK ADA panggilan jaringan di dalam transaksi. Pembuatan transaksi di
//   payment gateway dilakukan SETELAH commit, atau dilakukan sebelum transaksi
//   dengan status pending dan dikaitkan setelahnya. Memanggil gateway di dalam
//   transaksi berarti menahan koneksi database selama pihak lain berpikir, dan
//   rollback tidak membatalkan transaksi yang sudah terlanjur dibuat di sana.
//
//   Kunci baris stok dalam urutan konsisten (inventory yang mengurusnya).
//   Dua pembeli yang checkout barang sama persis akan bertemu di sini.
