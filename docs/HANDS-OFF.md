# Serah Terima — Kondisi Repo Saat Ini

Dokumen orientasi untuk agent atau orang yang baru masuk. Tujuannya satu: dalam
lima menit kamu tahu **di mana kita sekarang**, **apa yang mengikat**, dan **ranjau
apa yang sedang aktif** — tanpa perlu membaca 2.000 baris dokumen lebih dulu.

- **Terakhir diperbarui:** 2026-08-31
- **Branch:** `main` (sudah menggabungkan `feature/backend/master` dan `feature/frontend/master`)

Ini bukan pengganti `AGENTS.md`. Aturan tetap di sana. Ini peta situasi.

---

## 1. Repo ini apa

Klontong Net — **marketplace** kebutuhan sehari-hari. Banyak penjual, tiap penjual
punya beberapa outlet dengan stok sendiri, pembeli bisa pilih kurir lokal,
ekspedisi, atau ambil di toko.

Satu keranjang bisa berisi barang dari beberapa penjual. Karena itu **satu Order
punya banyak Suborder**, dan **Suborder-lah unit kerjanya** — pengiriman,
pembatalan, komisi, dan pencairan semuanya di level itu. Kalau kamu cuma ingat satu
hal dari dokumen ini, ingat yang ini (ADR-002).

Sistem ini menyentuh uang orang dan stok barang nyata.

```
apps/backend/    Go + chi + PostgreSQL — seluruh sistem
apps/frontend/   React + TS + Vite — storefront (pembeli) + dashboard (penjual, admin)
contracts/       OpenAPI 3.0.3 — kontrak API yang MENGIKAT kedua app
docs/            arsitektur, ADR, fase integrasi
```

---

## 2. Di mana kita sekarang

| Bagian | Status | Keterangan |
|---|---|---|
| `contracts/` | **Selesai** | M1 (buyer) + M2 (seller/admin/webhook). `make contracts-check` hijau |
| `apps/backend` | **Scaffold + fase 01 sedang jalan** | 419 file `.go`, hampir semuanya `package` + TODO kontrak |
| `apps/frontend` | **Fase 00 selesai** | Workspace, 2 app Vite, `packages/{api,ui}`, `schema.d.ts` sudah tergenerate |
| Dokumen tugas | **Selesai & sudah diselaraskan** | Backend 15 fase, frontend 11 fase, plus fase integrasi lintas-app |
| CI | **Tidak ada** | Belum ada `.github/`. Gate dijalankan manual |

**Semua `migrations/*.up.sql` (15 file) masih komentar spec, belum ada DDL.**
Tiap fase menuliskan DDL-nya sendiri sesuai spec yang sudah tertulis di file itu.

### Yang sudah benar-benar ditulis di backend

Commit `1d90f6e` ("temp: try add function money") menambahkan **test lebih dulu**
untuk sebagian `internal/platform`: `money`, `auth`, `clock`, `config`, `errs`,
`eventbus`, `id`, `postgres`. Implementasinya **belum menyusul**.

Jadi `go build ./...` **hijau**, tapi `go test ./...` **tidak kompilasi** — banyak
fungsi yang dipanggil test belum ada (`RoundDown`, `RoundNearest`, `id.New`,
`id.Parse`, `id.SuborderNumber`, dan lain-lain).

**Ini kondisi yang diharapkan**, bukan kerusakan. Fase 01 memang dikerjakan dengan
test lebih dulu karena menyentuh uang. Merah sekarang, hijau saat fase 01 selesai.

---

## 3. Yang mengikat

Baca `AGENTS.md` untuk aturan lengkap. Yang paling sering dilanggar:

- **Uang** `money.Amount` int64 rupiah. Tidak ada `float`, tidak pernah. Persen =
  basis poin (250 = 2,5%). Pembagian pakai `money.Distribute`. Sisa pembulatan
  tidak dibuang. Komisi bulat ke bawah.
- **Harga dan total dihitung server.** Struct request checkout tidak punya field
  harga/total. Klien kirim pilihan, bukan angka.
- **Status hanya lewat `Transition()`.** Status Order induk dihitung dari
  anak-anaknya, tidak pernah disetel langsung. Satu penjual menolak **tidak**
  membatalkan pesanan penjual lain.
- **Stok** ditahan saat checkout, di-commit saat kirim, dilepas saat batal. Kunci
  baris `ORDER BY outlet_id, variant_id` sebelum `FOR UPDATE`.
- **Transaksi hanya di layer `app`.** Tidak ada panggilan jaringan di dalam
  transaksi database.
- **Event lewat `outbox.Save` di transaksi yang sama.** Tidak pernah `bus.Publish`
  dari usecase.
- **Kontrak API mengikat** — lihat bagian 4.
- **Buku besar double entry.** Tidak ada kolom saldo yang di-UPDATE.

### Repo ini scaffold, dan itu disengaja

Komentar TODO di file `.go` berisi kontrak yang sudah dipikirkan. Hormati.
**Kerjakan hanya yang diminta.** Melihat TODO tetangga yang menggoda? Sebutkan,
jangan dikerjakan.

---

## 4. `contracts/` mengikat kedua sisi

Backend dan frontend dikerjakan **paralel dan tidak saling menunggu**. Yang
menyambungkan keduanya cuma `contracts/`.

- Backend mengimplementasikan apa yang tertulis di kontrak: path, verb, nama field,
  enum, bentuk envelope. Tiap file fase backend 02–15 menunjuk `paths/*.yaml`-nya.
- Frontend men-generate `packages/api/src/schema.d.ts` dari kontrak dan membangun
  seluruh alurnya di atas mock.
- Perlu menyimpang? **Ubah kontrak di PR yang sama**, `make contracts-check` hijau,
  frontend `pnpm gen:api` ulang. Jangan menambal di satu sisi.

Detail lengkap: **ADR-015** di `DECISIONS.md`, dan `contracts/README.md`.

Yang ditetapkan kontrak dan sering tidak tertulis di sisi Go: field JSON
`snake_case`, waktu RFC3339 UTC, uang `integer`, envelope `{data,meta}` /
`{error:{code,message,fields?,retryable?}}`, paginasi keyset `?limit`+`?cursor`
dengan `meta.next_cursor`+`meta.has_more`, `Authorization: Bearer`, `X-Request-ID`,
`Retry-After`, `Idempotency-Key` wajib di checkout.

---

## 5. Ranjau yang sedang aktif

### 5.1 `platform/money/money.go` berisi implementasi yang salah

Commit `1d90f6e` menaruh isi sementara di `money.go` yang **tidak boleh dipercaya**.
Ini file paling menentukan di sistem; jangan membangun apa pun di atasnya sebelum
dibetulkan:

| Yang ada sekarang | Masalahnya |
|---|---|
| `func (a Amount) Sub(b Amount) Amount { return a / b }` | **pembagian**, bukan pengurangan |
| `func (a Amount) Mul(b Amount) Amount` | rupiah × rupiah tidak bermakna; yang dibutuhkan pengali `int64` |
| `ApplyBPS` membagi **1000** | harusnya **10000**. Hasilnya **10× kelebihan** |
| `func (a Amount) MarshalJSON() int` | tanda tangannya salah — `json.Marshaler` butuh `([]byte, error)`. Diam-diam tidak dipakai `encoding/json` |
| `String()` mengembalikan angka mentah | komentar menjanjikan `"Rp12.500"`; ADR-005 bilang format itu urusan frontend — luruskan komentarnya |
| `Distribute` mengembalikan `[]Amount{}` | stub kosong. Ini fungsi yang membagi voucher ke suborder |
| `rounding.go` masih TODO | `RoundDown` / `RoundNearest` belum ada, karena itu test tidak kompilasi |

`money_test.go` dan `rounding_test.go` sudah ditulis dan menggambarkan perilaku yang
benar. **Perbaiki implementasinya sampai test hijau; jangan mengubah test supaya
cocok dengan implementasi yang salah.**

### 5.2 `session_token` keranjang tamu belum dikunci di kontrak

`contracts/openapi/paths/cart.yaml:1` masih menulis "cookie/header"; frontend
mengasumsikan header. Selama masih mock tidak ada yang gagal — MSW menuruti asumsi
frontend. Baru meledak saat lepas mock.

**Kunci satu transport di kontrak sebelum frontend mulai fase 04.**

### 5.3 `apps/frontend/node_modules` belum terinstall

`pnpm gen:api`, `pnpm typecheck`, dan `pnpm check` belum bisa dibuktikan hijau.
Jalankan `pnpm -C apps/frontend install` dulu.

### 5.4 Worker wajib jalan saat development

`make dev` hanya menjalankan API. Tanpa `make worker` di terminal lain: reservasi
stok tidak pernah lepas, pesanan tidak pernah kedaluwarsa, dana tidak pernah cair.
Ketiganya bergejala sama — **sistem terlihat normal padahal tidak ada yang
terjadi**. Ini sumber kebingungan nomor satu di project ini.

---

## 6. Langkah berikutnya

### Kalau kamu mengerjakan backend

Berikutnya **fase 01 — Platform**. Baca `apps/backend/docs/tasks/01-platform.md`.

1. Betulkan `money.go` dan tulis `rounding.go` sampai `go test ./internal/platform/money/...` hijau. Mulai dari `Distribute` — property test-nya sudah ada.
2. Lanjutkan sisa `internal/platform` yang test-nya sudah merah: `id`, `config`, `errs`, `auth`, `clock`, `eventbus`, `postgres`.
3. `/healthz` balas 200.

Jangan lompat ke fase 10 sebelum 6 dan 9 selesai. Urutan penuh di
`apps/backend/docs/TASKS.md`.

### Kalau kamu mengerjakan frontend

Fase 00 sudah selesai. Berikutnya **fase 01 — Lapis API & kerangka app**. Baca
`apps/frontend/docs/tasks/01-lapis-api-shell.md`.

**Jangan menunggu backend.** Seluruh 11 fase bisa selesai di atas MSW. Kalau sampai
fase 05 backend belum siap, lanjut ke fase 06.

### Kalau backend fase 12 dan frontend fase 05 sudah sama-sama selesai

Jalankan `docs/tasks/INTEGRASI.md` — lepas mock, sambungkan ke API asli, cocokkan
respons dengan kontrak. Fase itu **tidak memblokir** siapa pun.

---

## 7. Dokumen mana dibaca kapan

| Dokumen | Baca kalau |
|---|---|
| `AGENTS.md` | **selalu, sebelum menulis kode apa pun** |
| `apps/backend/AGENTS.md` | tugasnya Go |
| `apps/backend/docs/TASKS.md` | mau tahu fase mana berikutnya |
| `apps/backend/docs/tasks/NN-*.md` | mengerjakan fase itu — checklist per file |
| `apps/backend/docs/GUIDES.md` | "tambah endpoint / usecase / modul / event" |
| `apps/frontend/docs/TASKS.md` + `tasks/NN-*.md` | tugasnya frontend |
| `contracts/README.md` | menyentuh bentuk API |
| `docs/ARCHITECTURE.md` | tugas menyentuh lebih dari satu modul |
| `docs/DECISIONS.md` | menyentuh uang, stok, pesanan, pembayaran |
| `docs/tasks/INTEGRASI.md` | menyambungkan kedua app |

---

## 8. Perintah

```bash
# backend
make setup && make migrate
make dev                 # API :8080
make worker              # terminal lain — WAJIB
make check               # gate backend

# frontend (tidak butuh backend jalan)
pnpm -C apps/frontend install
make fe-gen              # generate tipe TS dari kontrak
make mock                # Prism :4010 — terminal lain
make fe-dev
make fe-check            # gate frontend

# kontrak
make contracts-check     # lint + bundle
make check-all           # gate lintas-app
```

---

## 9. Berhenti dan tanya kalau

- Sesuatu bisa membuat **uang salah hitung, salah bayar, atau hilang**. Ini yang
  paling utama.
- Sesuatu bisa membuat **stok salah**, terjual padahal tidak ada, atau tertahan
  selamanya.
- Menembus batas modul atau menambah dependensi melingkar.
- **Menerima harga atau total dari klien.**
- Menambah panggilan jaringan di dalam transaksi database.
- Mengubah keputusan yang sudah ada ADR-nya.
- Menambah dependensi baru, atau menonaktifkan linter.

Jelaskan apa yang diminta, kenapa bertentangan, dan tawarkan alternatif. Jangan
diam-diam menyimpang, jangan juga menolak mentah-mentah. Kalau keputusannya berubah,
**tulis ADR baru di `DECISIONS.md`**.

---

## 10. Keputusan yang masih terbuka

Belum diputuskan, dan tiga yang pertama menyentuh uang — jangan diasumsikan sendiri:

- Komisi dihitung dari nilai barang saja atau termasuk ongkir?
- Siapa menanggung gagal bayar COD? Harus diputuskan sebelum COD dinyalakan.
- Jarak antar lokal: Haversine dengan faktor pengali, atau API routing?
- Batas waktu konfirmasi penjual, dan masa komplain setelah barang diterima.
- Transport `session_token` keranjang tamu (lihat 5.2).

Daftar lengkap di bagian "Yang masih terbuka" di `DECISIONS.md`.
