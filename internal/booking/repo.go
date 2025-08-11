package booking

import "gorm.io/gorm"

type Repository interface {
	Create(b *Booking) error
	ListByCustomer(customer uint) ([]Booking, error)
	UpdateStatus(id uint, status BookingStatus) error
	GetId(id uint) (*Booking, error)
}

type repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) Repository { return &repo{db: db} }

func (r *repo) Create(b *Booking) error {
	return r.db.Create(b).Error
}
func (r repo) GetId(id uint) (*Booking, error) {
	var b Booking
	if err := r.db.First(&b, id).Error; err != nil {
		return nil, err
	}
	return &b, nil

}
func (r *repo) ListByCustomer(customer uint) ([]Booking, error) {
	var bs []Booking
	if err := r.db.Where("customer_id = ?", customer).Find(&bs).Error; err != nil {
		return nil, err
	}
	return bs, nil
}
func (r *repo) UpdateStatus(id uint, status BookingStatus) error {
	return r.db.Model(&Booking{}).Where("id = ?", id).Update("status", status).Error
}
