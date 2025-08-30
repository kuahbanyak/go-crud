package invoice

import "gorm.io/gorm"

type Repository interface {
	Create(i *Invoice) error
	Summary() (map[string]int64, error)

	CreateCustomBody(body *CustomInvoiceBody) error
	GetCustomBodies() ([]CustomInvoiceBody, error)
	GetCustomBodyByID(id string) (*CustomInvoiceBody, error)
	GetDefaultCustomBody() (*CustomInvoiceBody, error)
	UpdateCustomBody(body *CustomInvoiceBody) error
	DeleteCustomBody(id string) error
	SetDefaultCustomBody(id string) error
}

type repo struct{ db *gorm.DB }

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

func (r *repo) CreateCustomBody(body *CustomInvoiceBody) error {
	return r.db.Create(body).Error
}

func (r *repo) GetCustomBodies() ([]CustomInvoiceBody, error) {
	var bodies []CustomInvoiceBody
	err := r.db.Where("is_active = ?", true).Order("created_at DESC").Find(&bodies).Error
	return bodies, err
}

func (r *repo) GetCustomBodyByID(id string) (*CustomInvoiceBody, error) {
	var body CustomInvoiceBody
	err := r.db.Where("id = ?", id).First(&body).Error
	if err != nil {
		return nil, err
	}
	return &body, nil
}

func (r *repo) GetDefaultCustomBody() (*CustomInvoiceBody, error) {
	var body CustomInvoiceBody
	err := r.db.Where("is_default = ? AND is_active = ?", true, true).First(&body).Error
	if err != nil {
		return nil, err
	}
	return &body, nil
}

func (r *repo) UpdateCustomBody(body *CustomInvoiceBody) error {
	return r.db.Save(body).Error
}

func (r *repo) DeleteCustomBody(id string) error {
	return r.db.Where("id = ?", id).Delete(&CustomInvoiceBody{}).Error
}

func (r *repo) SetDefaultCustomBody(id string) error {

	tx := r.db.Begin()

	if err := tx.Model(&CustomInvoiceBody{}).Where("is_default = ?", true).Update("is_default", false).Error; err != nil {
		tx.Rollback()
		return err
	}

	if id != "" {
		if err := tx.Model(&CustomInvoiceBody{}).Where("id = ?", id).Update("is_default", true).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
