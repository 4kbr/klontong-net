package httpx

// Format response seragam.
//
// TODO:
//   sukses: {"data": ..., "meta": {...}}
//   gagal : {"error": {"code": "...", "message": "...", "fields": {...}}}
//
//   func JSON / OK / Created / Accepted / NoContent / Paginated
//
// Nilai uang dikirim sebagai ANGKA (int64 rupiah), bukan string terformat.
// Pemformatan "Rp12.500" adalah urusan frontend dan bergantung locale.
