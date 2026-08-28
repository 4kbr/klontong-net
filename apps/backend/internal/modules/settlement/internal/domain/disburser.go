package domain

// TODO:
//   type Disburser interface {
//       Transfer(ctx, TransferRequest) (reference string, err error)
//       Inquiry(ctx, reference string) (status string, err error)
//       ValidateAccount(ctx, bankCode, accountNumber string) (accountName string, err error)
//   }
//
// ValidateAccount dipanggil saat penjual menyetel rekening, bukan saat mencairkan.
// Rekening salah yang baru ketahuan saat transfer gagal berarti dana menggantung
// dan penjual menunggu tanpa penjelasan.
//
// Transfer WAJIB membawa idempotency key ke penyedia. Timeout tidak berarti
// transfer gagal — hanya berarti kita tidak tahu, dan mengulanginya tanpa kunci
// bisa berarti mengirim dua kali.
