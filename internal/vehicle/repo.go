package vehicle

import "gorm.io/gorm"

type Repository interface {
	Create(v *Vehicle) error
	ListByOwner(owner string) ([]Vehicle, error)
	Get(id string) (*Vehicle, error)
	Update(v *Vehicle) error
	Delete(id string) error
}

type repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) Repository { return &repo{db: db} }

func (r *repo) Create(v *Vehicle) error {
	return r.db.Create(v).Error
}
func (r *repo) ListByOwner(owner string) ([]Vehicle, error) {
	var vs []Vehicle
	if err := r.db.Where("owner_id = ?", owner).Find(&vs).Error; err != nil {
		return nil, err
	}
	return vs, nil
}
func (r *repo) Get(id string) (*Vehicle, error) {
	var v Vehicle
	if err := r.db.Where("id = ?", id).First(&v).Error; err != nil {
		return nil, err
	}
	return &v, nil
}
func (r *repo) Update(v *Vehicle) error {
	return r.db.Save(v).Error
}
func (r *repo) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&Vehicle{}).Error
}
