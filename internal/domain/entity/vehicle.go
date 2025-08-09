package entity

type Vehicle struct {
	ID        uint   `json:"id" gorm:"type:uniqueidentifier;primary_key;default:NEWID()" example:"550e8400-e29b-41d4-a716-446655440000"`
	Make      string `json:"make" gorm:"not null;size:100" validate:"required,min=2,max=100" example:"Toyota"`
	Model     string `json:"model" gorm:"not null;size:100" validate:"required,min=2,max=100" example:"Camry"`
	Year      int    `json:"year" gorm:"not null" validate:"required,min=1886,max=2023" example:"2020"`
	Color     string `json:"color" gorm:"size:50" example:"Blue"`
	VIN       string `json:"vin" gorm:"unique;not null;size:17" validate:"required,len=17" example:"1HGCM82633A123456"`
	IsActive  bool   `json:"is_active" gorm:"default:true" example:"true"`
	CreatedAt string `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt string `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	DeletedAt string `json:"-" gorm:"index"`
}

func (Vehicle) TableName() string {
	return "vehicles"
}

type CreateVehicleRequest struct {
	Make  string `json:"make" validate:"required,min=2,max=100" example:"Toyota"`
	Model string `json:"model" validate:"required,min=2,max=100" example:"Camry"`
	Year  int    `json:"year" validate:"required,min=1886,max=2023" example:"2020"`
	Color string `json:"color" example:"Blue"`
	VIN   string `json:"vin" validate:"required,len=17" example:"1HGCM82633A123456"`
}
type UpdateVehicleRequest struct {
	Make     *string `json:"make,omitempty" validate:"omitempty,min=2,max=100" example:"Toyota"`
	Model    *string `json:"model,omitempty" validate:"omitempty,min=2,max=100" example:"Camry"`
	Year     *int    `json:"year,omitempty" validate:"omitempty,min=1886,max=2023" example:"2020"`
	Color    *string `json:"color,omitempty" example:"Blue"`
	VIN      *string `json:"vin,omitempty" validate:"omitempty,len=17" example:"1HGCM82633A123456"`
	IsActive *bool   `json:"is_active,omitempty" example:"true"`
}

type VehicleResponse struct {
	ID        uint   `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Make      string `json:"make" example:"Toyota"`
	Model     string `json:"model" example:"Camry"`
	Year      int    `json:"year" example:"2020"`
	Color     string `json:"color" example:"Blue"`
	VIN       string `json:"vin" example:"1HGCM82633A123456"`
	IsActive  bool   `json:"is_active" example:"true"`
	CreatedAt string `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt string `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	DeletedAt string `json:"-" example:""`
}
