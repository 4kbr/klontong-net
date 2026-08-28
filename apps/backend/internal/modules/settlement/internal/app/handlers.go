package app

// TODO — handler event, semuanya harus IDEMPOTEN:
//   OnOrderPaid          -> buat Earning per suborder + RecordOrderPaid
//   OnSuborderCompleted  -> tetapkan MaturesAt earning
//   OnSuborderRejected   -> batalkan earning suborder itu
//   OnRefundCompleted    -> RecordRefund, reverse earning bila perlu
//   OnCODCollected       -> RecordCODCollected
//
// Idempotensi di sini bukan soal kerapian. Event yang diproses dua kali berarti
// penjual dibayar dua kali. Cara paling aman: unique index pada
// (event_id) di tabel jurnal, dan INSERT ... ON CONFLICT DO NOTHING.
