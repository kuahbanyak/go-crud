package repositories

import (
	"context"

	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
)

type SettingRepository interface {
	Create(ctx context.Context, setting *entities.Setting) error
	GetByKey(ctx context.Context, key string) (*entities.Setting, error)
	GetByCategory(ctx context.Context, category string) ([]*entities.Setting, error)
	GetAll(ctx context.Context) ([]*entities.Setting, error)
	GetPublic(ctx context.Context) ([]*entities.Setting, error)
	Update(ctx context.Context, setting *entities.Setting) error
	Delete(ctx context.Context, id types.MSSQLUUID) error
	SeedDefaults(ctx context.Context) error
}
