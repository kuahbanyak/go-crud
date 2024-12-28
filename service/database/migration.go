package database

import (
	"go-crud/entity"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	err := db.AutoMigrate(&entity.Account{}, &entity.Product{})
	if err != nil {
		return
	}
}
