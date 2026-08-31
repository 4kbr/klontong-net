# Fase 02 — Identity + Customer

> Prasyarat: fase 01. Guide: [`../GUIDES.md`](../GUIDES.md) §2 (tahap 2), §4–§6.
> ADR terkait: ADR-010 (snapshot alamat).
> Kontrak (**mengikat**, ADR-015): [`paths/auth.yaml`](../../../../contracts/openapi/paths/auth.yaml), [`me.yaml`](../../../../contracts/openapi/paths/me.yaml), [`addresses.yaml`](../../../../contracts/openapi/paths/addresses.yaml).

## Tujuan

Akun, sesi, peran, verifikasi; profil pembeli + alamat kirim. **Selesai kalau** bisa
daftar, login, dan menambah alamat.

## Aturan khusus fase ini

- Login **gagal selalu pesan identik** apa pun sebabnya (email tidak ada / password
  salah). Saat user tidak ada, tetap jalankan `Hasher.Compare` lawan dummy hash
  (timing-attack).
- Email dinormalisasi lowercase, phone ke E.164 **sebelum** cek keberadaan — kalau
  tidak, satu orang bisa punya banyak akun.
- Alamat **tidak pernah hard-delete** (order lama merujuknya) — `SoftDelete`. Order juga
  simpan salinan (ADR-010, dikerjakan di fase 10).
- Kunci keunikan email/phone = **unique index DB**; terjemahkan `23505` →
  `ErrEmailTaken`/`ErrPhoneTaken`.
- Verifikasi email **tidak memblokir** belanja.

## Urutan kerja

### identity

#### Domain
- [ ] `internal/modules/identity/internal/domain/user.go` — `User{Email,Phone,
      PasswordHash,Roles,Status}`; `NewUser` normalisasi email+phone; `AddRole`,
      `MarkEmailVerified/PhoneVerified`, `Suspend`.
- [ ] `.../domain/session.go` — `Session{RefreshTokenHash,ExpiresAt,RevokedAt,UA,IP}`
      (simpan **hash** refresh token); `IsActive`, `Revoke`.
- [ ] `.../domain/verification.go` — `Verification{Kind,TokenHash,ExpiresAt,UsedAt}`;
      `NewVerification` kembalikan token mentah sekali; `Consume` single-use di tx aksi.
- [ ] `.../domain/repository.go` — `UserRepository` (Create/Update/FindBy{ID,Email,
      Phone}/FindManyByID/ExistsBy{Email,Phone}/AddRole), `SessionRepository`
      (Create/FindByHash/Revoke/RevokeAllForUser/DeleteExpired), `VerificationRepository`.
- [ ] `.../domain/errors.go` — `ErrEmailTaken`, `ErrPhoneTaken`, `ErrInvalidCredential`,
      `ErrUserSuspended`, `ErrTokenExpired`, `ErrTokenUsed`, `ErrWeakPassword`.

#### Migrasi
- [ ] `migrations/000001_create_identity.up.sql` / `.down.sql` — `identity_users`,
      `identity_roles`, `identity_sessions` sesuai spec; unique index email & phone.

#### Infra
- [ ] `.../infra/user_repository.go`, `session_repository.go`,
      `verification_repository.go`, `mapper.go` — tiap method `ExecutorFrom` +
      `Translate`.

#### App (usecase)
- [ ] `.../app/register_usecase.go` — validasi kekuatan password → normalisasi → cek
      keberadaan (andalkan unique index) → hash + buat user peran `buyer` → `WithinTx`
      simpan user + verifikasi email + outbox `EventUserRegistered`.
- [ ] `.../app/login_usecase.go` — terima email ATAU phone satu field; dummy-compare saat
      user tidak ada; suspended → pesan berbeda; terbitkan access+refresh, simpan hash
      refresh + UA/IP. `Logout`, `LogoutAll`.
- [ ] `.../app/refresh_usecase.go` — rotasi: revoke lama, terbitkan baru; pemakaian token
      yang sudah di-revoke → indikasi pencurian → `RevokeAllForUser` + peringatan.
- [ ] `.../app/verification_usecase.go` — `SendEmailVerification`, `VerifyEmail`,
      `SendPhoneOTP`, `VerifyPhone`, `RequestPasswordReset` (selalu balas sukses),
      `ResetPassword`.
- [ ] `.../app/user_usecase.go` — `Me`, `UpdateProfile`, `ChangePassword` (lalu
      `RevokeAllForUser`), `GrantRole`.
- [ ] `.../app/service.go`, `dto.go`, `port_adapter.go`.

#### Transport
- [ ] `.../transport/rest/{routes,request,response,auth_handler,user_handler,handler}.go`
      — area `/api/v1` (login/register publik), `userID` dari context.

#### Port
- [ ] `internal/modules/identity/port.go` — `Get`, `GetMany` (**batch wajib**),
      `HasRole`, `GrantRole` (dipakai seller-signup: satu akun, tambah peran).
- [ ] `events.go` — `EventUserRegistered`, `EventEmailVerified`, `EventRoleGranted`.
- [ ] `module.go` — `Config{Pool,Tx,Outbox,Tokens,Hasher,Mailer,Clock,Logger}`,
      `New/RegisterRoutes/Port`.

### customer

#### Domain
- [ ] `.../customer/internal/domain/profile.go` — `Profile{UserID,DefaultAddressID,
      BirthDate,Gender}`.
- [ ] `.../domain/address.go` — `Address{...,Latitude/Longitude *float64,IsDefault,
      DeletedAt}`; `NewAddress` validasi recipient/phone/city/postal + phone E.164;
      `HasCoordinates`, `SoftDelete`. Cap ~20 per user.
- [ ] `.../domain/repository.go` — `ProfileRepository{Get,Upsert}`,
      `AddressRepository{Create,Update,FindByID,ListByUser,SoftDelete,SetDefault,
      CountByUser}`.
- [ ] `.../domain/errors.go` — not-found / validasi standar.

#### Migrasi
- [ ] `migrations/000002_create_customer.up.sql` / `.down.sql` — `customer_profiles`,
      `customer_addresses` (lat/lng nullable); pertimbangkan partial unique index
      `unique(user_id) where is_default`.

#### Infra + App + Transport
- [ ] `.../infra/{profile_repository,address_repository,mapper}.go`.
- [ ] `.../app/address_usecase.go` — `List/Get/Create/Update/Delete/SetDefaultAddress`;
      `SetDefault` satu tx (hapus lama, set baru); hapus default → alihkan atau tolak
      kalau satu-satunya.
- [ ] `.../app/profile_usecase.go` — `GetProfile`, `UpdateProfile`.
- [ ] `.../app/{service,dto,port_adapter}.go`.
- [ ] `.../transport/rest/{routes,request,response,profile_handler,address_handler,
      handler}.go` — area `/api/v1` (butuh login).

#### Port
- [ ] `internal/modules/customer/port.go` — `GetAddress(userID, addressID)` (**cek
      kepemilikan di dalam port**), `DefaultAddress(userID)`; `Address` bawa `*float64`
      lat/lng (nil ⇒ tidak ada antar lokal).
- [ ] `events.go` — `EventAddressAdded`, `EventAddressUpdated`.
- [ ] `module.go` — `Config{Pool,Tx,Outbox,Users identity.Port,Clock,Logger}`.

### Wiring
- [ ] `internal/app/registry.go` — build `identity` → lalu `customer` (terima
      `identityMod.Port()`).

## Test wajib

- Login gagal: email tidak ada vs password salah → response identik, durasi mirip.
- Register dua kali email sama → satu user, error `ErrEmailTaken` (bukan 500).
- Refresh token yang sudah di-revoke dipakai lagi → semua sesi user ter-revoke.
- `SetDefaultAddress` memindah default tepat satu baris.

## Sengaja TIDAK dikerjakan di fase ini

- Peran `seller` / registrasi toko (fase 3).
- Snapshot alamat di order (fase 10).
- Consumer `EventUserRegistered` di notification (fase 13).
