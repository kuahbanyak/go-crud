package services

import (
	"context"

	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
)

type UserService interface {
	Register(ctx context.Context, user *entities.User) error
	Login(ctx context.Context, email, password string) (*entities.User, string, error)
	GetProfile(ctx context.Context, userID types.MSSQLUUID) (*entities.User, error)
	UpdateProfile(ctx context.Context, user *entities.User) error
	ChangePassword(ctx context.Context, userID types.MSSQLUUID, oldPassword, newPassword string) error
	DeleteAccount(ctx context.Context, userID types.MSSQLUUID) error
	ListUsers(ctx context.Context, role entities.Role, limit, offset int) ([]*entities.User, error)
}

type AuthService interface {
	GenerateToken(userID types.MSSQLUUID, role entities.Role) (string, error)
	ValidateToken(token string) (types.MSSQLUUID, entities.Role, error)
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) error
}
