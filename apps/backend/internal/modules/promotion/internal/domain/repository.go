package domain

// TODO:
//   type VoucherRepository interface { Create/Update/FindByID/FindByCode/
//                                      FindByCodeForUpdate/ListActive/
//                                      IncrementUsed/DecrementUsed }
//   type RedemptionRepository interface { Create/CountByVoucherAndUser/
//                                         ListByOrder/DeleteByOrder }
//
// FindByCodeForUpdate memakai SELECT ... FOR UPDATE. Itu yang membuat
// pemotongan kuota aman terhadap konkurensi.
