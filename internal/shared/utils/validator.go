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
	// Password must be at least 8 characters long
	if len(password) < 8 {
		return false
	}

	// Check for at least one uppercase letter, one lowercase letter, and one digit
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)

	return hasUpper && hasLower && hasDigit
}

func (v *Validator) IsValidPhoneNumber(phone string) bool {
	// Remove all non-digit characters
	digits := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")

	// Phone number should be between 10-15 digits
	return len(digits) >= 10 && len(digits) <= 15
}

func (v *Validator) IsNotEmpty(value string) bool {
	return strings.TrimSpace(value) != ""
}

func (v *Validator) IsValidLength(value string, min, max int) bool {
	length := len(strings.TrimSpace(value))
	return length >= min && length <= max
}
