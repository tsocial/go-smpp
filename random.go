package smpp

import (
	"crypto/rand"
)

// RandomString generates a random string consisting of characters in the provided alphabet
func RandomString(n int, alphabet string) string {
	if n <= 0 {
		return ""
	}

	r := []rune(alphabet)
	k := byte(len(r))

	s := make([]rune, n)
	b := make([]byte, n)

	_, _ = rand.Read(b)
	for i := range b {
		s[i] = r[b[i]%k]
	}

	return string(s)
}
