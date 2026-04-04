package usecases_test

import (
	"context"
	"testing"

	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/google/uuid"
	"github.com/kuahbanyak/go-crud/internal/usecases"
	"github.com/kuahbanyak/go-crud/pkg/pagination"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRoleUsecase_CreateRole_Success(t *testing.T) {
	mockRoleRepo := new(MockRoleRepository)
	mockUserRepo := new(MockUserRepository)
	usecase := usecases.NewRoleUsecase(mockRoleRepo, mockUserRepo)
	ctx := context.Background()

	role := &entities.Role{
		Name:        "admin",
		DisplayName: "Administrator",
		Description: "Admin role",
		IsActive:    true,
	}

	mockRoleRepo.On("GetByName", ctx, "admin").Return(nil, nil)
	mockRoleRepo.On("Create", ctx, role).Return(nil)

	err := usecase.CreateRole(ctx, role)

	assert.NoError(t, err)
	mockRoleRepo.AssertExpectations(t)
}

func TestRoleUsecase_CreateRole_AlreadyExists(t *testing.T) {
	mockRoleRepo := new(MockRoleRepository)
	mockUserRepo := new(MockUserRepository)
	usecase := usecases.NewRoleUsecase(mockRoleRepo, mockUserRepo)
	ctx := context.Background()

	role := &entities.Role{Name: "admin"}
	existingRole := &entities.Role{Name: "admin"}

	mockRoleRepo.On("GetByName", ctx, "admin").Return(existingRole, nil)

	err := usecase.CreateRole(ctx, role)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
	mockRoleRepo.AssertExpectations(t)
}

func TestRoleUsecase_GetRoleByID_Success(t *testing.T) {
	mockRoleRepo := new(MockRoleRepository)
	mockUserRepo := new(MockUserRepository)
	usecase := usecases.NewRoleUsecase(mockRoleRepo, mockUserRepo)
	ctx := context.Background()

	roleID := uuid.UUID{}
	role := &entities.Role{
		ID:   roleID,
		Name: "admin",
	}

	mockRoleRepo.On("GetByID", ctx, roleID).Return(role, nil)

	result, err := usecase.GetRoleByID(ctx, roleID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "admin", result.Name)
	mockRoleRepo.AssertExpectations(t)
}

func TestRoleUsecase_GetRoleByID_NotFound(t *testing.T) {
	mockRoleRepo := new(MockRoleRepository)
	mockUserRepo := new(MockUserRepository)
	usecase := usecases.NewRoleUsecase(mockRoleRepo, mockUserRepo)
	ctx := context.Background()

	roleID := uuid.UUID{}

	mockRoleRepo.On("GetByID", ctx, roleID).Return(nil, nil)

	result, err := usecase.GetRoleByID(ctx, roleID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "role not found")
	mockRoleRepo.AssertExpectations(t)
}

func TestRoleUsecase_GetAllRoles_Success(t *testing.T) {
	mockRoleRepo := new(MockRoleRepository)
	mockUserRepo := new(MockUserRepository)
	usecase := usecases.NewRoleUsecase(mockRoleRepo, mockUserRepo)
	ctx := context.Background()

	roles := []*entities.Role{
		{Name: "admin"},
		{Name: "customer"},
	}

	mockRoleRepo.On("GetAll", ctx).Return(roles, nil)

	result, err := usecase.GetAllRoles(ctx)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockRoleRepo.AssertExpectations(t)
}

func TestRoleUsecase_GetAllRolesPaginated_Success(t *testing.T) {
	mockRoleRepo := new(MockRoleRepository)
	mockUserRepo := new(MockUserRepository)
	usecase := usecases.NewRoleUsecase(mockRoleRepo, mockUserRepo)
	ctx := context.Background()

	roles := []*entities.Role{
		{Name: "admin"},
	}

	pagParams := pagination.Params{Page: 1, PageSize: 10}
	filterParams := pagination.FilterParams{}

	mockRoleRepo.On("GetAllPaginated", ctx, pagParams, filterParams).Return(roles, 1, nil)

	result, total, err := usecase.GetAllRolesPaginated(ctx, pagParams, filterParams)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), total)
	mockRoleRepo.AssertExpectations(t)
}

func TestRoleUsecase_GetActiveRoles_Success(t *testing.T) {
	mockRoleRepo := new(MockRoleRepository)
	mockUserRepo := new(MockUserRepository)
	usecase := usecases.NewRoleUsecase(mockRoleRepo, mockUserRepo)
	ctx := context.Background()

	roles := []*entities.Role{
		{Name: "admin", IsActive: true},
	}

	mockRoleRepo.On("GetActive", ctx).Return(roles, nil)

	result, err := usecase.GetActiveRoles(ctx)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	mockRoleRepo.AssertExpectations(t)
}

func TestRoleUsecase_UpdateRole_Success(t *testing.T) {
	mockRoleRepo := new(MockRoleRepository)
	mockUserRepo := new(MockUserRepository)
	usecase := usecases.NewRoleUsecase(mockRoleRepo, mockUserRepo)
	ctx := context.Background()

	roleID := uuid.UUID{}
	existingRole := &entities.Role{
		ID:          roleID,
		Name:        "admin",
		DisplayName: "Old Admin",
		IsActive:    true,
	}

	updateData := &entities.Role{
		DisplayName: "New Admin",
		Description: "Updated description",
		IsActive:    false,
	}

	mockRoleRepo.On("GetByID", ctx, roleID).Return(existingRole, nil)
	mockRoleRepo.On("Update", ctx, mock.AnythingOfType("*entities.Role")).Return(nil)

	result, err := usecase.UpdateRole(ctx, roleID, updateData)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "New Admin", result.DisplayName)
	assert.Equal(t, "Updated description", result.Description)
	assert.False(t, result.IsActive)
	mockRoleRepo.AssertExpectations(t)
}

func TestRoleUsecase_UpdateRole_NotFound(t *testing.T) {
	mockRoleRepo := new(MockRoleRepository)
	mockUserRepo := new(MockUserRepository)
	usecase := usecases.NewRoleUsecase(mockRoleRepo, mockUserRepo)
	ctx := context.Background()

	roleID := uuid.UUID{}
	updateData := &entities.Role{DisplayName: "New Name"}

	mockRoleRepo.On("GetByID", ctx, roleID).Return(nil, nil)

	result, err := usecase.UpdateRole(ctx, roleID, updateData)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "role not found")
	mockRoleRepo.AssertExpectations(t)
}

func TestRoleUsecase_DeleteRole_Success(t *testing.T) {
	mockRoleRepo := new(MockRoleRepository)
	mockUserRepo := new(MockUserRepository)
	usecase := usecases.NewRoleUsecase(mockRoleRepo, mockUserRepo)
	ctx := context.Background()

	roleID := uuid.UUID{}
	role := &entities.Role{ID: roleID}

	mockRoleRepo.On("GetByID", ctx, roleID).Return(role, nil)
	mockRoleRepo.On("Delete", ctx, roleID).Return(nil)

	err := usecase.DeleteRole(ctx, roleID)

	assert.NoError(t, err)
	mockRoleRepo.AssertExpectations(t)
}

func TestRoleUsecase_DeleteRole_NotFound(t *testing.T) {
	mockRoleRepo := new(MockRoleRepository)
	mockUserRepo := new(MockUserRepository)
	usecase := usecases.NewRoleUsecase(mockRoleRepo, mockUserRepo)
	ctx := context.Background()

	roleID := uuid.UUID{}

	mockRoleRepo.On("GetByID", ctx, roleID).Return(nil, nil)

	err := usecase.DeleteRole(ctx, roleID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "role not found")
	mockRoleRepo.AssertExpectations(t)
}

func TestRoleUsecase_AssignRoleToUser_Success(t *testing.T) {
	mockRoleRepo := new(MockRoleRepository)
	mockUserRepo := new(MockUserRepository)
	usecase := usecases.NewRoleUsecase(mockRoleRepo, mockUserRepo)
	ctx := context.Background()

	userID := uuid.UUID{}
	roleID := uuid.UUID{}
	assignedBy := uuid.UUID{}

	user := &entities.User{ID: userID}
	role := &entities.Role{ID: roleID, IsActive: true}

	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil)
	mockRoleRepo.On("GetByID", ctx, roleID).Return(role, nil)
	mockRoleRepo.On("AssignRoleToUser", ctx, userID, roleID, assignedBy).Return(nil)

	err := usecase.AssignRoleToUser(ctx, userID, roleID, assignedBy)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
}

func TestRoleUsecase_AssignRoleToUser_UserNotFound(t *testing.T) {
	mockRoleRepo := new(MockRoleRepository)
	mockUserRepo := new(MockUserRepository)
	usecase := usecases.NewRoleUsecase(mockRoleRepo, mockUserRepo)
	ctx := context.Background()

	userID := uuid.UUID{}
	roleID := uuid.UUID{}
	assignedBy := uuid.UUID{}

	mockUserRepo.On("GetByID", ctx, userID).Return(nil, nil)

	err := usecase.AssignRoleToUser(ctx, userID, roleID, assignedBy)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
	mockUserRepo.AssertExpectations(t)
}

func TestRoleUsecase_AssignRoleToUser_RoleNotFound(t *testing.T) {
	mockRoleRepo := new(MockRoleRepository)
	mockUserRepo := new(MockUserRepository)
	usecase := usecases.NewRoleUsecase(mockRoleRepo, mockUserRepo)
	ctx := context.Background()

	userID := uuid.UUID{}
	roleID := uuid.UUID{}
	assignedBy := uuid.UUID{}

	user := &entities.User{ID: userID}

	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil)
	mockRoleRepo.On("GetByID", ctx, roleID).Return(nil, nil)

	err := usecase.AssignRoleToUser(ctx, userID, roleID, assignedBy)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "role not found")
	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
}

func TestRoleUsecase_AssignRoleToUser_InactiveRole(t *testing.T) {
	mockRoleRepo := new(MockRoleRepository)
	mockUserRepo := new(MockUserRepository)
	usecase := usecases.NewRoleUsecase(mockRoleRepo, mockUserRepo)
	ctx := context.Background()

	userID := uuid.UUID{}
	roleID := uuid.UUID{}
	assignedBy := uuid.UUID{}

	user := &entities.User{ID: userID}
	role := &entities.Role{ID: roleID, IsActive: false}

	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil)
	mockRoleRepo.On("GetByID", ctx, roleID).Return(role, nil)

	err := usecase.AssignRoleToUser(ctx, userID, roleID, assignedBy)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "inactive role")
	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
}

func TestRoleUsecase_RemoveRoleFromUser_Success(t *testing.T) {
	mockRoleRepo := new(MockRoleRepository)
	mockUserRepo := new(MockUserRepository)
	usecase := usecases.NewRoleUsecase(mockRoleRepo, mockUserRepo)
	ctx := context.Background()

	userID := uuid.UUID{}
	roleID := uuid.UUID{}

	user := &entities.User{ID: userID}
	role := &entities.Role{ID: roleID}

	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil)
	mockRoleRepo.On("GetByID", ctx, roleID).Return(role, nil)
	mockRoleRepo.On("RemoveRoleFromUser", ctx, userID, roleID).Return(nil)

	err := usecase.RemoveRoleFromUser(ctx, userID, roleID)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
}

func TestRoleUsecase_RemoveRoleFromUser_UserNotFound(t *testing.T) {
	mockRoleRepo := new(MockRoleRepository)
	mockUserRepo := new(MockUserRepository)
	usecase := usecases.NewRoleUsecase(mockRoleRepo, mockUserRepo)
	ctx := context.Background()

	userID := uuid.UUID{}
	roleID := uuid.UUID{}

	mockUserRepo.On("GetByID", ctx, userID).Return(nil, nil)

	err := usecase.RemoveRoleFromUser(ctx, userID, roleID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
	mockUserRepo.AssertExpectations(t)
}

func TestRoleUsecase_GetUserRoles_Success(t *testing.T) {
	mockRoleRepo := new(MockRoleRepository)
	mockUserRepo := new(MockUserRepository)
	usecase := usecases.NewRoleUsecase(mockRoleRepo, mockUserRepo)
	ctx := context.Background()

	userID := uuid.UUID{}
	user := &entities.User{ID: userID}
	roles := []*entities.Role{
		{Name: "admin"},
		{Name: "customer"},
	}

	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil)
	mockRoleRepo.On("GetUserRoles", ctx, userID).Return(roles, nil)

	result, err := usecase.GetUserRoles(ctx, userID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
}

func TestRoleUsecase_GetUserRoles_UserNotFound(t *testing.T) {
	mockRoleRepo := new(MockRoleRepository)
	mockUserRepo := new(MockUserRepository)
	usecase := usecases.NewRoleUsecase(mockRoleRepo, mockUserRepo)
	ctx := context.Background()

	userID := uuid.UUID{}

	mockUserRepo.On("GetByID", ctx, userID).Return(nil, nil)

	result, err := usecase.GetUserRoles(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "user not found")
	mockUserRepo.AssertExpectations(t)
}

func TestRoleUsecase_HasRole_Success(t *testing.T) {
	mockRoleRepo := new(MockRoleRepository)
	mockUserRepo := new(MockUserRepository)
	usecase := usecases.NewRoleUsecase(mockRoleRepo, mockUserRepo)
	ctx := context.Background()

	userID := uuid.UUID{}

	mockRoleRepo.On("HasRole", ctx, userID, "admin").Return(true, nil)

	result, err := usecase.HasRole(ctx, userID, "admin")

	assert.NoError(t, err)
	assert.True(t, result)
	mockRoleRepo.AssertExpectations(t)
}

func TestRoleUsecase_GetUsersByRole_Success(t *testing.T) {
	mockRoleRepo := new(MockRoleRepository)
	mockUserRepo := new(MockUserRepository)
	usecase := usecases.NewRoleUsecase(mockRoleRepo, mockUserRepo)
	ctx := context.Background()

	roleID := uuid.UUID{}
	role := &entities.Role{ID: roleID}
	users := []*entities.User{
		{Email: "user1@example.com"},
		{Email: "user2@example.com"},
	}

	mockRoleRepo.On("GetByID", ctx, roleID).Return(role, nil)
	mockRoleRepo.On("GetUsersByRole", ctx, roleID).Return(users, nil)

	result, err := usecase.GetUsersByRole(ctx, roleID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockRoleRepo.AssertExpectations(t)
}

func TestRoleUsecase_GetUsersByRole_RoleNotFound(t *testing.T) {
	mockRoleRepo := new(MockRoleRepository)
	mockUserRepo := new(MockUserRepository)
	usecase := usecases.NewRoleUsecase(mockRoleRepo, mockUserRepo)
	ctx := context.Background()

	roleID := uuid.UUID{}

	mockRoleRepo.On("GetByID", ctx, roleID).Return(nil, nil)

	result, err := usecase.GetUsersByRole(ctx, roleID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "role not found")
	mockRoleRepo.AssertExpectations(t)
}
