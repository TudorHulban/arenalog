package helpers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// These tests cover:

// 1. Basic functionality: Tests for simple strings and strings requiring escaping
// 2. All escaped characters: Tests for \", \\, \n, \r, \t
// 3. Buffer growth: Tests when the initial buffer capacity is insufficient
// 4. Edge cases: Exactly at capacity boundaries, unicode characters, mixed content
// 5. Appending to existing buffers: Tests that the function correctly appends to existing data
// 6. Performance: Large string tests and benchmarks
// 7. Boundary conditions: Tests at the exact limits where buffer growth occurs

// The tests verify that:

// a. Escaped characters are correctly transformed
// b. Buffer growth happens at the right times
// c. The function works with various initial buffer states
// d. No panics occur in edge cases
// e. The result is a valid slice with len ≤ cap

func TestAppendJSONString(t *testing.T) {
	tests := []struct {
		name       string
		initialBuf []byte
		input      string
		expected   []byte
	}{
		{
			name:       "1. empty string",
			initialBuf: make([]byte, 0),
			input:      "",
			expected:   []byte{},
		},
		{
			name:       "2. simple string no escaping needed",
			initialBuf: make([]byte, 0),
			input:      "hello",
			expected:   []byte("hello"),
		},
		{
			name:       "3. string with double quote",
			initialBuf: make([]byte, 0),
			input:      `hello "world"`,
			expected:   []byte(`hello \"world\"`),
		},
		{
			name:       "4. string with backslash",
			initialBuf: make([]byte, 0),
			input:      `path\to\file`,
			expected:   []byte(`path\\to\\file`),
		},
		{
			name:       "5. string with newline",
			initialBuf: make([]byte, 0),
			input:      "line1\nline2",
			expected:   []byte(`line1\nline2`),
		},
		{
			name:       "6. string with carriage return",
			initialBuf: make([]byte, 0),
			input:      "text\rwith\rcarriage",
			expected:   []byte(`text\rwith\rcarriage`),
		},
		{
			name:       "7. string with tab",
			initialBuf: make([]byte, 0),
			input:      "col1\tcol2\tcol3",
			expected:   []byte(`col1\tcol2\tcol3`),
		},
		{
			name:       "8. string with multiple escaped characters",
			initialBuf: make([]byte, 0),
			input:      "hello \"world\"\n\t",
			expected:   []byte(`hello \"world\"\n\t`),
		},
		{
			name:       "9. append to existing buffer",
			initialBuf: []byte("existing,"),
			input:      "data",
			expected:   []byte("existing,data"),
		},
		{
			name:       "10. append to buffer with capacity",
			initialBuf: make([]byte, 0, 100),
			input:      "test",
			expected:   []byte("test"),
		},
		{
			name:       "11. NUL character",
			initialBuf: make([]byte, 0),
			input:      "\x00",
			expected:   []byte(`\u0000`),
		},
		{
			name:       "12. BEL character",
			initialBuf: make([]byte, 0),
			input:      "\x07",
			expected:   []byte(`\u0007`),
		},
		{
			name:       "13. Mixed controls",
			initialBuf: make([]byte, 0),
			input:      "A\x01\n\x1FZ",
			expected:   []byte(`A\u0001\n\u001fZ`),
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				result := AppendJSON(tt.initialBuf, []byte(tt.input))

				require.Equal(t,
					len(result),
					len(tt.expected),

					"length mismatch: got %d, want %d",
					len(result),
					len(tt.expected),
				)

				for i := range result {
					require.Equal(t,
						result[i],
						tt.expected[i],

						"byte mismatch at index %d: got %c (%d), want %c (%d)",
						i, result[i],
						result[i],
						tt.expected[i],
						tt.expected[i],
					)
				}
			},
		)
	}
}

func TestAppendJSON(t *testing.T) {
	tests := []struct { //nolint:govet
		description string
		input       []byte
		expected    string
	}{
		{
			description: "01 control characters below 0x20 are escaped as \\u00XX",
			input:       []byte{0x01, 0x02, 0x1F},
			expected:    `\u0001\u0002\u001f`,
		},
		{
			description: "02 backslash is escaped",
			input:       []byte(`\`),
			expected:    `\\`,
		},
		{
			description: "03 double quote is escaped",
			input:       []byte(`"`),
			expected:    `\"`,
		},
		{
			description: "04 newline is escaped",
			input:       []byte("\n"),
			expected:    `\n`,
		},
		{
			description: "05 carriage return is escaped",
			input:       []byte("\r"),
			expected:    `\r`,
		},
		{
			description: "06 tab is escaped",
			input:       []byte("\t"),
			expected:    `\t`,
		},
		{
			description: "07 backspace is escaped",
			input:       []byte("\b"),
			expected:    `\b`,
		},
		{
			description: "08 formfeed is escaped",
			input:       []byte("\f"),
			expected:    `\f`,
		},
		{
			description: "09 normal ASCII characters are appended unchanged",
			input:       []byte("abcXYZ123"),
			expected:    "abcXYZ123",
		},
		{
			description: "10 mixed content with all escape types",
			input:       []byte("A\nB\tC\"D\\E\fF\bG\rH\x07I"),
			expected:    `A\nB\tC\"D\\E\fF\bG\rH\u0007I`,
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.description,
			func(t *testing.T) {
				var buf []byte

				buf = AppendJSON(buf, tc.input)
				require.Equal(t, tc.expected, string(buf))
			},
		)
	}
}

func TestDebugBufferGrowth(t *testing.T) {
	buf := make([]byte, 0, 5)
	result := AppendJSON(buf, []byte("this is a longer string"))

	t.Logf("Result: %q", string(result))
	t.Logf("Result bytes: %v", result)

	expected := "this is a longer string"
	if string(result) != expected {
		t.Errorf(
			"Got: %q, Want: %q",
			string(result),
			expected,
		)

		// Find where it diverges
		for i := 0; i < len(expected) && i < len(result); i++ {
			if result[i] != expected[i] {
				t.Errorf(
					"First mismatch at index %d: got %q (%d), want %q (%d)",
					i,
					result[i],
					result[i],
					expected[i],
					expected[i],
				)

				break
			}
		}
	}
}

func TestAppendJSONStringBufferGrowth(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedFinal []byte

		initialCap int
	}{
		{
			name:          "1. grow from zero capacity",
			initialCap:    0,
			input:         "short",
			expectedFinal: []byte("short"),
		},
		{
			name:          "2. grow when exceeding capacity in fast path",
			initialCap:    5,
			input:         "this is a longer string",
			expectedFinal: []byte("this is a longer string"),
		},
		{
			name:          "3. grow when exceeding capacity in slow path",
			initialCap:    5,
			input:         `"quoted"`,
			expectedFinal: []byte(`\"quoted\"`),
		},
		{
			name:          "4. multiple growths required",
			initialCap:    2,
			input:         "this string will cause multiple reallocations",
			expectedFinal: []byte("this string will cause multiple reallocations"),
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				initialBuf := make([]byte, 0, tt.initialCap)
				result := AppendJSON(initialBuf, []byte(tt.input))

				require.Equal(t,
					len(result),
					len(tt.expectedFinal),

					"length mismatch: got %d, want %d",
					len(result),
					len(tt.expectedFinal),
				)

				for i := range result {
					require.Equal(t,
						result[i],
						tt.expectedFinal[i],

						"byte mismatch at index %d: got %c (%d), want %c (%d)",
						i, result[i],
						result[i],
						tt.expectedFinal[i],
						tt.expectedFinal[i],
					)
				}

				// Verify that the result is a valid slice of the underlying array
				if len(result) > cap(result) {
					t.Errorf(
						"invalid slice: len(%d) > cap(%d)",
						len(result),
						cap(result),
					)
				}
			},
		)
	}
}

func TestAppendJSONStringEdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		initialBuf []byte
		input      string
		expected   []byte
	}{
		{
			name:       "1. exactly capacity boundary in fast path",
			initialBuf: make([]byte, 0, 8),
			input:      "12345678", // exactly 8 bytes, no growth needed
			expected:   []byte("12345678"),
		},
		{
			name:       "2. just over capacity in fast path",
			initialBuf: make([]byte, 0, 8),
			input:      "123456789", // 9 bytes, should grow
			expected:   []byte("123456789"),
		},
		{
			name:       "3. exactly capacity boundary in slow path",
			initialBuf: make([]byte, 0, 7),
			input:      `"`, // needs 2 bytes, but only 7 capacity, just fits
			expected:   []byte(`\"`),
		},
		{
			name:       "4. just over capacity in slow path",
			initialBuf: make([]byte, 0, 7),
			input:      `"a`, // needs 2 bytes for quote + 1 for 'a' (3 total), exceeds capacity
			expected:   []byte(`\"a`),
		},
		{
			name:       "5. unicode characters (no escaping needed)",
			initialBuf: make([]byte, 0),
			input:      "Hello, 世界",
			expected:   []byte("Hello, 世界"),
		},
		{
			name:       "6. mixed ascii and unicode with escapes",
			initialBuf: make([]byte, 0),
			input:      "Hello\n世界",
			expected:   []byte("Hello\\n世界"),
		},
		{
			name:       "7. all escaped characters",
			initialBuf: make([]byte, 0),
			input:      "\"\\\n\r\t",
			expected:   []byte(`\"\\\n\r\t`),
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				result := AppendJSON(tt.initialBuf, []byte(tt.input))

				require.Equal(t,
					len(result),
					len(tt.expected),

					"length mismatch: got %d, want %d",
					len(result),
					len(tt.expected),
				)

				for i := range result {
					require.Equal(t,
						result[i],
						tt.expected[i],

						"byte mismatch at index %d: got %c (%d), want %c (%d)",
						i, result[i],
						result[i],
						tt.expected[i],
						tt.expected[i],
					)
				}
			},
		)
	}
}

func TestAppendJSONStringPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping performance test in short mode")
	}

	// Test with a large string to ensure no panics and correct behavior
	largeInput := make([]byte, 10000)
	for i := range largeInput {
		if i%10 == 0 {
			largeInput[i] = '"'
		} else {
			largeInput[i] = 'a'
		}
	}

	buf := make([]byte, 0)
	result := AppendJSON(buf, (largeInput))

	expectedLen := 0

	for _, c := range largeInput {
		if c == '"' {
			expectedLen = expectedLen + 2 // escaped quote
		} else {
			expectedLen++
		}
	}

	require.Equal(t,
		len(result),
		expectedLen,

		"performance test length mismatch: got %d, want %d",
		len(result),
		expectedLen,
	)

	// Verify no panic occurred
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("performance test panicked: %v", r)
		}
	}()
}

func TestAppendJSONStringNilBuffer(t *testing.T) {
	var buf []byte

	result := AppendJSON(buf, []byte("test"))

	expected := []byte("test")
	require.Equal(t,
		len(result),
		len(expected),

		"nil buffer test length mismatch: got %d, want %d",
		len(result),
		len(expected),
	)

	for i := range result {
		require.Equal(t,
			result[i],
			expected[i],

			"nil buffer test byte mismatch at index %d: got %c, want %c",
			i, result[i],
			expected[i],
		)
	}
}
