package helpers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppendJSON_Quoted(t *testing.T) {
	tests := []struct {
		description string
		input       string
		expected    string
	}{
		{
			description: "01 empty string becomes two quotes",
			input:       "",
			expected:    `""`,
		},
		{
			description: "02 control characters below 0x20 are escaped as \\u00XX and wrapped in quotes",
			input:       string([]byte{0x00, 0x01, 0x1F}),
			expected:    `"\u0000\u0001\u001f"`,
		},
		{
			description: "03 backslash is escaped and wrapped in quotes",
			input:       `\`,
			expected:    "\"\\\\\"",
		},
		{
			description: "04 double quote is escaped and wrapped in quotes",
			input:       `"`,
			expected:    "\"\\\"\"",
		},

		{
			description: "05 newline is escaped and wrapped in quotes",
			input:       "\n",
			expected:    `"\n"`,
		},
		{
			description: "06 carriage return is escaped and wrapped in quotes",
			input:       "\r",
			expected:    `"\r"`,
		},
		{
			description: "07 tab is escaped and wrapped in quotes",
			input:       "\t",
			expected:    `"\t"`,
		},
		{
			description: "08 backspace is escaped and wrapped in quotes",
			input:       "\b",
			expected:    `"\b"`,
		},
		{
			description: "09 formfeed is escaped and wrapped in quotes",
			input:       "\f",
			expected:    `"\f"`,
		},
		{
			description: "10 normal ASCII characters are appended unchanged and wrapped in quotes",
			input:       "abcXYZ123",
			expected:    `"abcXYZ123"`,
		},
		{
			description: "11 mixed content with all escape types and surrounding quotes",
			input:       "A\nB\tC\"D\\E\fF\bG\rH\x07I",
			expected:    `"A\nB\tC\"D\\E\fF\bG\rH\u0007I"`,
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.description,
			func(t *testing.T) {
				var buf []byte

				buf = AppendJSON_Quoted(buf, []byte(tc.input))

				require.Equal(t, tc.expected, string(buf))
			},
		)
	}
}
