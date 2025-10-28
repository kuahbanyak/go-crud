package utils
import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)
type HashService struct {
	cost int
}
func NewHashService() *HashService {
	return &HashService{
		cost: bcrypt.DefaultCost,
	}
}
func (h *HashService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	return string(bytes), err
}
func (h *HashService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func (h *HashService) GenerateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
func (h *HashService) HashSHA256(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
func (h *HashService) GenerateAPIKey() (string, error) {
	token, err := h.GenerateRandomToken(32)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("sk_%s", token), nil
}

