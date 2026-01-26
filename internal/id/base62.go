package id

import (
	"strings"
)

// Base62Encode is an alias for Encode
func Base62Encode(num int64) string {
	return Encode(num)
}

// Base62Decode is an alias for Decode
func Base62Decode(s string) (int64, error) {
	return Decode(s)
}

// PadLeft pads a string on the left with a specific character
func PadLeft(s string, length int, padChar byte) string {
	for len(s) < length {
		s = string(padChar) + s
	}
	return s
}

// PadRight pads a string on the right with a specific character
func PadRight(s string, length int, padChar byte) string {
	for len(s) < length {
		s = s + string(padChar)
	}
	return s
}

// Contains checks if a character is in the alphabet
func Contains(r rune) bool {
	return strings.ContainsRune(alphabet, r)
}
