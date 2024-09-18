package utils

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGenerateID_ValidUUID checks if generateID returns a valid UUID
func TestGenerateID_ValidUUID(t *testing.T) {
	id := GenerateID()
	match, _ := regexp.MatchString(`^[a-f0-9\-]{36}$`, id)
	assert.True(t, match, "generateID should return a valid UUID")
}

// TestGenerateIDWithPrefix_ValidFormat checks if GenerateIDWithPrefix returns the correct format
func TestGenerateIDWithPrefix_ValidFormat(t *testing.T) {
	prefix := "prefix"
	id := GenerateIDWithPrefix(prefix)
	match, _ := regexp.MatchString(`^prefix-[a-f0-9\-]{36}$`, id)
	assert.True(t, match, "GenerateIDWithPrefix should return a valid prefixed UUID")
}

// TestGenerateIDWithPrefixes_ValidFormat checks if GenerateIDWithPrefixes returns the correct format
func TestGenerateIDWithPrefixes_ValidFormat(t *testing.T) {
	prefix1 := "prefix1"
	prefix2 := "prefix2"
	id := GenerateIDWithPrefixes(prefix1, prefix2)
	match, _ := regexp.MatchString(`^prefix1&prefix2-[a-f0-9\-]{36}$`, id)
	assert.True(t, match, "GenerateIDWithPrefixes should return a valid double-prefixed UUID")
}

// TestGenerateIDWithPrefixes_SamePrefixes checks if GenerateIDWithPrefixes returns different UUIDs when the prefixes are the same
func TestGenerateIDWithPrefixes_SamePrefixes(t *testing.T) {
	prefix := "prefix"
	id1 := GenerateIDWithPrefix(prefix)
	id2 := GenerateIDWithPrefix(prefix)
	assert.NotEqual(t, id1, id2, "GenerateIDWithPrefixes should return different UUIDs when the prefixes are the same")
}
