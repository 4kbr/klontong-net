package infra

// TODO: implementasi WebhookRepository.
// Record memakai INSERT ... ON CONFLICT (provider, event_id) DO NOTHING lalu
// memeriksa RowsAffected untuk mengetahui duplikat. Ini yang membuat pemrosesan
// ulang aman tanpa balapan.
