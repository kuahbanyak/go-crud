package usecases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
	"github.com/kuahbanyak/go-crud/internal/usecases"
	"github.com/kuahbanyak/go-crud/pkg/pagination"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock for UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id types.MSSQLUUID) (*entities.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) GetByRole(ctx context.Context, role string) ([]*entities.User, error) {
	args := m.Called(ctx, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.User), args.Error(1)
}

func (m *MockUserRepository) GetAll(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.User), args.Error(1)
}

func (m *MockUserRepository) GetAllPaginated(ctx context.Context, pagParams pagination.Params, filterParams pagination.FilterParams) ([]*entities.User, int64, error) {
	args := m.Called(ctx, pagParams, filterParams)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*entities.User), int64(args.Int(1)), args.Error(2)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id types.MSSQLUUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) Count(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

// MockRoleRepository is a mock for RoleRepository
type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) GetByName(ctx context.Context, name string) (*entities.Role, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Role), args.Error(1)
}

func (m *MockRoleRepository) GetByID(ctx context.Context, id types.MSSQLUUID) (*entities.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Role), args.Error(1)
}

func (m *MockRoleRepository) AssignRoleToUser(ctx context.Context, userID, roleID, assignedBy types.MSSQLUUID) error {
	args := m.Called(ctx, userID, roleID, assignedBy)
	return args.Error(0)
}

func (m *MockRoleRepository) Create(ctx context.Context, role *entities.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockRoleRepository) GetAll(ctx context.Context) ([]*entities.Role, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Role), args.Error(1)
}

func (m *MockRoleRepository) GetAllPaginated(ctx context.Context, pagParams pagination.Params, filterParams pagination.FilterParams) ([]*entities.Role, int64, error) {
	args := m.Called(ctx, pagParams, filterParams)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*entities.Role), int64(args.Int(1)), args.Error(2)
}

func (m *MockRoleRepository) GetActive(ctx context.Context) ([]*entities.Role, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Role), args.Error(1)
}

func (m *MockRoleRepository) Update(ctx context.Context, role *entities.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockRoleRepository) Delete(ctx context.Context, id types.MSSQLUUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRoleRepository) RemoveRoleFromUser(ctx context.Context, userID, roleID types.MSSQLUUID) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

func (m *MockRoleRepository) GetUserRoles(ctx context.Context, userID types.MSSQLUUID) ([]*entities.Role, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Role), args.Error(1)
}

func (m *MockRoleRepository) HasRole(ctx context.Context, userID types.MSSQLUUID, roleName string) (bool, error) {
	args := m.Called(ctx, userID, roleName)
	return args.Bool(0), args.Error(1)
}

func (m *MockRoleRepository) GetUsersByRole(ctx context.Context, roleID types.MSSQLUUID) ([]*entities.User, error) {
	args := m.Called(ctx, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.User), args.Error(1)
}

// MockAuthService is a mock for AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) ComparePassword(hashedPassword, password string) error {
	args := m.Called(hashedPassword, password)
	return args.Error(0)
}

func (m *MockAuthService) GenerateToken(userID types.MSSQLUUID, role string) (string, error) {
	args := m.Called(userID, role)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) ValidateToken(token string) (types.MSSQLUUID, string, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return types.MSSQLUUID{}, "", args.Error(2)
	}
	return args.Get(0).(types.MSSQLUUID), args.String(1), args.Error(2)
}

func TestUserUsecase_Register_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockAuthService := new(MockAuthService)
	usecase := usecases.NewUserUsecase(mockUserRepo, mockRoleRepo, mockAuthService)
	ctx := context.Background()

	user := &entities.User{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	customerRole := &entities.Role{
		ID:   types.MSSQLUUID{},
		Name: "customer",
	}

	mockUserRepo.On("GetByEmail", ctx, "test@example.com").Return(nil, errors.New("not found"))
	mockAuthService.On("HashPassword", "password123").Return("hashed_password", nil)
	mockUserRepo.On("Create", ctx, mock.AnythingOfType("*entities.User")).Return(nil)
	mockRoleRepo.On("GetByName", ctx, "customer").Return(customerRole, nil)
	mockRoleRepo.On("AssignRoleToUser", ctx, mock.AnythingOfType("types.MSSQLUUID"), customerRole.ID, mock.AnythingOfType("types.MSSQLUUID")).Return(nil)

	err := usecase.Register(ctx, user)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockAuthService.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
}

func TestUserUsecase_Register_UserAlreadyExists(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockAuthService := new(MockAuthService)
	usecase := usecases.NewUserUsecase(mockUserRepo, mockRoleRepo, mockAuthService)
	ctx := context.Background()

	existingUser := &entities.User{Email: "test@example.com"}
	user := &entities.User{
		Email:    "test@example.com",
		Password: "password123",
	}

	mockUserRepo.On("GetByEmail", ctx, "test@example.com").Return(existingUser, nil)

	err := usecase.Register(ctx, user)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
	mockUserRepo.AssertExpectations(t)
}

func TestUserUsecase_Register_HashPasswordError(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockAuthService := new(MockAuthService)
	usecase := usecases.NewUserUsecase(mockUserRepo, mockRoleRepo, mockAuthService)
	ctx := context.Background()

	user := &entities.User{
		Email:    "test@example.com",
		Password: "password123",
	}

	mockUserRepo.On("GetByEmail", ctx, "test@example.com").Return(nil, nil)
	mockAuthService.On("HashPassword", "password123").Return("", errors.New("hash error"))

	err := usecase.Register(ctx, user)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to hash password")
	mockUserRepo.AssertExpectations(t)
	mockAuthService.AssertExpectations(t)
}

func TestUserUsecase_Login_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockAuthService := new(MockAuthService)
	usecase := usecases.NewUserUsecase(mockUserRepo, mockRoleRepo, mockAuthService)
	ctx := context.Background()

	user := &entities.User{
		ID:       types.MSSQLUUID{},
		Email:    "test@example.com",
		Password: "hashed_password",
		Roles: []entities.Role{
			{Name: "customer"},
		},
	}

	mockUserRepo.On("GetByEmail", ctx, "test@example.com").Return(user, nil)
	mockAuthService.On("ComparePassword", "hashed_password", "password123").Return(nil)
	mockAuthService.On("GenerateToken", user.ID, "customer").Return("jwt_token", nil)

	resultUser, token, err := usecase.Login(ctx, "test@example.com", "password123")

	assert.NoError(t, err)
	assert.NotNil(t, resultUser)
	assert.Equal(t, "jwt_token", token)
	mockUserRepo.AssertExpectations(t)
	mockAuthService.AssertExpectations(t)
}

func TestUserUsecase_Login_InvalidEmail(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockAuthService := new(MockAuthService)
	usecase := usecases.NewUserUsecase(mockUserRepo, mockRoleRepo, mockAuthService)
	ctx := context.Background()

	mockUserRepo.On("GetByEmail", ctx, "test@example.com").Return(nil, errors.New("not found"))

	resultUser, token, err := usecase.Login(ctx, "test@example.com", "password123")

	assert.Error(t, err)
	assert.Nil(t, resultUser)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "invalid credentials")
	mockUserRepo.AssertExpectations(t)
}

func TestUserUsecase_Login_InvalidPassword(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockAuthService := new(MockAuthService)
	usecase := usecases.NewUserUsecase(mockUserRepo, mockRoleRepo, mockAuthService)
	ctx := context.Background()

	user := &entities.User{
		Email:    "test@example.com",
		Password: "hashed_password",
	}

	mockUserRepo.On("GetByEmail", ctx, "test@example.com").Return(user, nil)
	mockAuthService.On("ComparePassword", "hashed_password", "wrong_password").Return(errors.New("password mismatch"))

	resultUser, token, err := usecase.Login(ctx, "test@example.com", "wrong_password")

	assert.Error(t, err)
	assert.Nil(t, resultUser)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "invalid credentials")
	mockUserRepo.AssertExpectations(t)
	mockAuthService.AssertExpectations(t)
}

func TestUserUsecase_GetUserByID_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockAuthService := new(MockAuthService)
	usecase := usecases.NewUserUsecase(mockUserRepo, mockRoleRepo, mockAuthService)
	ctx := context.Background()

	userID := types.MSSQLUUID{}
	user := &entities.User{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Test User",
	}

	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil)

	result, err := usecase.GetUserByID(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test@example.com", result.Email)
	mockUserRepo.AssertExpectations(t)
}

func TestUserUsecase_UpdateUser_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockAuthService := new(MockAuthService)
	usecase := usecases.NewUserUsecase(mockUserRepo, mockRoleRepo, mockAuthService)
	ctx := context.Background()

	userID := types.MSSQLUUID{}
	existingUser := &entities.User{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Old Name",
		Phone: "111",
	}

	updateData := &entities.User{
		Name:  "New Name",
		Phone: "222",
	}

	mockUserRepo.On("GetByID", ctx, userID).Return(existingUser, nil)
	mockUserRepo.On("Update", ctx, mock.AnythingOfType("*entities.User")).Return(nil)

	result, err := usecase.UpdateUser(ctx, userID, updateData)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "New Name", result.Name)
	assert.Equal(t, "222", result.Phone)
	mockUserRepo.AssertExpectations(t)
}

func TestUserUsecase_UpdateUser_NotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockAuthService := new(MockAuthService)
	usecase := usecases.NewUserUsecase(mockUserRepo, mockRoleRepo, mockAuthService)
	ctx := context.Background()

	userID := types.MSSQLUUID{}
	updateData := &entities.User{Name: "New Name"}

	mockUserRepo.On("GetByID", ctx, userID).Return(nil, nil)

	result, err := usecase.UpdateUser(ctx, userID, updateData)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "user not found")
	mockUserRepo.AssertExpectations(t)
}

func TestUserUsecase_DeleteUser_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockAuthService := new(MockAuthService)
	usecase := usecases.NewUserUsecase(mockUserRepo, mockRoleRepo, mockAuthService)
	ctx := context.Background()

	userID := types.MSSQLUUID{}
	user := &entities.User{ID: userID}

	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil)
	mockUserRepo.On("Delete", ctx, userID).Return(nil)

	err := usecase.DeleteUser(ctx, userID)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

func TestUserUsecase_DeleteUser_NotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockAuthService := new(MockAuthService)
	usecase := usecases.NewUserUsecase(mockUserRepo, mockRoleRepo, mockAuthService)
	ctx := context.Background()

	userID := types.MSSQLUUID{}

	mockUserRepo.On("GetByID", ctx, userID).Return(nil, nil)

	err := usecase.DeleteUser(ctx, userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
	mockUserRepo.AssertExpectations(t)
}

func TestUserUsecase_GetUsers_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockAuthService := new(MockAuthService)
	usecase := usecases.NewUserUsecase(mockUserRepo, mockRoleRepo, mockAuthService)
	ctx := context.Background()

	users := []*entities.User{
		{Email: "user1@example.com"},
		{Email: "user2@example.com"},
	}

	mockUserRepo.On("GetAll", ctx, 10, 0).Return(users, nil)

	result, err := usecase.GetUsers(ctx, 10, 0)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockUserRepo.AssertExpectations(t)
}

func TestUserUsecase_GetUsersPaginated_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockAuthService := new(MockAuthService)
	usecase := usecases.NewUserUsecase(mockUserRepo, mockRoleRepo, mockAuthService)
	ctx := context.Background()

	users := []*entities.User{
		{Email: "user1@example.com"},
	}

	pagParams := pagination.Params{Page: 1, PageSize: 10}
	filterParams := pagination.FilterParams{}

	mockUserRepo.On("GetAllPaginated", ctx, pagParams, filterParams).Return(users, 1, nil)

	result, total, err := usecase.GetUsersPaginated(ctx, pagParams, filterParams)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), total)
	mockUserRepo.AssertExpectations(t)
}

func TestUserUsecase_RefreshToken_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockAuthService := new(MockAuthService)
	usecase := usecases.NewUserUsecase(mockUserRepo, mockRoleRepo, mockAuthService)
	ctx := context.Background()

	userID := types.MSSQLUUID{}
	refreshToken := "old_token"

	mockAuthService.On("ValidateToken", refreshToken).Return(userID, "customer", nil)
	mockAuthService.On("GenerateToken", userID, "customer").Return("new_token", nil)

	newToken, err := usecase.RefreshToken(ctx, refreshToken)

	assert.NoError(t, err)
	assert.Equal(t, "new_token", newToken)
	mockAuthService.AssertExpectations(t)
}

func TestUserUsecase_RefreshToken_InvalidToken(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockAuthService := new(MockAuthService)
	usecase := usecases.NewUserUsecase(mockUserRepo, mockRoleRepo, mockAuthService)
	ctx := context.Background()

	refreshToken := "invalid_token"

	mockAuthService.On("ValidateToken", refreshToken).Return(nil, "", errors.New("invalid token"))

	newToken, err := usecase.RefreshToken(ctx, refreshToken)

	assert.Error(t, err)
	assert.Empty(t, newToken)
	assert.Contains(t, err.Error(), "invalid refresh token")
	mockAuthService.AssertExpectations(t)
}

func TestUserUsecase_ListMechanics_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockAuthService := new(MockAuthService)
	usecase := usecases.NewUserUsecase(mockUserRepo, mockRoleRepo, mockAuthService)
	ctx := context.Background()

	mechanics := []*entities.User{
		{Email: "mechanic1@example.com"},
		{Email: "mechanic2@example.com"},
	}

	mockUserRepo.On("GetByRole", ctx, "mechanic").Return(mechanics, nil)

	result, err := usecase.ListMechanics(ctx)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockUserRepo.AssertExpectations(t)
}
