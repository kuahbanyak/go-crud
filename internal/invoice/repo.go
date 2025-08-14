package invoice

import "gorm.io/gorm"

type Repository interface {
    Create(i *Invoice) error
    Summary() (map[string]int64, error)
}

type repo struct { db *gorm.DB }

func NewRepo(db *gorm.DB) Repository { return &repo{db: db} }

func (r *repo) Create(i *Invoice) error {
    return r.db.Create(i).Error
}

func (r *repo) Summary() (map[string]int64, error) {
    var total int64
    if err := r.db.Model(&Invoice{}).Count(&total).Error; err != nil {
        return nil, err
    }
    return map[string]int64{"count": total}, nil
}
