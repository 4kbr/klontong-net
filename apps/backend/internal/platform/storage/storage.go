package storage

// TODO:
//   type Storage interface {
//       Put(ctx, key string, r io.Reader, size int64, mime string) error
//       PresignGet(ctx, key string, ttl time.Duration) (string, error)
//       PresignPut(ctx, key string, ttl time.Duration) (string, error)
//       Delete(ctx, key string) error
//   }
//   func NewS3(cfg config.StorageConfig) (Storage, error)
//
// PresignPut ada supaya unggah foto produk bisa langsung dari browser ke object
// storage tanpa melewati API kita. Foto produk banyak dan besar; menyalurkannya
// lewat API berarti membayar bandwidth dua kali dan menahan goroutine lama.
//
// Key berpola: `products/<seller_id>/<product_id>/<uuid>.<ext>`,
//              `documents/<seller_id>/<uuid>.<ext>`  (dokumen verifikasi, privat)
//
// Dokumen verifikasi penjual TIDAK BOLEH bisa diakses publik. Foto KTP yang
// bocor adalah insiden, bukan bug kecil.
