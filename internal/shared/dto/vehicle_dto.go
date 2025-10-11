package dto

type CreateVehicleRequest struct {
	Brand        string `json:"brand" validate:"required"`
	Model        string `json:"model" validate:"required"`
	Year         int    `json:"year" validate:"required,min=1900,max=2100"`
	LicensePlate string `json:"license_plate" validate:"required"`
	VIN          string `json:"vin"`
	Mileage      int    `json:"mileage" validate:"min=0"`
}

type UpdateVehicleRequest struct {
	Brand        string `json:"brand"`
	Model        string `json:"model"`
	Year         int    `json:"year" validate:"omitempty,min=1900,max=2100"`
	LicensePlate string `json:"license_plate"`
	VIN          string `json:"vin"`
	Mileage      int    `json:"mileage" validate:"omitempty,min=0"`
}

type VehicleResponse struct {
	ID           string `json:"id"`
	OwnerID      string `json:"owner_id"`
	Brand        string `json:"brand"`
	Model        string `json:"model"`
	Year         int    `json:"year"`
	LicensePlate string `json:"license_plate"`
	VIN          string `json:"vin"`
	Mileage      int    `json:"mileage"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}
