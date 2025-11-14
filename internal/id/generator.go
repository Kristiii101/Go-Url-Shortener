package id

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

const (
	// Base62 alphabet for URL-safe short codes
	alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	base     = int64(len(alphabet))
)

// Encode encodes a number to base62
func Encode(num int64) string {
	if num == 0 {
		return string(alphabet[0])
	}

	var result []byte
	for num > 0 {
		result = append(result, alphabet[num%base])
		num /= base
	}

	// Reverse
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}

// Pad pads a string with leading zeros to reach the desired length
func Pad(s string, length int) string {
	for len(s) < length {
		s = string(alphabet[0]) + s
	}
	return s
}

// Decode decodes a base62 string to number
func Decode(s string) (int64, error) {
	var result int64
	multiplier := int64(1)

	for i := len(s) - 1; i >= 0; i-- {
		pos := strings.IndexRune(alphabet, rune(s[i]))
		if pos == -1 {
			return 0, fmt.Errorf("invalid character: %c", s[i])
		}

		result += int64(pos) * multiplier
		multiplier *= base
	}

	return result, nil
}

// Generator generates unique IDs
type Generator struct {
	minLen int
	maxLen int
}

// NewGenerator creates a new ID generator
func NewGenerator(minLen, maxLen int) *Generator {
	return &Generator{
		minLen: minLen,
		maxLen: maxLen,
	}
}

// GenerateKey generates a random key of specified length
func (g *Generator) GenerateKey(length int) (string, error) {
	if length < g.minLen || length > g.maxLen {
		return "", fmt.Errorf("length must be between %d and %d", g.minLen, g.maxLen)
	}

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(base))
		if err != nil {
			return "", fmt.Errorf("random generation failed: %w", err)
		}
		result[i] = alphabet[num.Int64()]
	}

	return string(result), nil
}

// GenerateFromID generates a deterministic key from a numeric ID
func (g *Generator) GenerateFromID(id int64) string {
	encoded := Encode(id)
	return Pad(encoded, g.minLen)
}

// IsValidKey checks if a key is valid
func IsValidKey(key string, minLen, maxLen int) bool {
	if len(key) < minLen || len(key) > maxLen {
		return false
	}

	for _, r := range key {
		if !strings.ContainsRune(alphabet+"_-", r) {
			return false
		}
	}

	return true
}
