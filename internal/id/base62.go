package id

import (
	"errors"
	"strings"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const base = uint64(len(alphabet))

var (
	index = func() map[rune]uint64 {
		m := make(map[rune]uint64, len(alphabet))
		for i, r := range alphabet {
			m[r] = uint64(i)
		}
		return m
	}()
	ErrInvalidChar = errors.New("invalid base62 character")
)

func Encode(u uint64) string {
	if u == 0 {
		return string(alphabet[0])
	}
	var b []byte
	for u > 0 {
		r := u % base
		b = append(b, alphabet[r])
		u /= base
	}
	// reverse
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

func Decode(s string) (uint64, error) {
	var n uint64
	for _, r := range s {
		v, ok := index[r]
		if !ok {
			return 0, ErrInvalidChar
		}
		n = n*base + v
	}
	return n, nil
}

// Pad ensures a minimum length by left-padding with the first alphabet char (zero).
func Pad(s string, min int) string {
	if len(s) >= min {
		return s
	}
	return strings.Repeat(string(alphabet[0]), min-len(s)) + s
}
