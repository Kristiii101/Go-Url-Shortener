package id

import (
	"fmt"
	"math"
	"strings"
)

// Base62Encode encodes a number to base62
func Base62Encode(num int64) string {
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

// Base62Decode decodes a base62 string to number
func Base62Decode(s string) (int64, error) {
	var result int64

	for i, r := range s {
		pos := strings.IndexRune(alphabet, r)
		if pos == -1 {
			return 0, fmt.Errorf("invalid character: %c", r)
		}

		result += int64(pos) * int64(math.Pow(float64(base), float64(len(s)-i-1)))
	}

	return result, nil
}
