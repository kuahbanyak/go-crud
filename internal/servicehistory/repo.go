package servicehistory

import "gorm.io/gorm"

type Repository interface {
    Create(s *ServiceRecord) error
    ListByVehicle(vehicle uint) ([]ServiceRecord, error)
}

type repo struct { db *gorm.DB }

func NewRepo(db *gorm.DB) Repository { return &repo{db: db} }

func (r *repo) Create(s *ServiceRecord) error {
    return r.db.Create(s).Error
}
func (r *repo) ListByVehicle(vehicle uint) ([]ServiceRecord, error) {
    var ss []ServiceRecord
    if err := r.db.Where("vehicle_id = ?", vehicle).Find(&ss).Error; err != nil {
        return nil, err
    }
    return ss, nil
}
