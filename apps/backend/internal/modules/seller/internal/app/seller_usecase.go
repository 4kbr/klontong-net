package app

// TODO:
//   RegisterSeller(ctx, RegisterSellerInput) (*domain.Seller, error)
//     WithinTx: buat seller (pending) + tambahkan pemilik sebagai owner +
//     berikan peran `seller` lewat identity.Port + outbox event
//   GetSeller / GetPublicSeller / UpdateSeller / SetPayoutAccount (owner saja)
//   ListMySellers
