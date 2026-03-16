package shipment

type DriverInfo struct {
	DriverID string `validate:"required"`
	UnitID   string `validate:"required"`
}
