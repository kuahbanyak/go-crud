package user

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(u *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id uint) (*User, error)
	Update(u *User) error
	FindAll() ([]*User, error)
}

type repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) Repository { return &repo{db: db} }

func (r *repo) Create(u *User) error {
	return r.db.Create(u).Error
}
func (r *repo) FindByEmail(email string) (*User, error) {
	var u User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
func (r *repo) FindByID(id uint) (*User, error) {
	var u User
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
func (r *repo) Update(u *User) error {
	if u.ID == 0 {
		return errors.New("missing id")
	}
	return r.db.Save(u).Error
}
func (r *repo) FindAll() ([]*User, error) {
	var users []*User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
