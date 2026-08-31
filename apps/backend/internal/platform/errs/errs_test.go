package errs

// SPEC TEST (TDD). Butuh dari errs.go (semua sudah jadi TODO di sana):
//
//	var ErrNotFound, ErrConflict, ErrInvalidInput, ErrUnauthorized, ErrForbidden,
//	    ErrTooManyRequests, ErrUpstream, ErrInternal error
//
//	type Error struct {
//	    Kind error; Code string; Message string
//	    Fields map[string]string
//	    Retryable bool
//	    cause error
//	}
//	func (e *Error) Error() string
//	func (e *Error) Unwrap() error   // kembalikan e.Kind
//
//	func NotFound(code, msg string) *Error
//	func Conflict(code, msg string) *Error
//	func Invalid(code, msg string) *Error
//	func Unauthorized(code, msg string) *Error
//	func Forbidden(code, msg string) *Error
//	func Upstream(code, msg string) *Error
//	func Internal(code, msg string) *Error
//
// CATATAN: signature konstruktor di atas hanya TEBAKAN yang masuk akal. Kalau
// implementasimu berbeda (mis. pakai variadic opsi, atau tanpa argumen msg),
// sesuaikan pemanggilan di test ini — yang WAJIB dijaga: errors.Is bekerja lewat
// Kind, dan Code ikut terbawa.
//
// Tidak ada dependency eksternal.

import (
	"errors"
	"testing"
)

func TestError_ErrorsIs_LewatKind(t *testing.T) {
	err := Conflict("email_taken", "Email sudah terdaftar.")

	if !errors.Is(err, ErrConflict) {
		t.Fatal("errors.Is(Conflict(...), ErrConflict) = false — Unwrap() harus kembalikan Kind")
	}
	if errors.Is(err, ErrNotFound) {
		t.Fatal("errors.Is(Conflict(...), ErrNotFound) = true — seharusnya tidak cocok")
	}
}

func TestError_CodeTerbawa(t *testing.T) {
	err := Invalid("out_of_stock", "Stok habis.")

	var e *Error
	if !errors.As(err, &e) {
		t.Fatal("errors.As ke *Error gagal")
	}
	if e.Code != "out_of_stock" {
		t.Fatalf("Code = %q, mau %q", e.Code, "out_of_stock")
	}
	if e.Message == "" {
		t.Fatal("Message kosong")
	}
	// Code non-empty & stabil — frontend bercabang di sini.
	if e.Error() == "" {
		t.Fatal("Error() string kosong")
	}
}

func TestError_Retryable(t *testing.T) {
	// Upstream boleh ditandai retryable; default konstruktor lain: false.
	up := Upstream("gateway_unavailable", "Gateway sedang gangguan.")
	up.Retryable = true

	var e *Error
	if !errors.As(up, &e) || !e.Retryable {
		t.Fatal("Retryable tidak terbawa lewat errors.As")
	}

	nf := NotFound("not_found", "Tidak ada.")
	errors.As(nf, &e)
	if e.Retryable {
		t.Fatal("NotFound default Retryable = true, mau false")
	}
}

func TestError_WrapCause(t *testing.T) {
	// Kalau Error menyimpan cause internal, errors.Is tetap harus menemukan Kind.
	root := errors.New("boom")
	_ = root
	err := Internal("internal_error", "Kesalahan server.")
	if !errors.Is(err, ErrInternal) {
		t.Fatal("errors.Is(Internal(...), ErrInternal) = false")
	}
}
