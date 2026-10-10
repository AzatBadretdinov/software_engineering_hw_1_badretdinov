package fleet

import "errors"

var (
	ErrInspectionFailed   = errors.New("transport failed technical inspection")
	ErrDuplicateInventory = errors.New("inventory number already exists")
	ErrNilInventoryItem   = errors.New("inventory item must not be nil")
)
