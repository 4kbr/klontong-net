# AGENTS.md — backend

Aturan spesifik Go. **Baca [`../../AGENTS.md`](../../AGENTS.md) lebih dulu** —
di sana ada aturan yang berlaku untuk seluruh repo, termasuk aturan uang dan stok
yang paling penting di project ini.

File ini hanya memuat yang khas backend.

---

## Layer

```
transport ──▶ app ──▶ domain ◀── infra
```

- `domain` tidak boleh mengimpor `net/http`, `pgx`, atau `chi`. Ditegakkan
  `depguard`.
- Interface repository didefinisikan di `domain/repository.go`, diimplementasikan
  di `infra/`.
- `transport` tidak boleh memanggil `infra` langsung.
- `math/big` dilarang di luar `platform/money`, supaya jalur uang tidak diam-diam
  berubah.

---

## Repository

Setiap method dibuka dengan:

```go
db := postgres.ExecutorFrom(ctx, r.pool)
```

dan ditutup dengan `postgres.Translate(err)`.

Repository **tidak pernah** membuka transaksi dan **tidak pernah** memanggil
repository lain.

Untuk baris yang diperebutkan (stok, voucher, payout), pakai `FOR UPDATE` dengan
**urutan yang konsisten**.

---

## Usecase

Urutannya selalu: izin → validasi bisnis → transaksi.

```go
func (s *Service) DoSomething(ctx context.Context, in Input) (*domain.Thing, error) {
    if _, err := s.requireX(ctx, in.ID); err != nil { return nil, err }
    // validasi
    return nil, s.tx.WithinTx(ctx, func(ctx context.Context) error {
        // tulis + s.outbox.Save(ctx, evt)
        return nil
    })
}
```

**Tidak ada panggilan jaringan di dalam `WithinTx`.**

---

## Empat area HTTP

| Area | Middleware |
|---|---|
| `/api/v1/*` | `OptionalAuth` untuk katalog, `Authenticate` untuk keranjang & pesanan |
| `/api/v1/seller/*` | `Authenticate` + `RequireRole("seller")` |
| `/api/v1/admin/*` | `Authenticate` + `RequireRole("admin")` |
| `/webhook/*` | verifikasi signature, tanpa sesi, tanpa CORS |

Jangan memasang `Authenticate` untuk seluruh `/api/v1`. Halaman produk dan
pencarian harus bisa dibuka tanpa login.

---

## Worker

`cmd/worker`, proses terpisah. Delapan pekerjaan, semuanya gagal secara **senyap**
kalau tidak jalan.

Pekerjaan yang menyentuh uang butuh kunci terdistribusi di Redis **dan**
idempotency key di database.

---

## Perintah

```bash
make dev        # hot reload
make worker     # WAJIB jalan saat development
make test
make lint
make check
make migrate-create name=xxx
make migrate-up
```

Setelah mengubah kode, minimal `go build ./...`.

---

## Kesalahan yang sering terjadi di aplikasi ini

| Gejala | Penyebab |
|---|---|
| Nil pointer saat endpoint baru dipanggil | repository belum dirakit di `module.go` |
| Event handler tidak pernah jalan | `RegisterSubscriptions` lupa dipanggil di `registry.go` |
| Stok habis padahal barangnya ada | worker pelepas reservasi tidak jalan |
| Pesanan tidak pernah kedaluwarsa | worker tidak jalan |
| Deadlock saat checkout ramai | baris stok dikunci tanpa urutan tetap |
| Total order ≠ jumlah suborder | sisa pembulatan dibuang; cari pembagian tanpa `money.Distribute` |
| Rekonsiliasi tidak nol | ada jurnal tidak seimbang atau pergerakan uang tidak dicatat |
| `import cycle not allowed` | ada modul yang bergantung balik ke `order` |
