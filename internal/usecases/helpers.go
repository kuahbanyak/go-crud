package usecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// Repository interface for common operations
type Repository[T any] interface {
	GetByID(ctx context.Context, id uuid.UUID) (*T, error)
}

// ValidateAndGetByID is a helper function to validate UUID and get entity by ID
func ValidateAndGetByID[T any](ctx context.Context, repo Repository[T], id uuid.UUID, entityName string) (*T, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid " + entityName + " ID")
	}

	entity, err := repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if entity == nil {
		return nil, errors.New(entityName + " not found")
	}

	return entity, nil
}

// CheckExistsOrError checks if entity exists, returns error if not
func CheckExistsOrError[T any](entity *T, err error, entityName string) error {
	if err != nil {
		return err
	}
	if entity == nil {
		return errors.New(entityName + " not found")
	}
	return nil
}
