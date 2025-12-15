package mssql

import (
	"context"
	"errors"
	"time"

	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/domain/repositories"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
	"github.com/kuahbanyak/go-crud/pkg/pagination"
	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) repositories.RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(ctx context.Context, role *entities.Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *roleRepository) GetByID(ctx context.Context, id types.MSSQLUUID) (*entities.Role, error) {
	var role entities.Role
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) GetByName(ctx context.Context, name string) (*entities.Role, error) {
	var role entities.Role
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) GetAll(ctx context.Context) ([]*entities.Role, error) {
	var roles []*entities.Role
	err := r.db.WithContext(ctx).Order("name ASC").Find(&roles).Error
	return roles, err
}

func (r *roleRepository) GetAllPaginated(ctx context.Context, pagParams pagination.Params, filterParams pagination.FilterParams) ([]*entities.Role, int64, error) {
	var roles []*entities.Role
	var total int64

	// Base query
	query := r.db.WithContext(ctx).Model(&entities.Role{})

	// Apply filters
	if filterParams.Search != "" {
		query = pagination.ApplySearch(query, filterParams.Search, "name", "display_name", "description")
	}

	if filterParams.Status != "" {
		if filterParams.Status == "active" {
			query = query.Where("is_active = ?", true)
		} else if filterParams.Status == "inactive" {
			query = query.Where("is_active = ?", false)
		}
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	query = pagParams.Apply(query)

	// Execute query
	if err := query.Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

func (r *roleRepository) GetActive(ctx context.Context) ([]*entities.Role, error) {
	var roles []*entities.Role
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Order("name ASC").Find(&roles).Error
	return roles, err
}

func (r *roleRepository) Update(ctx context.Context, role *entities.Role) error {
	return r.db.WithContext(ctx).Save(role).Error
}

func (r *roleRepository) Delete(ctx context.Context, id types.MSSQLUUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.Role{}).Error
}

// AssignRoleToUser assigns a role to a user
func (r *roleRepository) AssignRoleToUser(ctx context.Context, userID, roleID, assignedBy types.MSSQLUUID) error {
	// Check if already assigned
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.UserRole{}).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Count(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("role already assigned to user")
	}

	// Create user role assignment
	userRole := &entities.UserRole{
		UserID:     userID,
		RoleID:     roleID,
		AssignedBy: assignedBy,
		AssignedAt: time.Now(),
	}

	return r.db.WithContext(ctx).Create(userRole).Error
}

// RemoveRoleFromUser removes a role from a user
func (r *roleRepository) RemoveRoleFromUser(ctx context.Context, userID, roleID types.MSSQLUUID) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&entities.UserRole{}).Error
}

// GetUserRoles gets all roles assigned to a user
func (r *roleRepository) GetUserRoles(ctx context.Context, userID types.MSSQLUUID) ([]*entities.Role, error) {
	var roles []*entities.Role
	err := r.db.WithContext(ctx).
		Table("roles").
		Joins("INNER JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Where("roles.is_active = ?", true).
		Find(&roles).Error
	return roles, err
}

// HasRole checks if a user has a specific role
func (r *roleRepository) HasRole(ctx context.Context, userID types.MSSQLUUID, roleName string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("user_roles").
		Joins("INNER JOIN roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ? AND roles.name = ? AND roles.is_active = ?", userID, roleName, true).
		Count(&count).Error
	return count > 0, err
}

// GetUsersByRole gets all users with a specific role
func (r *roleRepository) GetUsersByRole(ctx context.Context, roleID types.MSSQLUUID) ([]*entities.User, error) {
	var users []*entities.User
	err := r.db.WithContext(ctx).
		Preload("Roles").
		Table("users").
		Joins("INNER JOIN user_roles ON user_roles.user_id = users.id").
		Where("user_roles.role_id = ?", roleID).
		Find(&users).Error
	return users, err
}
