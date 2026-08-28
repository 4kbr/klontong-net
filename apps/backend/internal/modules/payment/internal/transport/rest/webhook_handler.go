package rest

// TODO: HandleWebhook.
//   Baca RAW BODY sebelum di-decode — verifikasi signature butuh byte aslinya.
//   Balas 200 secepat mungkin setelah menyimpan; pemrosesan lanjutan lewat event.
//   Balas 200 juga untuk duplikat. Balas 4xx HANYA untuk signature tidak sah.
