package arenalog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppendJSONKV(t *testing.T) {
	type testCase struct { //nolint:govet
		description string
		level       string
		file        string
		line        int
		msg         string
		kv          []any
		expected    string
	}

	tests := []testCase{
		{
			description: "1. Basic message with one key and int value",
			level:       "TRACE",
			file:        "/tmp/file.go",
			line:        42,
			msg:         "hello",
			kv:          []any{"key1", 1},
			expected:    `{"level":"TRACE","caller":"/tmp/file.go","line":42,"msg":"hello","key1":1}`,
		},
		{
			description: "2. Multiple keys with mixed types",
			level:       "INFO",
			file:        "/tmp/x.go",
			line:        7,
			msg:         "multi",
			kv:          []any{"a", "x", "b", 2, "c", true},
			expected:    `{"level":"INFO","caller":"/tmp/x.go","line":7,"msg":"multi","a":"x","b":2,"c":true}`,
		},
		{
			description: "3. Nil value",
			level:       "WARN",
			file:        "/tmp/y.go",
			line:        99,
			msg:         "nil-test",
			kv:          []any{"n", nil},
			expected:    `{"level":"WARN","caller":"/tmp/y.go","line":99,"msg":"nil-test","n":null}`,
		},
		{
			description: "4. No key-values",
			level:       "ERROR",
			file:        "/tmp/z.go",
			line:        1,
			msg:         "empty",
			kv:          nil,
			expected:    `{"level":"ERROR","caller":"/tmp/z.go","line":1,"msg":"empty"}`,
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.description,
			func(t *testing.T) {
				var l Logger

				buf := make([]byte, 0, 256)

				out := l.appendJSONKV(
					buf,
					tc.level,
					tc.file,
					tc.line,
					[]byte(tc.msg),
					tc.kv...,
				)

				require.Equal(t,
					tc.expected+"\n",
					string(out),
				)
			},
		)
	}
}
