package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/domain/repositories"
	"github.com/kuahbanyak/go-crud/internal/domain/services"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
	"github.com/kuahbanyak/go-crud/pkg/pagination"
)

type UserUsecase struct {
	userRepo    repositories.UserRepository
	authService services.AuthService
}

func NewUserUsecase(userRepo repositories.UserRepository, authService services.AuthService) *UserUsecase {
	return &UserUsecase{
		userRepo:    userRepo,
		authService: authService,
	}
}
func (u *UserUsecase) Register(ctx context.Context, user *entities.User) error {
	existingUser, _ := u.userRepo.GetByEmail(ctx, user.Email)
	if existingUser != nil {
		return errors.New("user with this email already exists")
	}
	hashedPassword, err := u.authService.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = hashedPassword
	return u.userRepo.Create(ctx, user)
}
func (u *UserUsecase) Login(ctx context.Context, email, password string) (*entities.User, string, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}
	if err := u.authService.ComparePassword(user.Password, password); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Get role name from Roles relationship (use first role if multiple, or empty string if none)
	roleName := ""
	if len(user.Roles) > 0 {
		roleName = user.Roles[0].Name
	}

	token, err := u.authService.GenerateToken(user.ID, roleName)
	if err != nil {
		return nil, "", errors.New("failed to generate token")
	}
	return user, token, nil
}
func (u *UserUsecase) GetUserByID(ctx context.Context, id types.MSSQLUUID) (*entities.User, error) {
	return u.userRepo.GetByID(ctx, id)
}
func (u *UserUsecase) UpdateUser(ctx context.Context, id types.MSSQLUUID, updateData *entities.User) (*entities.User, error) {
	existingUser, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		return nil, errors.New("user not found")
	}
	if updateData.Name != "" {
		existingUser.Name = updateData.Name
	}
	if updateData.Phone != "" {
		existingUser.Phone = updateData.Phone
	}
	err = u.userRepo.Update(ctx, existingUser)
	if err != nil {
		return nil, err
	}
	return existingUser, nil
}
func (u *UserUsecase) GetUsers(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	return u.userRepo.GetAll(ctx, limit, offset)
}

func (u *UserUsecase) GetUsersPaginated(ctx context.Context, pagParams pagination.Params, filterParams pagination.FilterParams) ([]*entities.User, int64, error) {
	return u.userRepo.GetAllPaginated(ctx, pagParams, filterParams)
}

func (u *UserUsecase) DeleteUser(ctx context.Context, id types.MSSQLUUID) error {
	existingUser, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not found")
	}
	return u.userRepo.Delete(ctx, id)
}
func (u *UserUsecase) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	userID, role, err := u.authService.ValidateToken(refreshToken)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}
	return u.authService.GenerateToken(userID, role)
}
func (u *UserUsecase) ListMechanics(ctx context.Context) ([]*entities.User, error) {
	return u.userRepo.GetByRole(ctx, "mechanic")
}
