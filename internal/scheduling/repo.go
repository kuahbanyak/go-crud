package scheduling

import (
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	CreateAvailability(availability *MechanicAvailability) error
	GetMechanicAvailability(mechanicID uint, date time.Time) ([]MechanicAvailability, error)
	UpdateAvailabilityStatus(id uint, status AvailabilityStatus) error

	CreateServiceType(serviceType *ServiceType) error
	GetServiceTypes() ([]ServiceType, error)
	GetServiceType(id uint) (*ServiceType, error)

	CreateReminder(reminder *MaintenanceReminder) error
	GetVehicleReminders(vehicleID uint) ([]MaintenanceReminder, error)
	GetDueReminders() ([]MaintenanceReminder, error)
	CompleteReminder(id uint) error

	AddToWaitlist(waitlist *BookingWaitlist) error
	GetWaitlistByDate(date time.Time) ([]BookingWaitlist, error)
	RemoveFromWaitlist(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateAvailability(availability *MechanicAvailability) error {
	return r.db.Create(availability).Error
}

func (r *repository) GetMechanicAvailability(mechanicID uint, date time.Time) ([]MechanicAvailability, error) {
	var availability []MechanicAvailability
	err := r.db.Where("mechanic_id = ? AND date = ?", mechanicID, date.Format("2006-01-02")).
		Order("start_time ASC").
		Find(&availability).Error
	return availability, err
}

func (r *repository) UpdateAvailabilityStatus(id uint, status AvailabilityStatus) error {
	return r.db.Model(&MechanicAvailability{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *repository) CreateServiceType(serviceType *ServiceType) error {
	return r.db.Create(serviceType).Error
}

func (r *repository) GetServiceTypes() ([]ServiceType, error) {
	var serviceTypes []ServiceType
	err := r.db.Find(&serviceTypes).Error
	return serviceTypes, err
}

func (r *repository) GetServiceType(id uint) (*ServiceType, error) {
	var serviceType ServiceType
	err := r.db.First(&serviceType, id).Error
	return &serviceType, err
}

func (r *repository) CreateReminder(reminder *MaintenanceReminder) error {
	return r.db.Create(reminder).Error
}

func (r *repository) GetVehicleReminders(vehicleID uint) ([]MaintenanceReminder, error) {
	var reminders []MaintenanceReminder
	err := r.db.Where("vehicle_id = ? AND is_completed = ?", vehicleID, false).
		Order("due_date ASC").
		Find(&reminders).Error
	return reminders, err
}

func (r *repository) GetDueReminders() ([]MaintenanceReminder, error) {
	var reminders []MaintenanceReminder
	today := time.Now().Format("2006-01-02")
	err := r.db.Where("due_date <= ? AND is_completed = ?", today, false).
		Find(&reminders).Error
	return reminders, err
}

func (r *repository) CompleteReminder(id uint) error {
	now := time.Now()
	return r.db.Model(&MaintenanceReminder{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_completed": true,
			"completed_at": &now,
		}).Error
}

func (r *repository) AddToWaitlist(waitlist *BookingWaitlist) error {
	return r.db.Create(waitlist).Error
}

func (r *repository) GetWaitlistByDate(date time.Time) ([]BookingWaitlist, error) {
	var waitlist []BookingWaitlist
	err := r.db.Where("preferred_date = ? AND is_notified = ?", date.Format("2006-01-02"), false).
		Find(&waitlist).Error
	return waitlist, err
}

func (r *repository) RemoveFromWaitlist(id uint) error {
	return r.db.Delete(&BookingWaitlist{}, id).Error
}
