package domain

import (
	"errors"
	"testing"
)

func TestNewScooterRejectsInvalidEase(t *testing.T) {
	_, err := NewScooter("Scooter", "TR-1", 11, true, 1.0)
	if !errors.Is(err, ErrInvalidBeginnerEase) {
		t.Fatalf("expected ErrInvalidBeginnerEase, got %v", err)
	}
}

func TestNewChargingStationRejectsNegativeEnergy(t *testing.T) {
	_, err := NewChargingStation("Station", "IT-1", -1)
	if !errors.Is(err, ErrInvalidEnergy) {
		t.Fatalf("expected ErrInvalidEnergy, got %v", err)
	}
}

func TestInventoryConstructorsRejectEmptyData(t *testing.T) {
	if _, err := NewHelmet("", "IT-1"); !errors.Is(err, ErrEmptyName) {
		t.Fatalf("expected ErrEmptyName, got %v", err)
	}
	if _, err := NewHelmet("Helmet", "  "); !errors.Is(err, ErrEmptyInventoryNumber) {
		t.Fatalf("expected ErrEmptyInventoryNumber, got %v", err)
	}
}

func TestTransportInterfaces(t *testing.T) {
	scooter, err := NewScooter("Scooter", "TR-1", 8, true, 1.5)
	if err != nil {
		t.Fatal(err)
	}
	bike, err := NewBicycle("Bike", "TR-2", 6, true)
	if err != nil {
		t.Fatal(err)
	}

	var _ Transport = scooter
	var _ EnergyConsumer = scooter
	var _ Transport = bike

	if scooter.DailyEnergyKWh() != 1.5 || scooter.BeginnerEase() != 8 || !scooter.IsOperational() {
		t.Fatal("unexpected scooter state")
	}
}

func TestElectricBikeAndEnergyValidation(t *testing.T) {
	bike, err := NewElectricBike("E-bike", "TR-3", 7, true, 2.5)
	if err != nil {
		t.Fatal(err)
	}
	if bike.DailyEnergyKWh() != 2.5 || bike.InventoryNumber() != "TR-3" || bike.Name() != "E-bike" {
		t.Fatal("unexpected electric bike state")
	}

	if _, err := NewElectricBike("E-bike", "TR-4", 7, true, -0.1); !errors.Is(err, ErrInvalidEnergy) {
		t.Fatalf("expected ErrInvalidEnergy, got %v", err)
	}
	if _, err := NewScooter("Scooter", "TR-5", 7, true, -0.1); !errors.Is(err, ErrInvalidEnergy) {
		t.Fatalf("expected ErrInvalidEnergy, got %v", err)
	}
}

func TestBicycleConstructorPropagatesValidationError(t *testing.T) {
	if _, err := NewBicycle("Bike", "TR-1", 0, true); !errors.Is(err, ErrInvalidBeginnerEase) {
		t.Fatalf("expected ErrInvalidBeginnerEase, got %v", err)
	}
}
