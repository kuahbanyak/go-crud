package dashboard

import (
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	UpdateVehicleHealth(health *VehicleHealthScore) error
	GetVehicleHealth(vehicleID uint) (*VehicleHealthScore, error)

	CreateRecommendation(rec *MaintenanceRecommendation) error
	GetVehicleRecommendations(vehicleID uint) ([]MaintenanceRecommendation, error)
	GetCustomerRecommendations(customerID uint) ([]MaintenanceRecommendation, error)

	UpdateBudget(budget *CustomerBudget) error
	GetCustomerBudget(customerID uint) (*CustomerBudget, error)
	UpdateSpending(customerID uint, amount int) error

	GetCustomerDashboard(customerID uint) (map[string]interface{}, error)
	GetVehicleDashboard(vehicleID uint) (map[string]interface{}, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) UpdateVehicleHealth(health *VehicleHealthScore) error {
	return r.db.Save(health).Error
}

func (r *repository) GetVehicleHealth(vehicleID uint) (*VehicleHealthScore, error) {
	var health VehicleHealthScore
	err := r.db.Where("vehicle_id = ?", vehicleID).First(&health).Error
	return &health, err
}

func (r *repository) CreateRecommendation(rec *MaintenanceRecommendation) error {
	return r.db.Create(rec).Error
}

func (r *repository) GetVehicleRecommendations(vehicleID uint) ([]MaintenanceRecommendation, error) {
	var recommendations []MaintenanceRecommendation
	err := r.db.Where("vehicle_id = ? AND is_completed = ?", vehicleID, false).
		Order("priority DESC, due_in_days ASC").
		Find(&recommendations).Error
	return recommendations, err
}

func (r *repository) GetCustomerRecommendations(customerID uint) ([]MaintenanceRecommendation, error) {
	var recommendations []MaintenanceRecommendation
	err := r.db.Joins("JOIN vehicles ON vehicles.id = maintenance_recommendations.vehicle_id").
		Where("vehicles.owner_id = ? AND maintenance_recommendations.is_completed = ?", customerID, false).
		Order("maintenance_recommendations.priority DESC").
		Find(&recommendations).Error
	return recommendations, err
}

func (r *repository) UpdateBudget(budget *CustomerBudget) error {
	return r.db.Save(budget).Error
}

func (r *repository) GetCustomerBudget(customerID uint) (*CustomerBudget, error) {
	var budget CustomerBudget
	err := r.db.Where("customer_id = ?", customerID).First(&budget).Error
	return &budget, err
}

func (r *repository) UpdateSpending(customerID uint, amount int) error {
	return r.db.Model(&CustomerBudget{}).
		Where("customer_id = ?", customerID).
		Updates(map[string]interface{}{
			"spent_this_month": gorm.Expr("spent_this_month + ?", amount),
			"spent_this_year":  gorm.Expr("spent_this_year + ?", amount),
		}).Error
}

func (r *repository) GetCustomerDashboard(customerID uint) (map[string]interface{}, error) {
	dashboard := make(map[string]interface{})

	var vehicleCount int64
	r.db.Model(&struct {
		OwnerID uint `gorm:"column:owner_id"`
	}{}).Table("vehicles").Where("owner_id = ?", customerID).Count(&vehicleCount)

	var activeBookings int64
	r.db.Model(&struct {
		CustomerID uint   `gorm:"column:customer_id"`
		Status     string `gorm:"column:status"`
	}{}).Table("bookings").Where("customer_id = ? AND status IN ?", customerID, []string{"scheduled", "in_progress"}).Count(&activeBookings)

	var pendingRecs int64
	r.db.Model(&MaintenanceRecommendation{}).
		Joins("JOIN vehicles ON vehicles.id = maintenance_recommendations.vehicle_id").
		Where("vehicles.owner_id = ? AND maintenance_recommendations.is_completed = ?", customerID, false).
		Count(&pendingRecs)

	dashboard["vehicle_count"] = vehicleCount
	dashboard["active_bookings"] = activeBookings
	dashboard["pending_recommendations"] = pendingRecs

	return dashboard, nil
}

func (r *repository) GetVehicleDashboard(vehicleID uint) (map[string]interface{}, error) {
	dashboard := make(map[string]interface{})

	var serviceCount int64
	r.db.Model(&struct {
		VehicleID uint `gorm:"column:vehicle_id"`
	}{}).Table("vehicle_service_histories").Where("vehicle_id = ?", vehicleID).Count(&serviceCount)

	var lastService struct {
		CompletedAt time.Time `gorm:"column:completed_at"`
	}
	r.db.Model(&struct {
		VehicleID   uint      `gorm:"column:vehicle_id"`
		CompletedAt time.Time `gorm:"column:completed_at"`
	}{}).Table("vehicle_service_histories").
		Where("vehicle_id = ?", vehicleID).
		Order("completed_at DESC").
		First(&lastService)

	dashboard["total_services"] = serviceCount
	dashboard["last_service"] = lastService.CompletedAt

	return dashboard, nil
}
