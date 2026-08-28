package domain

// TODO:
//   type Verification struct { ID; UserID; Kind; TokenHash; ExpiresAt; UsedAt }
//   func NewVerification(userID uuid.UUID, kind string, ttl time.Duration) (*Verification, raw string, err error)
//   func (v *Verification) Consume(now time.Time) error
// Token hanya boleh dipakai sekali. Hapus atau tandai used di dalam transaksi
// yang sama dengan aksinya.
