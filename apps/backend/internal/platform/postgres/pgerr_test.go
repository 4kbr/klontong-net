package postgres

// SPEC TEST (TDD). Butuh dari pgerr.go:
//
//	func Translate(err error) error
//	  pgx.ErrNoRows           -> errs.ErrNotFound
//	  23505 unique            -> errs.ErrConflict   (+ nama constraint di Message/Code)
//	  23503 foreign key       -> errs.ErrConflict
//	  23514 check             -> errs.ErrInvalidInput
//	  40001 serialization     -> errs.ErrConflict, Retryable = true
//	  55P03 lock_not_available -> errs.ErrConflict, Retryable = true
//	  nil                     -> nil
//	  error lain              -> diteruskan apa adanya (atau errs.ErrInternal — tetapkan, lalu sesuaikan test)
//
// PRASYARAT:
//   - `go get github.com/jackc/pgx/v5`  (dep resmi driver postgres)
//   - paket internal/platform/errs sudah ada  (fase 01 langkah 2)
// Sampai keduanya siap, error "cannot find package" itu wajar.

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/yourorg/klontong-net/apps/backend/internal/platform/errs"
)

func pgErr(code, constraint string) *pgconn.PgError {
	return &pgconn.PgError{Code: code, ConstraintName: constraint}
}

func TestTranslate_Mapping(t *testing.T) {
	cases := []struct {
		name      string
		in        error
		wantKind  error // errs sentinel; nil artinya "harap nil"
		retryable bool
	}{
		{name: "nil -> nil", in: nil, wantKind: nil},
		{name: "no rows -> NotFound", in: pgx.ErrNoRows, wantKind: errs.ErrNotFound},
		{name: "23505 unique -> Conflict", in: pgErr("23505", "identity_users_email_key"), wantKind: errs.ErrConflict},
		{name: "23503 fk -> Conflict", in: pgErr("23503", "order_items_order_id_fkey"), wantKind: errs.ErrConflict},
		{name: "23514 check -> InvalidInput", in: pgErr("23514", "chk_qty_positive"), wantKind: errs.ErrInvalidInput},
		{name: "40001 serialization -> Conflict + retryable", in: pgErr("40001", ""), wantKind: errs.ErrConflict, retryable: true},
		{name: "55P03 lock -> Conflict + retryable", in: pgErr("55P03", ""), wantKind: errs.ErrConflict, retryable: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Translate(tc.in)

			if tc.wantKind == nil {
				if got != nil {
					t.Fatalf("Translate(nil) = %v, mau nil", got)
				}
				return
			}

			if !errors.Is(got, tc.wantKind) {
				t.Fatalf("Translate(%v) tidak errors.Is(%v)\n  dapat: %v", tc.in, tc.wantKind, got)
			}

			if tc.retryable {
				var e *errs.Error
				if !errors.As(got, &e) || !e.Retryable {
					t.Fatalf("Translate(%v) harus Retryable=true, dapat: %+v", tc.in, got)
				}
			}
		})
	}
}

func TestTranslate_23505_MembawaNamaConstraint(t *testing.T) {
	got := Translate(pgErr("23505", "promotion_redemptions_voucher_id_order_id_key"))

	var e *errs.Error
	if !errors.As(got, &e) {
		t.Fatalf("Translate 23505 bukan *errs.Error: %v", got)
	}
	// nama constraint sering menjelaskan persis apa yang bentrok -> harus ikut terbawa
	// (di Code atau Message — sesuaikan dengan implementasimu).
	if !containsAny(e.Code+" "+e.Message, "promotion_redemptions_voucher_id_order_id_key", "voucher") {
		t.Fatalf("nama constraint hilang dari error: code=%q msg=%q", e.Code, e.Message)
	}
}

func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		for i := 0; i+len(sub) <= len(s); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
	}
	return false
}
