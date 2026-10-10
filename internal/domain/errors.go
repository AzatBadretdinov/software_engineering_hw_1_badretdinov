package domain

import "errors"

var (
	ErrEmptyName            = errors.New("name must not be empty")
	ErrEmptyInventoryNumber = errors.New("inventory number must not be empty")
	ErrInvalidBeginnerEase  = errors.New("beginner ease must be between 1 and 10")
	ErrInvalidEnergy        = errors.New("daily energy consumption must not be negative")
)
