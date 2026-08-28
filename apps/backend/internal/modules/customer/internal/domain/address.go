package domain

// TODO:
//   type Address struct { ID; UserID; Label; RecipientName; RecipientPhone;
//                         Province; City; District; Village; PostalCode;
//                         AddressLine; Notes; Latitude, Longitude *float64;
//                         IsDefault bool; DeletedAt; CreatedAt; UpdatedAt }
//   func NewAddress(...) (*Address, error)
//       validasi: nama & nomor penerima wajib, kota & kode pos wajib,
//       nomor HP dinormalisasi ke E.164
//   func (a *Address) HasCoordinates() bool
//   func (a *Address) SoftDelete(now)
//
// Alamat TIDAK PERNAH dihapus permanen. Pesanan lama merujuk padanya, dan
// riwayat tidak boleh berubah karena pembeli merapikan daftar alamatnya.
// Meski begitu, pesanan tetap menyimpan SALINAN alamat. Lihat ADR-010.
//
// Batasi jumlah alamat per pengguna (mis. 20). Tanpa batas, ada saja yang
// mengisinya sampai ribuan.
