package main

// Entry point HTTP API Klontong Net.
//
// TODO — urutan di func main:
//  1. config.Load()
//  2. logger.New(cfg.Log)
//  3. context root yang dibatalkan SIGINT/SIGTERM (signal.NotifyContext)
//  4. postgres.NewPool + Ping. Gagal ping = mati sekarang.
//  5. redis.New (cache katalog, rate limit, kunci idempotensi)
//  6. dependency bersama: txManager, eventBus, outboxStore, tokenIssuer,
//     hasher, clock, storage, paymentGateway, shippingProvider
//  7. app.NewRegistry(deps) — membangun semua modul dan menyambungkannya
//  8. httpx.NewRouter, lalu registry.MountRoutes(r)
//     Router punya EMPAT area dengan middleware berbeda:
//       /api/v1/*         pembeli (sebagian butuh login, sebagian publik)
//       /api/v1/seller/*  dasbor penjual, wajib login + peran penjual
//       /api/v1/admin/*   panel marketplace, wajib peran admin
//       /webhook/*        payment gateway & kurir; verifikasi signature, TANPA sesi
//  9. server.New(...).Run(ctx)
//
// Worker TIDAK dijalankan di sini. Lihat cmd/worker dan ADR-013.

func main() {
	// TODO
}
