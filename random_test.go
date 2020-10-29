package smpp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomString(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name     string
		Alphabet string
	}{
		{"upper_english", UpperEnglishLetters},
		{"lower_english", LowerEnglishLetters},
		{"decimal_digit", DecimalDigits},
		{"special", SpecialCharacters},
		{"english", EnglishLetters},
		{"alpha_numeric", AlphaNumericCharacters},
		{"keyboard", KeyboardCharacters},
		{"unicode_letters", UnicodeLetters},
	}
	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			for n := 0; n < 20; n++ {
				s := RandomString(n, tc.Alphabet)
				r := []rune(s)
				require.Len(t, r, n)
				require.True(t, AllInAlphabet(s, tc.Alphabet))
			}
		})
	}
}
