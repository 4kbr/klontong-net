package infra

// TODO: implementasi domain.Gateway untuk Midtrans.
//   - Charge lewat Core API atau Snap, tergantung pilihan
//   - VerifyWebhook: signature Midtrans adalah SHA512 dari
//     order_id + status_code + gross_amount + server_key.
//     Bandingkan dengan perbandingan waktu-konstan.
//   - Petakan status mereka ke status kita: capture/settlement -> settled,
//     pending -> pending, deny/cancel/expire -> failed/expired
//
// Panggil REST langsung dengan net/http, jangan tarik SDK. Permukaan yang kita
// pakai sempit, dan SDK mengikat kita ke bentuk mereka. Lihat ADR-014.
