package domain

// TODO:
//   type OwnerType string   // seller | marketplace | buyer | gateway | courier
//   type AccountKind string // pending | available | payable | revenue | clearing
//   type Account struct { ID; OwnerType; OwnerID *uuid.UUID; Kind; Currency }
//
// Akun yang minimal harus ada:
//   seller/<id>/pending     dana penjual yang belum matang
//   seller/<id>/available   dana penjual yang siap dicairkan
//   marketplace/revenue     komisi kita
//   marketplace/clearing    kas yang masuk dari gateway
//   courier/cod_receivable  piutang dari kurir untuk pesanan COD
