package auth

// SPEC TEST (TDD). Butuh dari jwt.go:
//
//	type Claims struct { UserID uuid.UUID; Roles []string; jwt.RegisteredClaims }
//	type JWT struct{ ... }
//	func NewJWT(cfg config.JWTConfig) *JWT
//	func (j *JWT) IssueAccess(userID uuid.UUID, roles []string) (token string, expiresAt time.Time, err error)
//	func (j *JWT) Verify(token string) (Claims, error)   // WAJIB tolak "alg: none"
//
// PRASYARAT:
//   - `go get github.com/golang-jwt/jwt/v5 github.com/google/uuid`
//   - config.JWTConfig sudah ada di paket config
// Sampai itu ada, error "cannot find package" wajar.
//
// >>> config.JWTConfig di bawah adalah TEBAKAN bentuknya. Sesuaikan field-nya
// dengan config.go milikmu. Yang WAJIB diuji: Verify menolak token beralgoritma
// "none", dan token hasil IssueAccess bisa di-Verify balik dengan klaim utuh.

import (
	"encoding/base64"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/yourorg/klontong-net/apps/backend/internal/platform/config"
)

func newTestJWT(t *testing.T) *JWT {
	t.Helper()
	// TEBAKAN bentuk config.JWTConfig — ganti sesuai punyamu.
	return NewJWT(config.JWTConfig{
		Secret:    "0123456789abcdef0123456789abcdef", // >= 32 byte
		AccessTTL: 15 * time.Minute,
		Issuer:    "klontong-net",
	})
}

func TestVerify_TolakAlgNone(t *testing.T) {
	j := newTestJWT(t)

	// Rakit token "alg: none" secara manual (tanpa signature).
	// header  = {"alg":"none","typ":"JWT"}
	// payload = {"sub":"...","roles":["buyer"]}
	enc := func(s string) string {
		return base64.RawURLEncoding.EncodeToString([]byte(s))
	}
	header := enc(`{"alg":"none","typ":"JWT"}`)
	payload := enc(`{"sub":"3f1a6c2e-9b7d-4e2a-8c11-0a2b3c4d5e6f","roles":["buyer"],"iss":"klontong-net"}`)
	algNoneToken := header + "." + payload + "." // signature kosong

	if _, err := j.Verify(algNoneToken); err == nil {
		t.Fatal("Verify menerima token alg:none — WAJIB ditolak (celah keamanan klasik)")
	}
}

func TestVerify_TolakTokenAsal(t *testing.T) {
	j := newTestJWT(t)

	for _, tok := range []string{
		"",
		"bukan.jwt",
		"a.b.c",
		strings.Repeat("x", 200),
	} {
		if _, err := j.Verify(tok); err == nil {
			t.Fatalf("Verify(%q) tidak error", tok)
		}
	}
}

func TestIssueThenVerify_Roundtrip(t *testing.T) {
	j := newTestJWT(t)
	userID := uuid.New()
	roles := []string{"buyer", "seller"}

	tok, exp, err := j.IssueAccess(userID, roles)
	if err != nil {
		t.Fatalf("IssueAccess error: %v", err)
	}
	if tok == "" {
		t.Fatal("IssueAccess mengembalikan token kosong")
	}
	if time.Until(exp) <= 0 || time.Until(exp) > 16*time.Minute {
		t.Fatalf("expiresAt = %v (dalam %v), mau ~15 menit ke depan", exp, time.Until(exp))
	}

	claims, err := j.Verify(tok)
	if err != nil {
		t.Fatalf("Verify token sendiri malah error: %v", err)
	}
	if claims.UserID != userID {
		t.Fatalf("claims.UserID = %v, mau %v", claims.UserID, userID)
	}
	if strings.Join(claims.Roles, ",") != strings.Join(roles, ",") {
		t.Fatalf("claims.Roles = %v, mau %v", claims.Roles, roles)
	}
}

func TestVerify_TolakSecretSalah(t *testing.T) {
	issuer := newTestJWT(t)
	tok, _, err := issuer.IssueAccess(uuid.New(), []string{"buyer"})
	if err != nil {
		t.Fatalf("IssueAccess error: %v", err)
	}

	other := NewJWT(config.JWTConfig{
		Secret:    "ffffffffffffffffffffffffffffffff", // secret berbeda
		AccessTTL: 15 * time.Minute,
		Issuer:    "klontong-net",
	})
	if _, err := other.Verify(tok); err == nil {
		t.Fatal("Verify dengan secret berbeda tidak error")
	}
}
