package id

// SPEC TEST (TDD). Butuh dari id.go:
//
//	func New() uuid.UUID
//	func Parse(s string) (uuid.UUID, error)
//	func OrderNumber(t time.Time, seq int64) string     // "KN-20260827-000123"
//	func SuborderNumber(orderNumber string, idx int) string  // "KN-20260827-000123-1"
//
// PRASYARAT: `go get github.com/google/uuid` (dep resmi, lihat go.mod).
// Sampai itu ada, error "cannot find package" wajar.
//
// Test ini sengaja TIDAK meng-import "github.com/google/uuid" langsung — cukup
// pakai bentuk string supaya permukaan yang diuji tetap di API paket id.

import (
	"regexp"
	"testing"
	"time"
)

var (
	reOrderNumber    = regexp.MustCompile(`^KN-\d{8}-\d{6}$`)
	reSuborderNumber = regexp.MustCompile(`^KN-\d{8}-\d{6}-\d+$`)
)

func TestOrderNumber_Format(t *testing.T) {
	tm := time.Date(2026, 8, 27, 10, 0, 0, 0, time.UTC)

	got := OrderNumber(tm, 123)
	if !reOrderNumber.MatchString(got) {
		t.Fatalf("OrderNumber = %q, mau cocok %s", got, reOrderNumber)
	}
	if got != "KN-20260827-000123" {
		t.Fatalf("OrderNumber = %q, mau %q (tanggal dari t, seq di-pad 6 digit)",
			got, "KN-20260827-000123")
	}

	// seq besar tetap valid formatnya.
	if g := OrderNumber(tm, 999999); !reOrderNumber.MatchString(g) {
		t.Fatalf("OrderNumber(seq=999999) = %q, tidak cocok format", g)
	}
}

func TestSuborderNumber_Format(t *testing.T) {
	got := SuborderNumber("KN-20260827-000123", 1)
	if !reSuborderNumber.MatchString(got) {
		t.Fatalf("SuborderNumber = %q, mau cocok %s", got, reSuborderNumber)
	}
	if got != "KN-20260827-000123-1" {
		t.Fatalf("SuborderNumber = %q, mau %q", got, "KN-20260827-000123-1")
	}
}

func TestNew_UnikDanTerParse(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		s := New().String()
		if s == "" || s == "00000000-0000-0000-0000-000000000000" {
			t.Fatalf("New() menghasilkan uuid kosong/nil: %q", s)
		}
		if seen[s] {
			t.Fatalf("New() menghasilkan duplikat: %q", s)
		}
		seen[s] = true

		if _, err := Parse(s); err != nil {
			t.Fatalf("Parse(%q) error: %v", s, err)
		}
	}

	if _, err := Parse("bukan-uuid"); err == nil {
		t.Fatal("Parse(\"bukan-uuid\") tidak error")
	}
}
