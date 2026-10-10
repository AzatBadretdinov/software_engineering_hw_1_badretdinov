package domain

// InventoryItem is anything that is kept on the service balance and has an inventory number.
type InventoryItem interface {
	Name() string
	InventoryNumber() string
}

// EnergyConsumer is implemented only by objects that consume energy.
type EnergyConsumer interface {
	DailyEnergyKWh() float64
}

// Transport is rental transport managed by the fleet.
// Not every Transport is an EnergyConsumer: a regular bicycle is a valid transport too.
type Transport interface {
	InventoryItem
	BeginnerEase() int
	IsOperational() bool
}
