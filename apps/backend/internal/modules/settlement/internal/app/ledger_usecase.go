package app

// Pencatatan buku besar. Semua dipanggil dari handler event, bukan dari HTTP.
//
// TODO — RecordOrderPaid(ctx, orderID): saat pembayaran gateway masuk
//   Debit  marketplace/clearing        sebesar total pembayaran
//   Credit seller/<id>/pending         sebesar pendapatan tiap penjual
//   Credit marketplace/revenue         sebesar komisi
//   (jumlah kredit = jumlah debit)
//
// TODO — RecordEarningMatured(ctx, earningID):
//   Debit  seller/<id>/pending
//   Credit seller/<id>/available
//
// TODO — RecordPayout(ctx, payoutID):
//   Debit  seller/<id>/available
//   Credit marketplace/clearing
//
// TODO — RecordRefund(ctx, refundID): arah dibalik.
//   Kalau dana sudah cair, hasilnya saldo negatif atau piutang ke penjual.
//   JANGAN menolak refund hanya karena dananya sudah keluar — itu keputusan
//   bisnis yang harus tetap bisa dijalankan.
//
// TODO — RecordCODCollected(ctx, orderID): lihat ADR-012, alurnya berbeda.
//
// Setiap fungsi membuat SATU jurnal seimbang. Validasi keseimbangannya sebelum
// menyimpan, dan gagalkan transaksi kalau tidak seimbang.
