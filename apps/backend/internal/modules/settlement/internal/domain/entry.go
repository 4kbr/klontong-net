package domain

// Buku besar double entry. Baca ADR-011 sebelum menyentuh file ini.
//
// TODO:
//   type Direction string   // debit | credit
//   type Entry struct { ID; JournalID uuid.UUID; AccountID; Direction;
//                       Amount money.Amount; ReferenceType string;
//                       ReferenceID uuid.UUID; Description string;
//                       OccurredAt; CreatedAt }
//   type Journal struct { ID uuid.UUID; Entries []Entry; Description string }
//   func (j *Journal) IsBalanced() bool
//         jumlah debit HARUS sama dengan jumlah kredit
//   func (j *Journal) Validate() error
//
// TIDAK ADA kolom saldo yang di-UPDATE. Saldo adalah hasil penjumlahan entri.
//
// Kenapa: kolom saldo yang diperbarui langsung akan menyimpang cepat atau
// lambat — satu bug, satu race, satu proses yang mati di tengah. Dan saat itu
// terjadi kamu tidak punya cara tahu angka mana yang benar. Dengan buku besar,
// saldo selalu bisa dihitung ulang dan setiap perubahan punya jejaknya.
//
// Kalau penjumlahan jadi lambat, buat snapshot saldo berkala. Snapshot tetap
// bisa dihitung ulang dari entri; kolom saldo tidak.
//
// Setiap jurnal WAJIB seimbang. Tegakkan di kode DAN pertimbangkan constraint
// deferred di database.
