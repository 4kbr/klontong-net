package domain

// TODO:
//   type AccountRepository interface { EnsureExists/FindByOwner/Balance(ctx, accountID) }
//   type EntryRepository interface { CreateJournal(ctx, *Journal) error
//                                    ListByAccount/SumByAccount }
//   type EarningRepository interface { Create/FindBySuborder/ListMatured(ctx, before, limit)/
//                                      UpdateStatus/ListBySeller }
//   type PayoutRepository interface { Create/Update/FindByID/ListPending/
//                                     ListBySeller/ExistsInProgress }
//
// SumByAccount adalah penjumlahan seluruh entri. Kalau jadi lambat, tambahkan
// snapshot — jangan mengganti dengan kolom saldo.
//
// ExistsInProgress mencegah dua pencairan berjalan bersamaan untuk penjual yang
// sama.
