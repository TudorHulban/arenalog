package query

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewLogset(t *testing.T) {
	input := `
Starting application...
{"ts": "2023-10-01T10:00:00Z", "level": "info", "msg": "hello world", "status": 200}
some random noise between logs
{"ts": "2023-10-01T10:00:01Z", "level": "error", "msg": "failed task", "retry": true}
`

	logSet, errCr := NewLogset(input)
	require.NoError(t,
		errCr,
		"failed to parse multi log",
	)

	// 1. Verify total count (4 non-empty lines)
	expectedCount := 4
	require.Len(t,
		logSet,
		expectedCount,

		"expected %d log records, got %d",
		expectedCount,
		len(logSet),
	)

	// 2. Test Raw Entry (First line)
	require.True(t,
		logSet[0].IsRAW(),

		"first line should be RAW",
	)

	require.Equal(t,
		"Starting application...",

		logSet[0].String(),
		"unexpected raw content",
	)

	// 3. Test JSON Parsing (Second line)
	second := logSet[1]
	require.Equal(t,
		"2023-10-01T10:00:00Z",

		second.timestamp,
		"expected timestamp mismatch",
	)

	exists, val := second.HasKey("level")
	require.True(t, exists, "key 'level' should exist")
	require.Equal(t,
		"info",
		val,
		"expected key 'level' to be 'info'",
	)

	// 4. Test Types (integers/floats in JSON)
	// Note: encoding/json unmarshals numbers into float64 by default in map[string]any
	_, statusVal := second.HasKey("status")
	require.Equal(t,
		200.0,
		statusVal,

		"expected status 200 (float64), got %v (%T)",
		statusVal,
		statusVal,
	)

	// 5. Verify filtering functionality works with the result
	withTs := logSet.WithTimestamp()
	require.Len(t,
		withTs,
		2,

		"expected 2 logs with timestamps",
	)
}

func TestNewLogset_AnsiCleaning(t *testing.T) {
	// Input with ANSI color codes: [32m is green
	input := "\x1b[32m{\"ts\": \"2023-10-01T10:00:00Z\", \"msg\": \"colored\"}\x1b[0m"

	logs, errCr := NewLogset(input)
	require.NoError(t,
		errCr,
		"failed to parse colored log",
	)

	require.Len(t,
		logs,
		1,

		"expected 1 log record",
	)

	require.Equal(t,
		"2023-10-01T10:00:00Z",
		logs[0].timestamp,

		"failed to parse JSON after ANSI stripping, got ts: %q",
		logs[0].timestamp,
	)

	// Ensure the 'raw' field preserves the original ANSI for display/debugging
	require.Contains(t,
		logs[0].raw,
		"\x1b[32m",

		"raw field should preserve ANSI codes",
	)
}

func TestNewLogset_Empty(t *testing.T) {
	logs, errCr := NewLogset("   \n\n   ")
	require.NoError(t,
		errCr,
		"expected no error on empty input",
	)

	require.Empty(t,
		logs,

		"expected 0 logs, got %d",
		len(logs),
	)
}
