package types

import (
	"github.com/google/uuid"
)

// NewUUID generates a new random UUID
func NewUUID() uuid.UUID {
	return uuid.New()
}

// ParseUUID parses a string into a UUID
func ParseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
