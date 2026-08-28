package app

// TODO — satu handler per peristiwa yang layak dicatat:
//   OnOrderPlaced, OnOrderCancelled, OnSuborderRejected,
//   OnPaymentSettled, OnRefundCompleted, OnPayoutCompleted,
//   OnStockAdjusted, OnPriceChanged,
//   OnSellerVerified, OnSellerSuspended, OnVoucherRedeemed
//
// IDEMPOTEN. Simpan event id dengan unique index, INSERT ... ON CONFLICT DO NOTHING.
//
// Payload rusak -> log dan return nil. Mengembalikan error untuk sesuatu yang
// tidak akan sembuh dengan retry membuat relay mencoba selamanya.
