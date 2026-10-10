package fleet

import (
	"errors"
	"testing"

	"vyshkat/internal/domain"
)

type fakeInspector struct {
	decision bool
	calls    int
}

func (f *fakeInspector) Inspect(domain.Transport) bool {
	f.calls++
	return f.decision
}

func mustScooter(t *testing.T, id string, ease int, energy float64) *domain.Scooter {
	t.Helper()
	value, err := domain.NewScooter("Scooter "+id, id, ease, true, energy)
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func TestAddTransportAcceptedByInjectedInspector(t *testing.T) {
	// Arrange
	inspector := &fakeInspector{decision: true}
	service := NewService(inspector)
	transport := mustScooter(t, "TR-1", 8, 1.5)

	// Act
	err := service.AddTransport(transport)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if service.TransportCount() != 1 {
		t.Fatalf("expected 1 transport, got %d", service.TransportCount())
	}
	if inspector.calls != 1 {
		t.Fatalf("expected inspector to be called once, got %d", inspector.calls)
	}
}

func TestAddTransportRejectedByInjectedInspector(t *testing.T) {
	// Arrange: fake dependency always rejects.
	inspector := &fakeInspector{decision: false}
	service := NewService(inspector)
	transport := mustScooter(t, "TR-1", 8, 1.5)

	// Act
	err := service.AddTransport(transport)

	// Assert
	if !errors.Is(err, ErrInspectionFailed) {
		t.Fatalf("expected ErrInspectionFailed, got %v", err)
	}
	if service.TransportCount() != 0 {
		t.Fatalf("rejected transport entered fleet")
	}
	if len(service.Inventory()) != 0 {
		t.Fatalf("rejected transport entered inventory")
	}
}

func TestBeginnerFriendlyTransportsUsesThresholdSix(t *testing.T) {
	service := NewService(&fakeInspector{decision: true})

	for _, transport := range []domain.Transport{
		mustScooter(t, "TR-1", 5, 1.0),
		mustScooter(t, "TR-2", 6, 1.0),
		mustScooter(t, "TR-3", 9, 1.0),
	} {
		if err := service.AddTransport(transport); err != nil {
			t.Fatal(err)
		}
	}

	got := service.BeginnerFriendlyTransports()
	if len(got) != 2 {
		t.Fatalf("expected 2 beginner-friendly transports, got %d", len(got))
	}
	if got[0].InventoryNumber() != "TR-2" || got[1].InventoryNumber() != "TR-3" {
		t.Fatalf("unexpected beginner-friendly list")
	}
}

func TestTotalDailyEnergyIncludesEveryEnergyConsumerOnBalance(t *testing.T) {
	service := NewService(&fakeInspector{decision: true})

	scooter := mustScooter(t, "TR-1", 7, 1.8)
	bike, err := domain.NewBicycle("Bike", "TR-2", 7, true)
	if err != nil {
		t.Fatal(err)
	}
	station, err := domain.NewChargingStation("Station", "IT-1", 4.2)
	if err != nil {
		t.Fatal(err)
	}
	helmet, err := domain.NewHelmet("Helmet", "IT-2")
	if err != nil {
		t.Fatal(err)
	}

	for _, transport := range []domain.Transport{scooter, bike} {
		if err := service.AddTransport(transport); err != nil {
			t.Fatal(err)
		}
	}
	for _, item := range []domain.InventoryItem{station, helmet} {
		if err := service.AddItem(item); err != nil {
			t.Fatal(err)
		}
	}

	if got := service.TotalDailyEnergyKWh(); got != 6.0 {
		t.Fatalf("expected 6.0 kWh, got %.1f", got)
	}
}

func TestInventoryContainsTransportAndThingsAndIsSorted(t *testing.T) {
	service := NewService(&fakeInspector{decision: true})

	helmet, err := domain.NewHelmet("Helmet", "IT-2")
	if err != nil {
		t.Fatal(err)
	}
	station, err := domain.NewChargingStation("Station", "IT-1", 1)
	if err != nil {
		t.Fatal(err)
	}
	transport := mustScooter(t, "TR-1", 8, 1)

	if err := service.AddItem(helmet); err != nil {
		t.Fatal(err)
	}
	if err := service.AddTransport(transport); err != nil {
		t.Fatal(err)
	}
	if err := service.AddItem(station); err != nil {
		t.Fatal(err)
	}

	got := service.Inventory()
	want := []string{"IT-1", "IT-2", "TR-1"}
	if len(got) != len(want) {
		t.Fatalf("expected %d items, got %d", len(want), len(got))
	}
	for i, item := range got {
		if item.InventoryNumber() != want[i] {
			t.Fatalf("at %d: expected %s, got %s", i, want[i], item.InventoryNumber())
		}
	}
}

func TestDuplicateInventoryNumberRejected(t *testing.T) {
	service := NewService(&fakeInspector{decision: true})
	transport := mustScooter(t, "SAME", 8, 1)
	helmet, err := domain.NewHelmet("Helmet", "SAME")
	if err != nil {
		t.Fatal(err)
	}

	if err := service.AddTransport(transport); err != nil {
		t.Fatal(err)
	}
	if err := service.AddItem(helmet); !errors.Is(err, ErrDuplicateInventory) {
		t.Fatalf("expected ErrDuplicateInventory, got %v", err)
	}
}

func TestNilInventoryRejected(t *testing.T) {
	service := NewService(&fakeInspector{decision: true})
	if err := service.AddTransport(nil); !errors.Is(err, ErrNilInventoryItem) {
		t.Fatalf("expected ErrNilInventoryItem, got %v", err)
	}
	if err := service.AddItem(nil); !errors.Is(err, ErrNilInventoryItem) {
		t.Fatalf("expected ErrNilInventoryItem, got %v", err)
	}
}
