package eventbus

// TODO: implementasi in-process.
//   type MemoryBus struct { mu sync.RWMutex; handlers map[string][]Handler; log }
//
// Publish menjalankan handler satu per satu. Error satu handler dicatat tapi
// tidak menghentikan yang lain. Sinkron dulu — gampang di-debug dan di-test.
//
// Subscribe hanya dipanggil saat start-up, di app/registry.go.
