package utils

import (
	"fmt"

	"github.com/google/uuid"
)

type IDGenerator interface {
	GenerateID() string
	GenerateIDWithPrefixes(prefix1 string, prefix2 string) string
}

type uuidGenerator struct{}

// Create instance as pointer reference so it can be used as singleton.
// Replaced for an interface of "any" other module when needed without breaking modularity principles (direction of dependency).
func NewUUIDGenerator() IDGenerator {
	return &uuidGenerator{}
}

func (g *uuidGenerator) GenerateID() string {
	return uuid.NewString()
}

func (g *uuidGenerator) GenerateIDWithPrefixes(prefix1 string, prefix2 string) string {
	return fmt.Sprintf("%s-%s-%s", prefix1, prefix2, g.GenerateID())
}
