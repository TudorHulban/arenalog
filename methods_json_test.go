package arenalog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppendJSON(t *testing.T) {
	type testCase struct { //nolint:govet
		description string
		timestamp   string
		level       string
		file        string
		line        int
		msg         string
		expected    string
	}

	tests := []testCase{
		{
			description: "1. Full log entry with timestamp and caller",
			timestamp:   "2026-05-11T18:00:00Z",
			level:       "INFO",
			file:        "/src/main.go",
			line:        100,
			msg:         "application started",
			expected:    `{"ts":"2026-05-11T18:00:00Z","level":"INFO","caller":"/src/main.go","line":100,"msg":"application started"}`,
		},
		{
			description: "2. No timestamp, only level and message",
			timestamp:   "",
			level:       "DEBUG",
			file:        "",
			line:        0,
			msg:         "debugging simple logic",
			expected:    `{"level":"DEBUG","msg":"debugging simple logic"}`,
		},
		{
			description: "3. No caller info (empty file)",
			timestamp:   "2026-05-11T18:05:00Z",
			level:       "WARN",
			file:        "",
			line:        50, // Line exists but file is empty
			msg:         "warning message",
			expected:    `{"ts":"2026-05-11T18:05:00Z","level":"WARN","msg":"warning message"}`,
		},
		{
			description: "4. Message requiring JSON escaping",
			timestamp:   "",
			level:       "ERROR",
			file:        "app.go",
			line:        12,
			msg:         `found a "quote" and a \ backslash`,
			expected:    `{"level":"ERROR","caller":"app.go","line":12,"msg":"found a \"quote\" and a \\ backslash"}`,
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.description,
			func(t *testing.T) {
				l := &Logger{}

				if tc.timestamp != "" {
					l.fnTimestamp = func(b []byte) []byte {
						return append(b, tc.timestamp...)
					}
				}

				buf := make([]byte, 0, 512)
				out := l.appendJSON(
					buf,
					tc.level,
					tc.file,
					tc.line,
					[]byte(tc.msg),
				)

				// We expect the result to have a trailing newline
				require.Equal(t, tc.expected+"\n", string(out))
			},
		)
	}
}
