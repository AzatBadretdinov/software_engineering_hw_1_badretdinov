package servicecenter

import (
	"testing"

	"vyshkat/internal/domain"
)

func TestServiceCenterInspect(t *testing.T) {
	center := New()

	working, err := domain.NewBicycle("Working", "TR-1", 7, true)
	if err != nil {
		t.Fatal(err)
	}
	broken, err := domain.NewBicycle("Broken", "TR-2", 7, false)
	if err != nil {
		t.Fatal(err)
	}

	if !center.Inspect(working) {
		t.Fatal("operational transport must be accepted")
	}
	if center.Inspect(broken) {
		t.Fatal("broken transport must be rejected")
	}
	if center.Inspect(nil) {
		t.Fatal("nil transport must be rejected")
	}
}
