package services

import (
	"github.com/google/uuid"
)

type AuthService interface {
	GenerateToken(userID uuid.UUID, role string) (string, error)
	ValidateToken(token string) (uuid.UUID, string, error)
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) error
}
