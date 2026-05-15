package helpers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func floatEqualStr(t *testing.T, got string, want float64, relTol float64) {
	t.Helper()

	parsed, err := strconv.ParseFloat(got, 64)
	require.NoError(t, err, "failed to parse float token %q", got)

	// relative tolerance check
	diff := parsed - want
	if diff < 0 {
		diff = -diff
	}

	// allow small absolute tolerance for values near zero
	abs := parsed
	if abs < 0 {
		abs = -abs
	}

	if abs < 1 {
		abs = 1
	}

	require.LessOrEqual(t,
		diff/abs,
		relTol,

		"float %q not within relative tolerance of %v",
		got,
		want,
	)
}

func TestAppendArgs_RobustFloatHandling(t *testing.T) {
	tests := []struct {
		description string
		args        []any
		// For floats, provide the numeric expected value and tolerance.
		// For non-floats, provide the exact expected token string.
		expectedTokens []any // token per arg: string for exact match, float64 for numeric compare
	}{
		{
			description:    "01 no args returns empty",
			args:           []any{},
			expectedTokens: []any{},
		},
		{
			description:    "02 single string",
			args:           []any{"hello"},
			expectedTokens: []any{"hello"},
		},
		{
			description:    "03 multiple strings separated by space",
			args:           []any{"hello", "world"},
			expectedTokens: []any{"hello", "world"},
		},
		{
			description:    "04 []byte is appended as bytes",
			args:           []any{[]byte("bytez")},
			expectedTokens: []any{"bytez"},
		},
		{
			description:    "05 ints and uints",
			args:           []any{int(42), int64(-7), int32(123), uint(77), uint64(9999999999)},
			expectedTokens: []any{"42", "-7", "123", "77", "9999999999"},
		},
		{
			description:    "06 floats float64 and float32",
			args:           []any{float64(3.14159265358979), float32(1.23456)},
			expectedTokens: []any{3.14159265358979, 1.23456}, // numeric expectations
		},
		{
			description:    "07 bool true and false",
			args:           []any{true, false},
			expectedTokens: []any{"true", "false"},
		},
		{
			description:    "08 error becomes its Error string",
			args:           []any{errors.New("boom")},
			expectedTokens: []any{"boom"},
		},
		{
			description:    "09 nil becomes null",
			args:           []any{nil},
			expectedTokens: []any{"null"},
		},
		{
			description: "10 mixed types with spacing",
			args: []any{
				"start",
				[]byte("bytes"),
				123,
				float64(2.5),
				true,
				errors.New("err"),
				nil,
				struct{ A int }{A: 5},
			},
			expectedTokens: []any{
				"start",
				"bytes",
				"123",
				2.5,
				"true",
				"err",
				"null",
				fmt.Sprint(struct{ A int }{A: 5}),
			},
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.description,
			func(t *testing.T) {
				var dst []byte

				dst = AppendArgs(dst, tc.args...)
				out := string(dst)

				// split tokens by single space (AppendArgs uses single space between args)
				var tokens []string
				if out == "" {
					tokens = []string{}
				} else {
					tokens = strings.Split(out, " ")
				}

				require.Equal(t, len(tc.expectedTokens), len(tokens), "token count mismatch: got %q", out)

				for i, want := range tc.expectedTokens {
					got := tokens[i]

					switch w := want.(type) {
					case string:
						require.Equal(t, w, got, "token %d mismatch", i)
					case float64:
						// use a small relative tolerance; 1e-12 is strict but safe for typical formatting differences
						floatEqualStr(t, got, w, 1e-9)
					default:
						// fallback: compare stringified expected
						require.Equal(t, fmt.Sprint(w), got, "token %d mismatch (fallback)", i)
					}
				}
			},
		)
	}
}

func TestAppendf(t *testing.T) {
	tests := []struct { //nolint:govet
		description string
		format      string
		args        []any
		expected    string
	}{
		{
			description: "01 plain text without verbs",
			format:      "hello world",
			args:        nil,
			expected:    "hello world",
		},
		{
			description: "02 percent escape sequence %% produces single percent",
			format:      "100%% sure",
			args:        nil,
			expected:    "100% sure",
		},
		{
			description: "03 trailing percent is emitted literally",
			format:      "end%",
			args:        nil,
			expected:    "end%",
		},
		{
			description: "04 missing argument for verb emits percent+verb literally",
			format:      "%s %d",
			args:        []any{"only"},
			expected:    "only %d",
		},
		{
			description: "05 %s accepts []byte",
			format:      "%s",
			args:        []any{[]byte("bytes")},
			expected:    "bytes",
		},
		{
			description: "06 %s falls back to fmt.Sprint for exotic types",
			format:      "%s",
			args:        []any{123},
			expected:    "123",
		},
		{
			description: "07 %d handles signed and unsigned integers",
			format:      "%d %d",
			args:        []any{int(42), uint64(7)},
			expected:    "42 7",
		},
		{
			description: "08 %v handles many types including float, bool, error, nil, string",
			format:      "%v %v %v %v %v",
			args:        []any{float64(1.23), true, errors.New("boom"), nil, "x"},
			expected:    "1.23 true boom null x",
		},
		{
			description: "09 %t prints bool and falls back for non-bool",
			format:      "%t %t",
			args:        []any{true, 5},
			expected:    "true 5",
		},
		{
			description: "10 %f prints floats and falls back for non-float",
			format:      "f1:%f f2:%f f3:%f",
			args:        []any{float64(2.5), float32(3.5), "x"},
			expected:    "f1:2.5 f2:3.5 f3:x",
		},
		{
			description: "11 unsupported verb is emitted literally and does not consume arg",
			format:      "%q",
			args:        []any{"a"},
			expected:    "%q",
		},
		{
			description: "12 %% does not consume following argument",
			format:      "%% %s",
			args:        []any{"ok"},
			expected:    "% ok",
		},
		{
			description: "13 interleaved verbs and literals",
			format:      "A:%s B:%d C:%v D:%t E:%f Z:%%",
			args:        []any{"s", 7, "v", false, float32(1.5)},
			expected:    "A:s B:7 C:v D:false E:1.5 Z:%",
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.description,
			func(t *testing.T) {
				var dst []byte

				dst = Appendf(dst, tc.format, tc.args)
				require.Equal(t, tc.expected, string(dst))
			},
		)
	}
}
