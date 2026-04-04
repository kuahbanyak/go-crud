package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/pkg/pagination"
)

type RoleRepository interface {
	Create(ctx context.Context, role *entities.Role) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Role, error)
	GetByName(ctx context.Context, name string) (*entities.Role, error)
	GetAll(ctx context.Context) ([]*entities.Role, error)
	GetAllPaginated(ctx context.Context, pagParams pagination.Params, filterParams pagination.FilterParams) ([]*entities.Role, int64, error)
	GetActive(ctx context.Context) ([]*entities.Role, error)
	Update(ctx context.Context, role *entities.Role) error
	Delete(ctx context.Context, id uuid.UUID) error

	// User role assignments
	AssignRoleToUser(ctx context.Context, userID, roleID, assignedBy uuid.UUID) error
	RemoveRoleFromUser(ctx context.Context, userID, roleID uuid.UUID) error
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]*entities.Role, error)
	HasRole(ctx context.Context, userID uuid.UUID, roleName string) (bool, error)
	GetUsersByRole(ctx context.Context, roleID uuid.UUID) ([]*entities.User, error)
}
