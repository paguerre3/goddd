package utils

import (
	"fmt"

	"github.com/google/uuid"
)

type UUIDGenerator struct{}

// Create instance as pointer reference so it can be used as singleton.
// Replaced for an interface of "any" other module when needed without breaking modularity principles (direction of dependency).
func NewUUIDGenerator() *UUIDGenerator {
	return &UUIDGenerator{}
}

func (g *UUIDGenerator) GenerateID() string {
	return uuid.New().String()
}

func (g *UUIDGenerator) GenerateIDWithPrefixes(prefix1 string, prefix2 string) string {
	return fmt.Sprintf("%s&%s-%s", prefix1, prefix2, g.GenerateID())
}
