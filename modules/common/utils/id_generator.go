package utils

import (
	"fmt"

	"github.com/google/uuid"
)

func GenerateID() string {
	return uuid.New().String()
}

func GenerateIDWithPrefix(prefix string) string {
	return fmt.Sprintf("%s-%s", prefix, GenerateID())
}

func GenerateIDWithPrefixes(prefix1 string, prefix2 string) string {
	return fmt.Sprintf("%s&%s-%s", prefix1, prefix2, GenerateID())
}
