package services

import (
	"github.com/kuahbanyak/go-crud/internal/shared/types"
)

type AuthService interface {
	GenerateToken(userID types.MSSQLUUID, role string) (string, error)
	ValidateToken(token string) (types.MSSQLUUID, string, error)
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) error
}
