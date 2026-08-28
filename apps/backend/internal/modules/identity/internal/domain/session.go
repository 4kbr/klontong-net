package domain

// TODO:
//   type Session struct { ID; UserID; RefreshTokenHash; ExpiresAt; RevokedAt;
//                         UserAgent; IP; CreatedAt }
//   func (s *Session) IsActive(now) bool
//   func (s *Session) Revoke(now)
// Simpan HASH refresh token, bukan tokennya.
