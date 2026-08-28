package server

// TODO:
//   func New(cfg config.HTTPConfig, h http.Handler, log *slog.Logger) *Server
//   func (s *Server) Run(ctx context.Context) error
//     ListenAndServe di goroutine, tunggu ctx.Done(), Shutdown dengan batas waktu.
//
// Shutdown rapi penting: checkout yang sedang berjalan harus selesai. Pesanan
// yang terpotong di tengah transaksi memang akan rollback, tapi panggilan ke
// payment gateway yang sudah terlanjur terkirim tidak bisa ditarik.
