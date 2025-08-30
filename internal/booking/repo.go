package booking

import "gorm.io/gorm"

type Repository interface {
	Create(b *Booking) error
	ListByCustomer(customer string) ([]Booking, error)
	UpdateStatus(id string, status BookingStatus) error
	GetId(id string) (*Booking, error)
}

type repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) Repository { return &repo{db: db} }

func (r *repo) Create(b *Booking) error {
	return r.db.Create(b).Error
}
func (r repo) GetId(id string) (*Booking, error) {
	var b Booking
	if err := r.db.Where("id = ?", id).First(&b).Error; err != nil {
		return nil, err
	}
	return &b, nil

}
func (r *repo) ListByCustomer(customer string) ([]Booking, error) {
	var bs []Booking
	if err := r.db.Where("customer_id = ?", customer).Find(&bs).Error; err != nil {
		return nil, err
	}
	return bs, nil
}
func (r *repo) UpdateStatus(id string, status BookingStatus) error {
	return r.db.Model(&Booking{}).Where("id = ?", id).Update("status", status).Error
}
