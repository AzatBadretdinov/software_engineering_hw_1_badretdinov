package domain

// Scooter is an electric scooter and therefore consumes energy.
type Scooter struct {
	transportBase
	dailyEnergyKWh float64
}

func NewScooter(name, inventoryNumber string, beginnerEase int, operational bool, dailyEnergyKWh float64) (*Scooter, error) {
	base, err := newTransportBase(name, inventoryNumber, beginnerEase, operational)
	if err != nil {
		return nil, err
	}
	if dailyEnergyKWh < 0 {
		return nil, ErrInvalidEnergy
	}

	return &Scooter{transportBase: base, dailyEnergyKWh: dailyEnergyKWh}, nil
}

func (s *Scooter) DailyEnergyKWh() float64 {
	return s.dailyEnergyKWh
}

// Bicycle is a regular bicycle. It intentionally does not implement EnergyConsumer.
type Bicycle struct {
	transportBase
}

func NewBicycle(name, inventoryNumber string, beginnerEase int, operational bool) (*Bicycle, error) {
	base, err := newTransportBase(name, inventoryNumber, beginnerEase, operational)
	if err != nil {
		return nil, err
	}
	return &Bicycle{transportBase: base}, nil
}

// ElectricBike is an electric bicycle and therefore consumes energy.
type ElectricBike struct {
	transportBase
	dailyEnergyKWh float64
}

func NewElectricBike(name, inventoryNumber string, beginnerEase int, operational bool, dailyEnergyKWh float64) (*ElectricBike, error) {
	base, err := newTransportBase(name, inventoryNumber, beginnerEase, operational)
	if err != nil {
		return nil, err
	}
	if dailyEnergyKWh < 0 {
		return nil, ErrInvalidEnergy
	}

	return &ElectricBike{transportBase: base, dailyEnergyKWh: dailyEnergyKWh}, nil
}

func (b *ElectricBike) DailyEnergyKWh() float64 {
	return b.dailyEnergyKWh
}
