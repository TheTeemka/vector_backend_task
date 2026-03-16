package domain

import "errors"

var (
	ErrShipmentNotFound        = errors.New("shipment not found")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrDuplicateReference      = errors.New("shipment with this reference number already exists")
	ErrInvalidInput            = errors.New("invalid input")
)
