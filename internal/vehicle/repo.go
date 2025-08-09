package vehicle

import "gorm.io/gorm"

type Repository interface {
	Create(v *Vehicle) error
	ListByOwner(owner uint) ([]Vehicle, error)
	Get(id uint) (*Vehicle, error)
	Update(v *Vehicle) error
	Delete(id uint) error
}

type repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) Repository { return &repo{db: db} }

func (r *repo) Create(v *Vehicle) error {
	return r.db.Create(v).Error
}
func (r *repo) ListByOwner(owner uint) ([]Vehicle, error) {
	var vs []Vehicle
	if err := r.db.Where("owner_id = ?", owner).Find(&vs).Error; err != nil {
		return nil, err
	}
	return vs, nil
}
func (r *repo) Get(id uint) (*Vehicle, error) {
	var v Vehicle
	if err := r.db.First(&v, id).Error; err != nil {
		return nil, err
	}
	return &v, nil
}
func (r *repo) Update(v *Vehicle) error {
	return r.db.Save(v).Error
}
func (r *repo) Delete(id uint) error {
	return r.db.Delete(&Vehicle{}, id).Error
}
