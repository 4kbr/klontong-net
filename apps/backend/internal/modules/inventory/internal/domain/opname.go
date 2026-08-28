package domain

// TODO:
//   type Opname struct { ID; OutletID; Status; StartedBy; StartedAt; FinishedAt; Note }
//   type OpnameItem struct { ID; OpnameID; VariantID;
//                            SystemQuantity, CountedQuantity, Difference decimal }
//   func (o *Opname) Finish(items []OpnameItem) ([]Movement, error)
//
// Hasil opname masuk sebagai Movement kind=adjustment, BUKAN mengubah
// QuantityOnHand langsung. Dengan begitu selisihnya tercatat beserta waktunya
// dan siapa yang menghitung.
