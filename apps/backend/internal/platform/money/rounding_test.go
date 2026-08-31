package money

// SPEC TEST (TDD). Butuh dari rounding.go (lihat TODO di sana):
//
//	func RoundDown(v int64, unit int64) int64
//	func RoundNearest(v int64, unit int64) int64
//
// Tidak ada dependency. Jalankan: go test ./internal/platform/money/ -run TestRound

import "testing"

func TestRoundDown(t *testing.T) {
	tests := []struct {
		v, unit, want int64
	}{
		{v: 12_345, unit: 100, want: 12_300},
		{v: 12_300, unit: 100, want: 12_300}, // sudah kelipatan
		{v: 99, unit: 100, want: 0},
		{v: 0, unit: 100, want: 0},
		{v: 12_345, unit: 1, want: 12_345}, // unit 1 = tidak berubah
	}
	for _, tc := range tests {
		if got := RoundDown(tc.v, tc.unit); got != tc.want {
			t.Fatalf("RoundDown(%d, %d) = %d, mau %d", tc.v, tc.unit, got, tc.want)
		}
	}
}

func TestRoundNearest(t *testing.T) {
	tests := []struct {
		v, unit, want int64
	}{
		{v: 12_345, unit: 100, want: 12_300}, // .45 -> ke bawah
		{v: 12_355, unit: 100, want: 12_400}, // .55 -> ke atas
		{v: 12_350, unit: 100, want: 12_400}, // tepat setengah -> ke atas (tetapkan aturannya, konsisten)
		{v: 0, unit: 100, want: 0},
	}
	for _, tc := range tests {
		if got := RoundNearest(tc.v, tc.unit); got != tc.want {
			t.Fatalf("RoundNearest(%d, %d) = %d, mau %d "+
				"(kalau aturan setengah-ke-bawah yang kamu pilih, ubah baris ekspektasi ini)",
				tc.v, tc.unit, got, tc.want)
		}
	}
}
