package query

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRawset(t *testing.T) {
	input := `2026-05-12T13:43:58+03:00 /path/file.go Line68 service=auth req_id=12345 cache_hit=true area="some area" g=8 msg=finished`

	logSet, errCr := NewRawset(input)
	require.NoError(t, errCr)
	require.Len(t, logSet, 1)

	entry := logSet[0]

	// 1. Verify Timestamp
	require.Equal(t,
		"2026-05-12T13:43:58+03:00",
		entry.timestamp,

		"timestamp extraction failed",
	)

	// 2. Verify String Values
	require.Equal(t,
		"auth",
		entry.keyValues["service"],
	)

	// 3. Verify Quoted Strings
	require.Equal(t,
		"some area",
		entry.keyValues["area"],

		"quoted string parsing failed",
	)

	// 4. Verify Booleans
	require.Equal(t,
		true,
		entry.keyValues["cache_hit"],
	)

	// 5. Verify Numbers (float64 to match JSON parser)
	require.Equal(t,
		8.0,
		entry.keyValues["g"],
	)
}
