package main

// Runner migrasi programatik. Untuk startup container dan integration test.
// Kerja sehari-hari tetap `make migrate-up`.
//
// TODO: flag -direction=up|down, -steps=N, -force=VERSION; baca DATABASE_URL
// lewat config.Load; exit code non-zero saat gagal.

func main() {
	// TODO
}
