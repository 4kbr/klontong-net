package app

// TODO: ListAddresses, GetAddress, CreateAddress, UpdateAddress, DeleteAddress,
//       SetDefaultAddress.
//
// SetDefaultAddress harus dalam satu transaksi: matikan is_default yang lama,
// nyalakan yang baru. Dua request bersamaan bisa menghasilkan dua alamat default
// kalau tidak dikunci — pertimbangkan unique index parsial
// `unique (user_id) where is_default` sebagai penjaga.
//
// DeleteAddress pada alamat default harus memindahkan default ke alamat lain,
// atau menolak kalau itu satu-satunya.
