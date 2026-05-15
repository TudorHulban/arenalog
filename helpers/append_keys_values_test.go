package helpers

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppendKeyValues(t *testing.T) {
	type testCase struct { //nolint:govet
		description string
		kv          []any
		expected    string
	}

	tests := []testCase{
		{
			description: "1. Basic strings and integers",
			kv:          []any{"user", "bob", "id", 123},
			expected:    "user=bob id=123",
		},
		{
			description: "2. Byte slices for keys and values",
			kv:          []any{[]byte("raw_key"), []byte("raw_val")},
			expected:    "raw_key=raw_val",
		},
		{
			description: "3. Various integer widths and signs",
			kv: []any{
				"u64", uint64(500),
				"i32", int32(-10),
				"bool_val", true,
			},
			expected: "u64=500 i32=-10 bool_val=true",
		},
		{
			description: "4. Floats and Nils",
			kv: []any{
				"pi", 3.14159,
				"nothing", nil,
			},
			// Note: Float precision depends on your AppendFloat implementation
			expected: "pi=3.141589999999 nothing=null",
		},
		{
			description: "5. Errors and Exotic types",
			kv: []any{
				"err", errors.New("database timeout"),
				"map", map[string]int{"a": 1},
			},
			expected: "err=database timeout map=map[a:1]",
		},
		{
			description: "6. Non-string keys (cold path)",
			kv: []any{
				100, "hundred",
				false, "is_false",
			},
			expected: "100=hundred false=is_false",
		},
		{
			description: "7. Single pair",
			kv:          []any{"key", "value"},
			expected:    "key=value",
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.description,
			func(t *testing.T) {
				buf := make([]byte, 0, 128)
				out := AppendKeyValues(buf, tc.kv...)

				require.Equal(t, tc.expected, string(out))
			},
		)
	}
}

func TestAppendJSONKeyValuesIntoObject(t *testing.T) {
	type testCase struct { //nolint:govet
		description string
		kv          []any
		expected    string
	}

	tests := []testCase{
		{
			description: "1. Standard string and integer types",
			kv:          []any{"name", "alice", "age", 30},
			expected:    `"name":"alice","age":30,`,
		},
		{
			description: "2. Byte slices for keys and values",
			kv:          []any{[]byte("hex"), []byte("deadbeef")},
			expected:    `"hex":"deadbeef",`,
		},
		{
			description: "3. Booleans, floats, and nulls",
			kv: []any{
				"enabled", true,
				"score", 95.5,
				"meta", nil,
			},
			expected: `"enabled":true,"score":95.500000000000,"meta":null,`,
		},
		{
			description: "4. Error types (should be quoted strings)",
			kv:          []any{"status", errors.New("connection_failed")},
			expected:    `"status":"connection_failed",`,
		},
		{
			description: "5. Non-string keys and exotic values (cold path)",
			kv: []any{
				123, "numeric_key",
				"data", []int{1, 2, 3}, // fmt.Sprint fallback
			},
			expected: `"123":"numeric_key","data":"[1 2 3]",`,
		},
		{
			description: "6. Empty slice (should do nothing)",
			kv:          []any{},
			expected:    "",
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.description,
			func(t *testing.T) {
				buf := make([]byte, 0, 256)
				out := AppendJSONKeyValuesIntoObject(buf, tc.kv...)

				require.Equal(t, tc.expected, string(out))
			},
		)
	}
}
