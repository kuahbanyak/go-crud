package utils
import (
	"regexp"
	"strings"
)
type Validator struct{}
func NewValidator() *Validator {
	return &Validator{}
}
func (v *Validator) IsValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}
func (v *Validator) IsValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	return hasUpper && hasLower && hasDigit
}
func (v *Validator) IsValidPhoneNumber(phone string) bool {
	digits := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")
	return len(digits) >= 10 && len(digits) <= 15
}
func (v *Validator) IsNotEmpty(value string) bool {
	return strings.TrimSpace(value) != ""
}
func (v *Validator) IsValidLength(value string, min, max int) bool {
	length := len(strings.TrimSpace(value))
	return length >= min && length <= max
}

