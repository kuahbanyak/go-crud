package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/domain/services"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
	"golang.org/x/crypto/bcrypt"
)

type jwtService struct {
	secretKey  string
	expiration time.Duration
}

func NewJWTService(secretKey string, expirationHours int) services.AuthService {
	return &jwtService{
		secretKey:  secretKey,
		expiration: time.Duration(expirationHours) * time.Hour,
	}
}

type Claims struct {
	UserID types.MSSQLUUID `json:"user_id"`
	Role   entities.Role   `json:"role"`
	jwt.RegisteredClaims
}

func (j *jwtService) GenerateToken(userID types.MSSQLUUID, role entities.Role) (string, error) {
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "go-crud-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *jwtService) ValidateToken(tokenString string) (types.MSSQLUUID, entities.Role, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return types.MSSQLUUID{}, "", err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, claims.Role, nil
	}

	return types.MSSQLUUID{}, "", errors.New("invalid token")
}

func (j *jwtService) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (j *jwtService) ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
