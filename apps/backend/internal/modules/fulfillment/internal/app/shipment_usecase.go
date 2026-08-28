package app

// TODO:
//   CreateShipment(ctx, CreateShipmentInput) — dipanggil saat penjual menekan kirim
//     courier        -> CourierProvider.CreateBooking, simpan nomor resi
//     local_delivery -> catat nama & nomor kurir penjual
//     pickup         -> terbitkan PickupCode dengan masa berlaku
//   UpdateShipmentStatus
//   ConfirmPickup(ctx, shipmentID, code string) — verifikasi kode, tandai diambil
//   UploadDeliveryProof
//
// ConfirmPickup: bandingkan kode dengan perbandingan waktu-konstan. Kode ambil
// pendek dan mudah ditebak kalau ada yang mencoba berulang — beri rate limit
// per shipment.
