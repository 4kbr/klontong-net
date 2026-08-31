package clock

// SPEC TEST (TDD). Butuh dari clock.go:
//
//	type Clock interface{ Now() time.Time }
//	func New() Clock                 // pakai time.Now()
//	type Fixed struct{ T time.Time } // Now() selalu kembalikan T; untuk test
//	func (f Fixed) Now() time.Time
//
// Ini paket paling kecil di fase 01 — bagus untuk "pemanasan": bikin hijau
// dalam 5 menit lalu lanjut ke yang berat.
//
// Tidak ada dependency eksternal.

import (
	"testing"
	"time"
)

func TestFixed_SelaluSama(t *testing.T) {
	want := time.Date(2026, 8, 27, 9, 30, 0, 0, time.UTC)
	var c Clock = Fixed{T: want}

	if got := c.Now(); !got.Equal(want) {
		t.Fatalf("Fixed.Now() = %v, mau %v", got, want)
	}
	// dipanggil lagi -> tetap sama (tidak maju).
	if got := c.Now(); !got.Equal(want) {
		t.Fatalf("Fixed.Now() panggilan kedua = %v, mau %v", got, want)
	}
}

func TestNew_PakaiWaktuNyata(t *testing.T) {
	c := New()
	before := time.Now()
	got := c.Now()
	after := time.Now()

	if got.Before(before) || got.After(after) {
		t.Fatalf("New().Now() = %v, mau di antara %v dan %v", got, before, after)
	}
}
