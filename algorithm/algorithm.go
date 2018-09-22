package algorithm

import (
	"fmt"
	"strings"
)

// All characters
const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	base   = int64(len(alphabet))
)

// Encode number to base62
func Encode(n int64) string {
	if n == 0 {
		return string(alphabet[0])
	}

	s := ""
	for ; n > 0; n = n / base {
		s = string(alphabet[n%base]) + s
	}
	return s
}

// Decode from base62 to int.
func Decode(key string) (int64, error) {
	var n int64
	for _, c := range []byte(key) {
		i := strings.IndexByte(alphabet, c)
		if i < 0 {
			return 0, fmt.Errorf("unexpected character %c in base62 literal", c)
		}
		n = base*n + int64(i)
	}
	return n, nil
}
