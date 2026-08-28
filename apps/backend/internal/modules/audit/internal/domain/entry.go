package domain

// TODO:
//   type ActorType string   // user | seller | admin | system | gateway | courier
//   type Entry struct { ID; ActorType; ActorID *uuid.UUID; ActorLabel string
//                       Action string; TargetType string; TargetID uuid.UUID
//                       Before, After map[string]any
//                       IP, UserAgent string; CreatedAt }
//
// Append-only. Tidak ada method Update.
//
// `Before` dan `After` menyimpan keadaan sebelum dan sesudah. Untuk perubahan
// harga dan stok, inilah yang menjawab "kok berubah" — dan pertanyaan itu pasti
// datang.
//
// JANGAN menyimpan password hash, token, atau server key di Before/After.
// Saring field sensitif sebelum menulis.
