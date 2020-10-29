package smpp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAllInAlphabet(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Text     string
		Alphabet string
		Expected bool
	}{
		{RandomString(10, KeyboardCharacters), KeyboardCharacters, true},
		{RandomString(13, SpecialCharacters), AlphaNumericCharacters, false},
	}
	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.Text, func(t *testing.T) {
			t.Parallel()
			result := AllInAlphabet(tc.Text, tc.Alphabet)
			require.Equal(t, tc.Expected, result)
		})
	}
}

func TestContainsUnicodeChar(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Value     string
		IsUnicode bool
	}{
		{"ascii_string", "this is an ascii string", false},
		{"unicode_string", "đây là một chuỗi unicode", true},
	}

	for _, testcase := range testCases {
		tc := testcase
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			isUnicode := ContainsUnicodeChar(tc.Value)
			require.Equal(t, tc.IsUnicode, isUnicode)
		})
	}
}

func TestRuneLength(t *testing.T) {
	t.Parallel()

	text1 := "xin chào các bạn"
	text2 := "hello friends"

	require.Equal(t, 16, RuneLength(text1))
	require.Equal(t, 13, RuneLength(text2))
}
