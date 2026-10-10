package fleet

import (
	"sort"

	"vyshkat/internal/domain"
)

const BeginnerEaseThreshold = 6

// VehicleInspector is the small abstraction Fleet Service needs from a service center.
// The concrete ServiceCenter is injected from outside.
type VehicleInspector interface {
	Inspect(domain.Transport) bool
}

// Service manages the active fleet and all inventory items on the balance.
type Service struct {
	inspector  VehicleInspector
	transports []domain.Transport
	inventory  []domain.InventoryItem
	knownIDs   map[string]struct{}
}

func NewService(inspector VehicleInspector) *Service {
	return &Service{
		inspector: inspector,
		knownIDs:  make(map[string]struct{}),
	}
}

// AddTransport runs technical inspection before adding transport to the active fleet.
func (s *Service) AddTransport(transport domain.Transport) error {
	if transport == nil {
		return ErrNilInventoryItem
	}
	if _, exists := s.knownIDs[transport.InventoryNumber()]; exists {
		return ErrDuplicateInventory
	}
	if !s.inspector.Inspect(transport) {
		return ErrInspectionFailed
	}

	s.transports = append(s.transports, transport)
	s.registerInventory(transport)
	return nil
}

// AddItem puts a non-transport inventory item on the balance.
func (s *Service) AddItem(item domain.InventoryItem) error {
	if item == nil {
		return ErrNilInventoryItem
	}
	if _, exists := s.knownIDs[item.InventoryNumber()]; exists {
		return ErrDuplicateInventory
	}

	s.registerInventory(item)
	return nil
}

func (s *Service) registerInventory(item domain.InventoryItem) {
	s.inventory = append(s.inventory, item)
	s.knownIDs[item.InventoryNumber()] = struct{}{}
}

func (s *Service) TransportCount() int {
	return len(s.transports)
}

// TotalDailyEnergyKWh sums energy consumption of every energy-consuming object on the balance.
func (s *Service) TotalDailyEnergyKWh() float64 {
	var total float64
	for _, item := range s.inventory {
		consumer, ok := item.(domain.EnergyConsumer)
		if !ok {
			continue
		}
		total += consumer.DailyEnergyKWh()
	}
	return total
}

// BeginnerFriendlyTransports returns accepted rental transport with ease >= 6.
func (s *Service) BeginnerFriendlyTransports() []domain.Transport {
	result := make([]domain.Transport, 0)
	for _, transport := range s.transports {
		if transport.BeginnerEase() >= BeginnerEaseThreshold {
			result = append(result, transport)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].InventoryNumber() < result[j].InventoryNumber()
	})
	return result
}

// Inventory returns all accepted transport and things sorted by inventory number.
func (s *Service) Inventory() []domain.InventoryItem {
	result := append([]domain.InventoryItem(nil), s.inventory...)
	sort.Slice(result, func(i, j int) bool {
		return result[i].InventoryNumber() < result[j].InventoryNumber()
	})
	return result
}
