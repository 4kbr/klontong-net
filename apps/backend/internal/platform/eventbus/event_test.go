package eventbus

// SPEC TEST (TDD). Butuh dari event.go:
//
//	type Event struct {
//	    ID uuid.UUID; Type string; AggregateType string
//	    AggregateID uuid.UUID; Payload json.RawMessage; OccurredAt time.Time
//	}
//	func New(eventType, aggType string, aggID uuid.UUID, payload any) (Event, error)
//	func Decode[T any](e Event) (T, error)
//
// PRASYARAT: `go get github.com/google/uuid`.
//
// Fokus: New() menandai event dengan benar dan mem-boks payload, lalu Decode[T]
// mengembalikannya utuh (payload = kontrak publik).

import (
	"testing"

	"github.com/google/uuid"
)

type suborderShipped struct {
	OrderID    string `json:"order_id"`
	SuborderID string `json:"suborder_id"`
	SellerID   string `json:"seller_id"`
	OutletID   string `json:"outlet_id"`
}

func TestNew_MengisiMetadata(t *testing.T) {
	aggID := uuid.New()
	in := suborderShipped{
		OrderID: "o1", SuborderID: "s1", SellerID: "seller1", OutletID: "outlet1",
	}

	e, err := New("order.suborder.shipped", "suborder", aggID, in)
	if err != nil {
		t.Fatalf("New error: %v", err)
	}
	if e.ID == uuid.Nil {
		t.Fatal("Event.ID belum di-generate")
	}
	if e.Type != "order.suborder.shipped" {
		t.Fatalf("Event.Type = %q", e.Type)
	}
	if e.AggregateType != "suborder" || e.AggregateID != aggID {
		t.Fatalf("agregat salah: type=%q id=%v", e.AggregateType, e.AggregateID)
	}
	if e.OccurredAt.IsZero() {
		t.Fatal("Event.OccurredAt kosong")
	}
	if len(e.Payload) == 0 {
		t.Fatal("Event.Payload kosong")
	}
}

func TestDecode_Roundtrip(t *testing.T) {
	in := suborderShipped{
		OrderID: "o1", SuborderID: "s1", SellerID: "seller1", OutletID: "outlet1",
	}
	e, err := New("order.suborder.shipped", "suborder", uuid.New(), in)
	if err != nil {
		t.Fatalf("New error: %v", err)
	}

	out, err := Decode[suborderShipped](e)
	if err != nil {
		t.Fatalf("Decode error: %v", err)
	}
	if out != in {
		t.Fatalf("Decode = %+v, mau %+v", out, in)
	}
}

func TestDecode_TipeTidakCocok_Error(t *testing.T) {
	e, err := New("x.y.z", "y", uuid.New(), map[string]any{"a": "bukan angka"})
	if err != nil {
		t.Fatalf("New error: %v", err)
	}

	type wrong struct {
		A int `json:"a"`
	}
	if _, err := Decode[wrong](e); err == nil {
		t.Fatal("Decode ke tipe yang tidak cocok tidak error")
	}
}
