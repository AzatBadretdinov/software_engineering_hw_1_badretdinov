package domain

// Helmet is an inventory item but does not consume energy.
type Helmet struct {
	inventoryBase
}

func NewHelmet(name, inventoryNumber string) (*Helmet, error) {
	base, err := newInventoryBase(name, inventoryNumber)
	if err != nil {
		return nil, err
	}
	return &Helmet{inventoryBase: base}, nil
}

// ChargingStation is an inventory item that consumes energy.
type ChargingStation struct {
	inventoryBase
	dailyEnergyKWh float64
}

func NewChargingStation(name, inventoryNumber string, dailyEnergyKWh float64) (*ChargingStation, error) {
	base, err := newInventoryBase(name, inventoryNumber)
	if err != nil {
		return nil, err
	}
	if dailyEnergyKWh < 0 {
		return nil, ErrInvalidEnergy
	}

	return &ChargingStation{inventoryBase: base, dailyEnergyKWh: dailyEnergyKWh}, nil
}

func (s *ChargingStation) DailyEnergyKWh() float64 {
	return s.dailyEnergyKWh
}
