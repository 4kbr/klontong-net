package domain

// TODO:
//   type Document struct { ID; SellerID; Kind; StorageKey; Status;
//                          ReviewedBy; ReviewedAt; RejectionReason; CreatedAt }
//
// Dokumen verifikasi (KTP, NPWP) TIDAK BOLEH bisa diakses publik. Selalu lewat
// presigned URL berumur pendek, dan hanya untuk admin.
