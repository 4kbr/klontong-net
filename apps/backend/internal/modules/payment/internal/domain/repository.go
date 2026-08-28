package domain

// TODO:
//   type PaymentRepository interface { Create/Update/FindByID/FindByOrder/
//                                      FindByProviderRef/ListPending(ctx, olderThan) }
//   type WebhookRepository interface { Record(ctx, *WebhookEvent) (isDuplicate bool, err error)
//                                      MarkProcessed/MarkFailed }
//   type RefundRepository interface { Create/Update/ListByPayment/SumByPayment }
//
// Record mengembalikan penanda duplikat alih-alih error, karena webhook ganda
// adalah kondisi NORMAL, bukan kesalahan.
