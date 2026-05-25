package helpers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetEstimatedFormatless(t *testing.T) {
	tests := []struct {
		description string
		args        []any
		expected    uint32
	}{
		{
			description: "1. single long string argument",
			args:        []any{"hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh"},
			expected:    33, // Only the string, no separators needed
		},
		{
			description: "2. multiple string arguments with separator spaces",
			args:        []any{"abc", "def"},
			expected:    7, // len("abc") + len("def") + 1 separator space = 7
		},
		{
			description: "3. mixed string and float64",
			args:        []any{"abc", 2.1},
			expected:    7, // len("abc") + len("2.1") + 1 separator space = 7
		},
		{
			description: "4. empty args slice",
			args:        []any{},
			expected:    0,
		},
		{
			description: "5. single negative integer",
			args:        []any{-420},
			expected:    4, // len("-420") = 4
		},
		{
			description: "6. booleans",
			args:        []any{true, false},
			expected:    10, // len("true") + len("false") + 1 separator space = 10
		},
		{
			description: "7. handling nil value",
			args:        []any{nil},
			expected:    4, // len("<nil>") = 4
		},
		{
			description: "8. complex type fallback",
			args:        []any{struct{ ID int }{ID: 9}},
			expected:    3, // fmt.Sprint output is "{9}", len = 3
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.description,
			func(t *testing.T) {
				estimated := GetEstimatedArgsSize(tc.args)
				if estimated != tc.expected {
					t.Errorf(
						"expected: %d versus estimated: %d for args %v",
						tc.expected,
						estimated,
						tc.args,
					)
				}
			},
		)
	}
}

func TestGetEstimatedMessageSize(t *testing.T) {
	tests := []struct {
		description string
		format      string
		args        []any
		expected    uint32
	}{
		{
			description: "01a plain text counts characters",
			format:      "hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh",
			args:        nil,
			expected:    33,
		},
		{
			description: "02 percent escape sequence %% produces single percent",
			format:      "100%% sure",
			args:        nil,
			expected:    9, // "100% sure"
		},
		{
			description: "03 trailing percent is emitted literally",
			format:      "end%",
			args:        nil,
			expected:    4, // 'e','n','d' + lone '%'
		},
		{
			description: "04 missing argument for verb emits percent+verb literally",
			format:      "%s %d",
			args:        []any{"only"},
			expected:    7, // len("only") + 1(space) + 2(for "%d")
		},
		{
			description: "05 %s accepts []byte",
			format:      "%s",
			args:        []any{[]byte("bytes")},
			expected:    5,
		},
		{
			description: "06 %d handles negative int (includes sign)",
			format:      "%d",
			args:        []any{int(-42)},
			expected:    3, // "-42" -> 3 chars
		},
		{
			description: "07 %d handles uint64 zero",
			format:      "%d",
			args:        []any{uint64(0)},
			expected:    1, // "0"
		},
		{
			description: "08 %t",
			format:      "%t",
			args:        []any{true},
			expected:    4,
		},
		{
			description: "09 %f uses float64Len conservative bound",
			format:      "%f",
			args:        []any{123.456},
			expected:    15, // integer part 123 -> 3 digits + 6 = 9
		},
		{
			description: "10 %v uses fmt.Sprint length for exotic types",
			format:      "X:%v",
			args:        []any{struct{ A int }{A: 5}},
			expected:    5, // "X:" (2) + fmt.Sprint(struct{A int}{5}) -> "{5}" (3) => 5
		},
		{
			description: "11 unsupported verb falls back to literal %x size 2",
			format:      "%q",
			args:        []any{"ignored"},
			expected:    2, // "%q" counted as two literal chars
		},
		{
			description: "12 multiple mixed verbs and literals",
			format:      "A:%s B:%d C:%f D:%t Z:%%",
			args:        []any{"s", 7, float32(1.5), false},
			// A: (2) + len("s")=1 -> 3
			// space between tokens are part of format
			// " B:" (3) + digits of 7 = 1 -> 4
			// " C:" (3) + float32Len(1.5) -> integer part 1 -> digits 1 + 6 = 7 -> 10
			// " D:" (3) + t worst-case 5 -> 8 -> 18
			// " Z:" (3) + %% -> single '%' (1) -> 4 -> total 22
			expected: 35,
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.description,
			func(t *testing.T) {
				estimated := GetEstimatedMessageSize(tc.format, tc.args)

				require.Equal(t,
					tc.expected,
					estimated,

					"expected: %d versus estimated: %d",
					tc.expected,
					estimated,
				)
			},
		)
	}
}
