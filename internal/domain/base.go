package domain

import "strings"

type inventoryBase struct {
	name            string
	inventoryNumber string
}

func newInventoryBase(name, inventoryNumber string) (inventoryBase, error) {
	name = strings.TrimSpace(name)
	inventoryNumber = strings.TrimSpace(inventoryNumber)

	if name == "" {
		return inventoryBase{}, ErrEmptyName
	}
	if inventoryNumber == "" {
		return inventoryBase{}, ErrEmptyInventoryNumber
	}

	return inventoryBase{name: name, inventoryNumber: inventoryNumber}, nil
}

func (b inventoryBase) Name() string {
	return b.name
}

func (b inventoryBase) InventoryNumber() string {
	return b.inventoryNumber
}

type transportBase struct {
	inventoryBase
	beginnerEase int
	operational  bool
}

func newTransportBase(name, inventoryNumber string, beginnerEase int, operational bool) (transportBase, error) {
	base, err := newInventoryBase(name, inventoryNumber)
	if err != nil {
		return transportBase{}, err
	}
	if beginnerEase < 1 || beginnerEase > 10 {
		return transportBase{}, ErrInvalidBeginnerEase
	}

	return transportBase{
		inventoryBase: base,
		beginnerEase:  beginnerEase,
		operational:   operational,
	}, nil
}

func (b transportBase) BeginnerEase() int {
	return b.beginnerEase
}

func (b transportBase) IsOperational() bool {
	return b.operational
}
