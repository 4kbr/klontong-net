package app

// TODO — handler event:
//   OnPaymentSettled  -> order.MarkPaid, seluruh suborder ke awaiting_confirmation,
//                        outbox EventOrderPaid
//   OnPaymentFailed   -> order.Cancel, lepas stok dan voucher
//   OnPaymentExpired  -> sama dengan di atas
//   OnShipmentDelivered -> suborder.Transition(SubDelivered)
//
// OnPaymentSettled adalah handler paling penting di sistem dan WAJIB idempoten.
// Webhook gateway datang berkali-kali; memproses dua kali tidak boleh membuat
// dua notifikasi ke penjual atau dua catatan pembukuan.
