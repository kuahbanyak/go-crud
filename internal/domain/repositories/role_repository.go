package repositories

import (
	"context"

	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
	"github.com/kuahbanyak/go-crud/pkg/pagination"
)

type RoleRepository interface {
	Create(ctx context.Context, role *entities.Role) error
	GetByID(ctx context.Context, id types.MSSQLUUID) (*entities.Role, error)
	GetByName(ctx context.Context, name string) (*entities.Role, error)
	GetAll(ctx context.Context) ([]*entities.Role, error)
	GetAllPaginated(ctx context.Context, pagParams pagination.Params, filterParams pagination.FilterParams) ([]*entities.Role, int64, error)
	GetActive(ctx context.Context) ([]*entities.Role, error)
	Update(ctx context.Context, role *entities.Role) error
	Delete(ctx context.Context, id types.MSSQLUUID) error

	// User role assignments
	AssignRoleToUser(ctx context.Context, userID, roleID, assignedBy types.MSSQLUUID) error
	RemoveRoleFromUser(ctx context.Context, userID, roleID types.MSSQLUUID) error
	GetUserRoles(ctx context.Context, userID types.MSSQLUUID) ([]*entities.Role, error)
	HasRole(ctx context.Context, userID types.MSSQLUUID, roleName string) (bool, error)
	GetUsersByRole(ctx context.Context, roleID types.MSSQLUUID) ([]*entities.User, error)
}
