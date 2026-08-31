package money

// =============================================================================
// SPEC TEST (TDD) — ditulis SEBELUM implementasi.
//
// File ini sengaja BELUM bisa di-compile sampai simbol berikut ada di money.go
// (semua sudah tertulis sebagai TODO di money.go + docs/tasks/01-platform.md §1):
//
//	type Amount int64
//	func (a Amount) Add(b Amount) Amount
//	func (a Amount) Sub(b Amount) Amount
//	func (a Amount) Mul(qty int64) Amount
//	func (a Amount) IsZero() bool
//	func (a Amount) IsNegative() bool
//	func (a Amount) MarshalJSON() ([]byte, error)   // ANGKA, bukan string
//	type BasisPoints int                            // 250 = 2,5%
//	func (a Amount) ApplyBPS(bps BasisPoints) Amount
//	func Distribute(total Amount, weights []Amount) []Amount
//
// CARA PAKAI:
//  1. `go test ./internal/platform/money/`  -> MERAH (gagal compile). Itu wajar.
//  2. Isi money.go sedikit demi sedikit.
//  3. Ulangi sampai HIJAU. Tidak ada yang perlu di-uncomment di sini.
//
// money TIDAK punya dependency eksternal — begitu money.go jadi, test ini jalan.
// =============================================================================

import (
	"encoding/json"
	"math/rand"
	"testing"
)

// ----------------------------------------------------------------------------
// Distribute — fungsi paling kritis. ADR-005 + docs/tasks/01-platform.md.
// ----------------------------------------------------------------------------

// Properti yang HARUS selalu benar untuk input apa pun (total >= 0, semua bobot >= 0,
// minimal satu bobot > 0):
//   - len(hasil) == len(bobot)
//   - jumlah seluruh hasil == total  (tidak ada rupiah yang hilang / muncul)
//   - tidak ada bagian negatif
//   - deterministik: input sama -> output sama persis
func TestDistribute_Properties(t *testing.T) {
	r := rand.New(rand.NewSource(1)) // seed tetap -> test reproducible

	for i := 0; i < 2000; i++ {
		total := Amount(r.Int63n(10_000_000)) // 0 .. 10 juta rupiah
		n := 1 + r.Intn(6)                    // 1 .. 6 bagian
		weights := make([]Amount, n)
		var sumW Amount
		for j := range weights {
			weights[j] = Amount(r.Int63n(1000))
			sumW += weights[j]
		}
		if sumW == 0 { // pastikan minimal satu bobot > 0
			weights[0] = 1
		}

		got := Distribute(total, weights)

		if len(got) != len(weights) {
			t.Fatalf("iter %d: len hasil = %d, mau %d", i, len(got), len(weights))
		}
		var sum Amount
		for _, part := range got {
			if part < 0 {
				t.Fatalf("iter %d: ada bagian negatif: %v (total=%d weights=%v)", i, got, total, weights)
			}
			sum += part
		}
		if sum != total {
			t.Fatalf("iter %d: jumlah bagian = %d, mau = %d (weights=%v hasil=%v)",
				i, sum, total, weights, got)
		}

		// deterministik
		again := Distribute(total, weights)
		for j := range got {
			if got[j] != again[j] {
				t.Fatalf("iter %d: Distribute tidak deterministik: %v vs %v", i, got, again)
			}
		}
	}
}

func TestDistribute_ContohKonkret(t *testing.T) {
	tests := []struct {
		name    string
		total   Amount
		weights []Amount
		want    []Amount // urutan hasil harus sesuai urutan weights
	}{
		{
			name:    "Rp10.000 dibagi rata tiga -> sisa 1 rupiah tidak dibuang",
			total:   10_000,
			weights: []Amount{1, 1, 1},
			// jumlah harus 10000. Alokasi sisa deterministik.
			// Contoh alokasi "sisa ke bagian terbesar/paling awal": {3334, 3333, 3333}.
			// Kalau implementasimu memilih aturan lain, sesuaikan angka di sini —
			// yang WAJIB: jumlahnya 10000 dan pola-nya konsisten tiap kali dipanggil.
			want: []Amount{3334, 3333, 3333},
		},
		{
			name:    "bobot proporsional",
			total:   1_000,
			weights: []Amount{1, 4}, // 20% : 80%
			want:    []Amount{200, 800},
		},
		{
			name:    "total nol",
			total:   0,
			weights: []Amount{3, 7},
			want:    []Amount{0, 0},
		},
		{
			name:    "satu bagian",
			total:   777,
			weights: []Amount{5},
			want:    []Amount{777},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Distribute(tc.total, tc.weights)

			var sum Amount
			for _, p := range got {
				sum += p
			}
			if sum != tc.total {
				t.Fatalf("jumlah bagian = %d, mau %d (hasil=%v)", sum, tc.total, got)
			}
			if len(got) != len(tc.want) {
				t.Fatalf("len hasil = %d, mau %d", len(got), len(tc.want))
			}
			for i := range tc.want {
				if got[i] != tc.want[i] {
					t.Fatalf("hasil = %v, mau %v", got, tc.want)
				}
			}
		})
	}
}

// ----------------------------------------------------------------------------
// Aritmetika Amount
// ----------------------------------------------------------------------------

func TestAmount_Aritmetika(t *testing.T) {
	if got := Amount(2_500).Add(Amount(1_000)); got != 3_500 {
		t.Fatalf("Add = %d, mau 3500", got)
	}
	if got := Amount(2_500).Sub(Amount(1_000)); got != 1_500 {
		t.Fatalf("Sub = %d, mau 1500", got)
	}
	if got := Amount(1_250).Mul(4); got != 5_000 {
		t.Fatalf("Mul = %d, mau 5000", got)
	}
	if !Amount(0).IsZero() {
		t.Fatal("Amount(0).IsZero() = false")
	}
	if !Amount(-1).IsNegative() {
		t.Fatal("Amount(-1).IsNegative() = false")
	}
	if Amount(1).IsNegative() {
		t.Fatal("Amount(1).IsNegative() = true")
	}
}

// ----------------------------------------------------------------------------
// MarshalJSON — WAJIB angka, bukan string terformat. ADR-005.
// ----------------------------------------------------------------------------

func TestAmount_MarshalJSON_KirimAngka(t *testing.T) {
	b, err := json.Marshal(Amount(12_500))
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	if string(b) != "12500" {
		t.Fatalf("JSON = %q, mau %q (bukan \"Rp12.500\", bukan \"12500\")", string(b), "12500")
	}

	// Di dalam struct pun tetap angka.
	type resp struct {
		Total Amount `json:"total"`
	}
	b, err = json.Marshal(resp{Total: 12_500})
	if err != nil {
		t.Fatalf("Marshal struct error: %v", err)
	}
	if string(b) != `{"total":12500}` {
		t.Fatalf("JSON = %s, mau {\"total\":12500}", string(b))
	}
}

// ----------------------------------------------------------------------------
// ApplyBPS — persen sebagai basis poin, pembulatan DITENTUKAN (ke bawah).
// ----------------------------------------------------------------------------

func TestAmount_ApplyBPS(t *testing.T) {
	tests := []struct {
		amount Amount
		bps    BasisPoints
		want   Amount
	}{
		{amount: 100_000, bps: 250, want: 2_500},   // 2,5%
		{amount: 100_000, bps: 1250, want: 12_500}, // 12,5%
		{amount: 999, bps: 250, want: 24},          // 24,975 -> ke bawah -> 24
		{amount: 0, bps: 500, want: 0},
		{amount: 100_000, bps: 0, want: 0},
	}
	for _, tc := range tests {
		if got := tc.amount.ApplyBPS(tc.bps); got != tc.want {
			t.Fatalf("Amount(%d).ApplyBPS(%d) = %d, mau %d (pembulatan harus KE BAWAH)",
				tc.amount, tc.bps, got, tc.want)
		}
	}
}
