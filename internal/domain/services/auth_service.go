package services
import (
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
)
type AuthService interface {
	GenerateToken(userID types.MSSQLUUID, role entities.Role) (string, error)
	ValidateToken(token string) (types.MSSQLUUID, entities.Role, error)
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) error
}

