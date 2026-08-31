package config

// SPEC TEST (TDD). Butuh dari config.go:
//
//	func Load() (Config, error)   // VALIDASI di akhir, gagalkan bila kombinasi tidak aman
//
// Test ini men-set environment variable lalu memanggil Load() dan memastikan
// Load() MENOLAK kombinasi yang berbahaya (docs/tasks/01-platform.md §2, ADR-003).
//
// >>> PENTING: nama environment variable di bawah adalah TEBAKAN. Samakan dengan
// nama yang benar-benar kamu baca di config.go. Yang diuji BUKAN nama env-nya,
// tapi PERILAKU Load(): kombinasi tidak valid -> error, kombinasi valid -> nil.
// Kalau nama env berbeda, ganti key di helper setEnv() di bawah — struktur test
// tidak perlu berubah.
//
// Tidak ada dependency eksternal (hanya os + testing + time).

import (
	"testing"
)

// baseEnv = satu set nilai yang VALID. Tiap sub-test menimpanya dengan satu
// nilai buruk lalu memastikan Load() error.
func baseEnv() map[string]string {
	return map[string]string{
		"APP_ENV": "test", // bukan "development" -> godotenv tidak dipakai

		"HTTP_ADDR":     ":8080",
		"PUBLIC_BASE_URL": "http://localhost:8080",

		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"DB_NAME":     "klontong",
		"DB_USER":     "klontong",
		"DB_PASSWORD": "secret",
		"DB_SSLMODE":  "disable",

		"REDIS_ADDR": "localhost:6379",

		"JWT_SECRET": "0123456789abcdef0123456789abcdef", // 32 byte -> valid

		"PAYMENT_PROVIDER":       "midtrans",
		"PAYMENT_SERVER_KEY":     "sk_test",
		"PAYMENT_CLIENT_KEY":     "ck_test",
		"PAYMENT_WEBHOOK_SECRET": "whsec_test",
		"PAYMENT_IS_PRODUCTION":  "false",
		"PAYMENT_EXPIRY":         "60m",

		"COMMISSION_DEFAULT_BPS":  "250",
		"SETTLEMENT_HOLD_PERIOD":  "72h",
		"PAYOUT_MINIMUM":          "50000",

		"STOCK_RESERVATION_TTL": "30m", // < PAYMENT_EXPIRY -> valid

		"STORAGE_BUCKET": "klontong-dev",
		"SMTP_ADDR":      "localhost:1025",
		"LOG_LEVEL":      "info",
	}
}

func applyEnv(t *testing.T, env map[string]string) {
	t.Helper()
	for k, v := range env {
		t.Setenv(k, v) // otomatis di-restore setelah test
	}
}

func TestLoad_KombinasiValid_TidakError(t *testing.T) {
	applyEnv(t, baseEnv())

	if _, err := Load(); err != nil {
		t.Fatalf("Load() dengan env valid malah error: %v", err)
	}
}

func TestLoad_MenolakKombinasiBerbahaya(t *testing.T) {
	cases := []struct {
		name     string
		override map[string]string
		// potongan kata yang diharapkan muncul di pesan error (opsional, boleh "")
		wantErrContains string
	}{
		{
			name:            "STOCK_RESERVATION_TTL >= PAYMENT_EXPIRY (ADR-003)",
			override:        map[string]string{"STOCK_RESERVATION_TTL": "90m", "PAYMENT_EXPIRY": "60m"},
			wantErrContains: "",
		},
		{
			name:            "STOCK_RESERVATION_TTL == PAYMENT_EXPIRY juga ditolak",
			override:        map[string]string{"STOCK_RESERVATION_TTL": "60m", "PAYMENT_EXPIRY": "60m"},
			wantErrContains: "",
		},
		{
			name:            "JWT_SECRET kurang dari 32 byte",
			override:        map[string]string{"JWT_SECRET": "pendek"},
			wantErrContains: "",
		},
		{
			name:            "PAYMENT_WEBHOOK_SECRET kosong padahal PAYMENT_PROVIDER terisi",
			override:        map[string]string{"PAYMENT_WEBHOOK_SECRET": ""},
			wantErrContains: "",
		},
		{
			name:            "COMMISSION_DEFAULT_BPS di atas 10000",
			override:        map[string]string{"COMMISSION_DEFAULT_BPS": "10001"},
			wantErrContains: "",
		},
		{
			name:            "COMMISSION_DEFAULT_BPS negatif",
			override:        map[string]string{"COMMISSION_DEFAULT_BPS": "-1"},
			wantErrContains: "",
		},
		{
			name: "produksi: DB_SSLMODE=disable tidak boleh",
			override: map[string]string{
				"PAYMENT_IS_PRODUCTION": "true",
				"APP_ENV":               "production",
				"DB_SSLMODE":            "disable",
				"PUBLIC_BASE_URL":       "https://api.klontong.example",
			},
			wantErrContains: "",
		},
		{
			name: "produksi: PUBLIC_BASE_URL harus https",
			override: map[string]string{
				"PAYMENT_IS_PRODUCTION": "true",
				"APP_ENV":               "production",
				"DB_SSLMODE":            "require",
				"PUBLIC_BASE_URL":       "http://api.klontong.example",
			},
			wantErrContains: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			env := baseEnv()
			for k, v := range tc.override {
				env[k] = v
			}
			applyEnv(t, env)

			_, err := Load()
			if err == nil {
				t.Fatalf("Load() TIDAK error padahal kombinasi tidak valid: %s", tc.name)
			}
			if tc.wantErrContains != "" && !contains(err.Error(), tc.wantErrContains) {
				t.Fatalf("pesan error = %q, mau memuat %q", err.Error(), tc.wantErrContains)
			}
		})
	}
}

func contains(s, sub string) bool {
	return len(sub) == 0 || (len(s) >= len(sub) && indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
