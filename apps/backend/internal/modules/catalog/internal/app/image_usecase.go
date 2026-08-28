package app

// TODO: RequestUploadURL (presigned PUT), AttachImage, DeleteImage, ReorderImages.
// Unggah langsung dari browser ke object storage; API hanya menerbitkan URL dan
// mencatat hasilnya. Menyalurkan foto lewat API berarti membayar bandwidth dua
// kali dan menahan goroutine lama.
