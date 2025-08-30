package servicepackage

import (
	"gorm.io/gorm"
)

type Repository interface {
	CreatePackage(pkg *ServicePackage) error
	GetPackages() ([]ServicePackage, error)
	GetPackageByID(id uint) (*ServicePackage, error)
	UpdatePackage(pkg *ServicePackage) error

	AddServiceToPackage(packageServiceType *PackageServiceType) error
	GetPackageServices(packageID uint) ([]PackageServiceType, error)

	CreateCategory(category *ServiceCategory) error
	GetCategories() ([]ServiceCategory, error)

	CreateServiceHistory(history *VehicleServiceHistory) error
	GetVehicleHistory(vehicleID uint) ([]VehicleServiceHistory, error)
	GetVehicleHistoryWithDetails(vehicleID uint) ([]VehicleServiceHistory, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreatePackage(pkg *ServicePackage) error {
	return r.db.Create(pkg).Error
}

func (r *repository) GetPackages() ([]ServicePackage, error) {
	var packages []ServicePackage
	err := r.db.Where("is_active = ?", true).Find(&packages).Error
	return packages, err
}

func (r *repository) GetPackageByID(id uint) (*ServicePackage, error) {
	var pkg ServicePackage
	err := r.db.First(&pkg, id).Error
	return &pkg, err
}

func (r *repository) UpdatePackage(pkg *ServicePackage) error {
	return r.db.Save(pkg).Error
}

func (r *repository) AddServiceToPackage(packageServiceType *PackageServiceType) error {
	return r.db.Create(packageServiceType).Error
}

func (r *repository) GetPackageServices(packageID uint) ([]PackageServiceType, error) {
	var services []PackageServiceType
	err := r.db.Where("package_id = ?", packageID).Find(&services).Error
	return services, err
}

func (r *repository) CreateCategory(category *ServiceCategory) error {
	return r.db.Create(category).Error
}

func (r *repository) GetCategories() ([]ServiceCategory, error) {
	var categories []ServiceCategory
	err := r.db.Order("sort_order ASC").Find(&categories).Error
	return categories, err
}

func (r *repository) CreateServiceHistory(history *VehicleServiceHistory) error {
	return r.db.Create(history).Error
}

func (r *repository) GetVehicleHistory(vehicleID uint) ([]VehicleServiceHistory, error) {
	var history []VehicleServiceHistory
	err := r.db.Where("vehicle_id = ?", vehicleID).
		Order("completed_at DESC").
		Find(&history).Error
	return history, err
}

func (r *repository) GetVehicleHistoryWithDetails(vehicleID uint) ([]VehicleServiceHistory, error) {
	var history []VehicleServiceHistory
	err := r.db.Where("vehicle_id = ?", vehicleID).
		Order("completed_at DESC").
		Find(&history).Error
	return history, err
}
