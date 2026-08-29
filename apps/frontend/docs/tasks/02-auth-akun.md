# Fase 02 — Auth & akun

> Prasyarat: 01 | Ref: kontrak `paths/auth.yaml`, `paths/me.yaml`,
> `paths/addresses.yaml`; ADR terkait: sesi/refresh di `TASKS.md` aturan global.

## Tujuan

Alur autentikasi penuh + manajemen akun & alamat untuk pembeli, plus pintu login
yang sama dipakai `dashboard` (seller/admin). **Selesai kalau** pengguna bisa
daftar, login dengan email atau telepon, sesi bertahan lewat refresh, logout
mencabut sesi, dan bisa mengelola profil + alamat; rute terproteksi mengarahkan
ke login lalu kembali ke tujuan.

## Aturan khusus fase ini

- Login memakai **satu field `identifier`** (email atau telepon). Semua kegagalan
  kredensial memunculkan pesan identik (`invalid_credential`); akun tersuspensi →
  `user_suspended` (pesan berbeda).
- `forgot-password` **selalu** memberi feedback sukses (202) — jangan bocorkan
  keberadaan akun.
- Verifikasi email **tidak memblokir belanja** — jangan paksa verifikasi sebelum
  masuk katalog/keranjang.
- Setelah login sukses di storefront: bila ada `session_token` keranjang tamu →
  `POST /api/v1/cart/merge` lalu hapus token (idempoten; aman bila keranjang kosong).
- `fields` pada error 400/422 dipetakan ke error field RHF.
- Access token TTL 15 menit — jangan tampilkan "sesi berakhir" saat refresh masih bisa; hanya logout bila refresh gagal.

## Urutan kerja

### Tipe & API (packages/api)
- [ ] `src/endpoints/auth.api.ts` — `register`, `login`, `refresh`, `logout`,
      `verifyEmail`, `requestOtp`, `verifyOtp`, `forgotPassword`, `resetPassword`.
- [ ] `src/endpoints/account.api.ts` — `getMe`, `changePassword` (`/me/password`), `getProfile`/`updateProfile` (`/profile`).
- [ ] `src/endpoints/address.api.ts` — `list`, `create`, `get`, `update`, `remove`, `setDefault`.
- [ ] `src/hooks/` — `useMe`, `useLogin`, `useRegister`, `useLogout`, `useAddresses`, `useAddressMutations` (+ invalidation `queryKeys.me`, `queryKeys.addresses`).
- [ ] `src/schemas/auth.ts`, `src/schemas/account.ts`, `src/schemas/address.ts` — Zod untuk tiap form (email/telp/E.164, password policy, koordinat alamat opsional).

### State (Zustand / Query)
- [ ] `authStore`: `signIn(session)` set token + `queryClient.setQueryData(queryKeys.me, user)`; `signOut()` panggil `logout` + `clearSession()` + `queryClient.clear()`.
- [ ] Aksi `mergeGuestCartAfterLogin()` (storefront) dipanggil dari `signIn`.

### Komponen (packages/ui)
- [ ] `AuthCard`, `OtpInput`, `PasswordField` (toggle + meter opsional), `FormField` (label + error + hint) bila belum ada.

### Halaman & rute
Storefront (`src/routes/(auth)/`, `AuthLayout`):
- [ ] `/login`, `/register`, `/verify-email`, `/otp` (request + verify), `/forgot-password`, `/reset-password?token=`.
- [ ] `src/routes/(account)/` di balik `RequireAuth`: `/account` (profil), `/account/password`, `/account/addresses` (+ tambah/edit dialog, set default, hapus).
Dashboard:
- [ ] `/login` memakai komponen auth yang sama; setelah login cek peran → `RequireRole` di layout; peran salah → halaman "akses ditolak".

### Wiring
- [ ] MSW handlers `auth` (stateful: user store, token rotasi, reuse-detection → 401 + cabut sesi), `me`, `addresses`.
- [ ] Fixtures user buyer/seller/admin untuk dev & test.
- [ ] Redirect pasca-login: honor `?next=`.

## Test wajib

- Register → sesi langsung aktif (201, `AuthSession`), `useMe` terisi tanpa fetch ulang.
- Login salah sandi vs email tak ada → **pesan & `code` identik** (`invalid_credential`).
- Login akun tersuspensi → `user_suspended` ditangani beda.
- Refresh token dirotasi: setelah refresh, token lama dipakai lagi → MSW balas 401 dan `authStore` jadi `anon`.
- Logout → `logout` dipanggil, storage bersih, rute terproteksi menolak.
- `forgot-password` untuk email tak terdaftar → tetap tampil sukses.
- Alamat: create → muncul di list; `setDefault` → hanya satu default; hapus default → backend/menu menetapkan default lain (ikuti respons).
- Login di storefront dengan `session_token` ada → `cart/merge` terpanggil sekali.
- Error 422 dengan `fields` → error muncul di field RHF yang tepat.

## Sengaja TIDAK dikerjakan

- Halaman katalog/keranjang — Fase 03/04.
- Onboarding seller (`seller/register`) — Fase 06.
- Preferensi notifikasi — Fase 05.
