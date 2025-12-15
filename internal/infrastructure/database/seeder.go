package database

import (
	"context"
	"log"

	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"gorm.io/gorm"
)

// SeedDefaultRoles seeds the default roles into the database
func SeedDefaultRoles(db *gorm.DB) error {
	ctx := context.Background()

	defaultRoles := []entities.Role{
		{
			Name:        "admin",
			DisplayName: "Administrator",
			Description: "Full system access and management capabilities",
			IsActive:    true,
		},
		{
			Name:        "mechanic",
			DisplayName: "Mechanic",
			Description: "Can perform maintenance tasks and update service status",
			IsActive:    true,
		},
		{
			Name:        "customer",
			DisplayName: "Customer",
			Description: "Basic user with access to customer features",
			IsActive:    true,
		},
		{
			Name:        "manager",
			DisplayName: "Manager",
			Description: "Can manage operations and view reports",
			IsActive:    true,
		},
	}

	for _, role := range defaultRoles {
		var existingRole entities.Role
		result := db.WithContext(ctx).Where("name = ?", role.Name).First(&existingRole)

		if result.Error == gorm.ErrRecordNotFound {
			// Role doesn't exist, create it
			if err := db.WithContext(ctx).Create(&role).Error; err != nil {
				log.Printf("Failed to seed role %s: %v", role.Name, err)
				return err
			}
			log.Printf("Seeded role: %s", role.Name)
		} else if result.Error != nil {
			// Some other error occurred
			log.Printf("Error checking role %s: %v", role.Name, result.Error)
			return result.Error
		}
		// Role already exists, skip
	}

	log.Println("Default roles seeding completed")
	return nil
}
