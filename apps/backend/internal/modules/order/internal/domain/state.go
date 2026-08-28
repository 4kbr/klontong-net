package domain

// Peta transisi status. Ditulis eksplisit sebagai data, bukan tersebar sebagai
// `if` di berbagai usecase.
//
// TODO:
//   var suborderTransitions = map[SuborderStatus][]SuborderStatus{
//       SubAwaitingConfirmation: {SubConfirmed, SubRejected, SubCancelled},
//       SubConfirmed:            {SubPacking, SubCancelled},
//       SubPacking:              {SubShipped, SubReadyForPickup, SubCancelled},
//       SubReadyForPickup:       {SubDelivered, SubCancelled},
//       SubShipped:              {SubDelivered},
//       SubDelivered:            {SubCompleted},
//       SubCompleted:            {},
//       SubCancelled:            {},
//       SubRejected:             {},
//   }
//   func CanTransition(from, to SuborderStatus) bool
//   func TerminalStatuses() []SuborderStatus
//
// Perhatikan: dari SubShipped TIDAK BISA dibatalkan. Barang sudah di jalan;
// yang bisa terjadi hanya pengembalian, dan itu alur terpisah.
//
// SubDelivered -> SubCompleted terjadi otomatis setelah masa komplain lewat,
// atau segera kalau pembeli menekan "pesanan diterima". Jangan menyamakan
// keduanya — `delivered` berarti barang sampai, `completed` berarti tidak ada
// lagi sengketa dan dana boleh cair.
//
// Menulis peta ini sebagai data membuat satu tempat untuk diperiksa saat
// bertanya "kenapa pesanan ini tidak bisa dibatalkan".
