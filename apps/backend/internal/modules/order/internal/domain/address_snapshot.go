package domain

// TODO:
//   type AddressSnapshot struct { RecipientName, RecipientPhone,
//                                 Province, City, District, Village, PostalCode,
//                                 AddressLine, Notes string
//                                 Latitude, Longitude *float64 }
//   func FromCustomerAddress(a customer.Address) AddressSnapshot
//
// Disimpan sebagai jsonb di pesanan. Pembeli yang mengubah atau menghapus
// alamatnya tidak boleh mengubah tujuan pengiriman pesanan yang sudah jalan.
