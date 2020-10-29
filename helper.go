package smpp

import (
	"strings"
	"time"
)

// AllInAlphabet checks if all characters of s are in alphabet
func AllInAlphabet(s string, alphabet string) bool {
	for _, r := range s {
		if !strings.ContainsRune(alphabet, r) {
			return false
		}
	}
	return true
}

// ContainsUnicodeChar detects whether text contains unicode character or not.
func ContainsUnicodeChar(text string) bool {
	for _, char := range text {
		if char > 127 {
			return true
		}
	}
	return false
}

// RuneLength returns the number of runes (characters) of a string
func RuneLength(text string) int {
	return len([]rune(text))
}

// WaitSeconds sleeps for some specific seconds
func WaitSeconds(seconds int) {
	time.Sleep(time.Second * time.Duration(seconds))
}
