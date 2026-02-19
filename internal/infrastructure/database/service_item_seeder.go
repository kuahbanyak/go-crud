package database

import (
	"context"

	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"gorm.io/gorm"
)

func SeedDefaultServiceItems(db *gorm.DB) error {
	ctx := context.Background()

	defaultServices := []entities.ServiceItem{
		// Engine Services
		{
			Name:            "Oil Change",
			Description:     "Engine oil and filter replacement",
			Category:        "Engine",
			EstimatedTime:   30,
			EstimatedCost:   150000,
			DisplayOrder:    1,
			RequiresBooking: true,
			IsActive:        true,
		},
		{
			Name:            "Engine Tune-Up",
			Description:     "Complete engine tune-up service",
			Category:        "Engine",
			EstimatedTime:   120,
			EstimatedCost:   500000,
			DisplayOrder:    2,
			RequiresBooking: true,
			IsActive:        true,
		},
		{
			Name:            "Air Filter Replacement",
			Description:     "Replace engine air filter",
			Category:        "Engine",
			EstimatedTime:   15,
			EstimatedCost:   75000,
			DisplayOrder:    3,
			RequiresBooking: false,
			IsActive:        true,
		},
		// Brake Services
		{
			Name:            "Brake Pad Replacement",
			Description:     "Replace front or rear brake pads",
			Category:        "Brakes",
			EstimatedTime:   60,
			EstimatedCost:   400000,
			DisplayOrder:    1,
			RequiresBooking: true,
			IsActive:        true,
		},
		{
			Name:            "Brake Fluid Change",
			Description:     "Complete brake fluid replacement",
			Category:        "Brakes",
			EstimatedTime:   30,
			EstimatedCost:   150000,
			DisplayOrder:    2,
			RequiresBooking: true,
			IsActive:        true,
		},
		{
			Name:            "Brake Inspection",
			Description:     "Comprehensive brake system inspection",
			Category:        "Brakes",
			EstimatedTime:   20,
			EstimatedCost:   50000,
			DisplayOrder:    3,
			RequiresBooking: false,
			IsActive:        true,
		},
		// Tire Services
		{
			Name:            "Tire Rotation",
			Description:     "Rotate all four tires",
			Category:        "Tires",
			EstimatedTime:   30,
			EstimatedCost:   100000,
			DisplayOrder:    1,
			RequiresBooking: true,
			IsActive:        true,
		},
		{
			Name:            "Tire Replacement",
			Description:     "Replace one or more tires",
			Category:        "Tires",
			EstimatedTime:   45,
			EstimatedCost:   600000,
			DisplayOrder:    2,
			RequiresBooking: true,
			IsActive:        true,
		},
		{
			Name:            "Wheel Alignment",
			Description:     "Four-wheel alignment service",
			Category:        "Tires",
			EstimatedTime:   60,
			EstimatedCost:   250000,
			DisplayOrder:    3,
			RequiresBooking: true,
			IsActive:        true,
		},
		// General Maintenance
		{
			Name:            "General Inspection",
			Description:     "Complete vehicle inspection",
			Category:        "General",
			EstimatedTime:   45,
			EstimatedCost:   100000,
			DisplayOrder:    1,
			RequiresBooking: true,
			IsActive:        true,
		},
		{
			Name:            "Battery Check",
			Description:     "Battery health check and replacement if needed",
			Category:        "Electrical",
			EstimatedTime:   20,
			EstimatedCost:   50000,
			DisplayOrder:    1,
			RequiresBooking: false,
			IsActive:        true,
		},
		{
			Name:            "AC Service",
			Description:     "Air conditioning system service and recharge",
			Category:        "AC",
			EstimatedTime:   60,
			EstimatedCost:   350000,
			DisplayOrder:    1,
			RequiresBooking: true,
			IsActive:        true,
		},
	}

	for _, service := range defaultServices {
		var existingService entities.ServiceItem
		result := db.WithContext(ctx).
			Where("name = ? AND category = ?", service.Name, service.Category).
			First(&existingService)

		if result.Error == gorm.ErrRecordNotFound {
			if err := db.WithContext(ctx).Create(&service).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
