package domain

// TODO:
//   type User struct { ID; Email; Phone; PasswordHash; FullName; AvatarURL;
//                      EmailVerifiedAt; PhoneVerifiedAt; Status; Roles []string;
//                      CreatedAt; UpdatedAt }
//   func NewUser(email, phone, fullName, passwordHash string) (*User, error)
//       normalisasi email ke lowercase, normalisasi nomor HP ke format E.164
//   func (u *User) AddRole(role string) error
//   func (u *User) MarkEmailVerified(now) / MarkPhoneVerified(now)
//   func (u *User) Suspend(reason string) error
//
// Normalisasi nomor HP penting: "08123456789", "+628123456789", dan
// "628123456789" adalah orang yang sama. Tanpa normalisasi, satu orang bisa
// punya tiga akun dan tidak bisa login dengan nomor yang ia ingat.
