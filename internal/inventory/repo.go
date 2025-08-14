package inventory

import "gorm.io/gorm"

type Repository interface {
    Create(p *Part) error
    List() ([]Part, error)
    Update(p *Part) error
}

type repo struct { db *gorm.DB }

func NewRepo(db *gorm.DB) Repository { return &repo{db: db} }

func (r *repo) Create(p *Part) error {
    return r.db.Create(p).Error
}
func (r *repo) List() ([]Part, error) {
    var ps []Part
    if err := r.db.Find(&ps).Error; err != nil {
        return nil, err
    }
    return ps, nil
}
func (r *repo) Update(p *Part) error {
    return r.db.Save(p).Error
}
