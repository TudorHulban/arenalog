package query

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogSet_Methods(t *testing.T) {
	input := "RAW entry 1\n" +
		`{"ts": "2026-05-12T12:00:00Z", "msg": "json entry"}` + "\n" +
		"RAW entry 2"

	logSet, errCr := NewLogset(input)
	require.NoError(t, errCr)
	require.Len(t, logSet, 3)

	t.Run(
		"1. Stringer",
		func(t *testing.T) {
			// String() should join the raw values back with newlines
			expected := "RAW entry 1\n" +
				`{"ts": "2026-05-12T12:00:00Z", "msg": "json entry"}` + "\n" +
				"RAW entry 2"

			require.Equal(t,
				expected,
				logSet.String(),

				"String() output should match original raw lines joined by newline",
			)
		},
	)

	t.Run(
		"2. WithTimestamp",
		func(t *testing.T) {
			filtered := logSet.WithTimestamp()

			require.Len(t,
				filtered,
				1,
				"expected 1 entry with timestamp",
			)

			require.Equal(t,
				"2026-05-12T12:00:00Z",
				filtered[0].timestamp,
			)
		},
	)

	t.Run(
		"3. WithNoTimestamp",
		func(t *testing.T) {
			filtered := logSet.WithNoTimestamp()

			require.Len(t,
				filtered,
				2,
				"expected 2 entries with no timestamp",
			)

			require.True(t,
				filtered[0].IsRAW(),
				"first filtered entry should be RAW",
			)

			require.True(t,
				filtered[1].IsRAW(),
				"second filtered entry should be RAW",
			)
		},
	)
}

func TestLogSet_KeyVerification(t *testing.T) {
	input := `{"ts": "1", "app": "web", "cluster": "us-east", "status": 200}
{"ts": "2", "app": "db", "cluster": "us-east", "status": 500}
{"ts": "3", "app": "web", "cluster": "us-west", "status": 200}
non-json-record`

	logSet, errCr := NewLogset(input)
	require.NoError(t, errCr)

	t.Run(
		"1. HasKey",
		func(t *testing.T) {
			// 'app' exists in 3 records
			require.NoError(t, logSet.HasKey("app", 3))

			// 'status' exists in 3 records
			require.NoError(t, logSet.HasKey("status", 3))

			// Error case: 'cluster' exists 3 times, not 1
			errCluster := logSet.HasKey("cluster", 1)
			require.Error(t, errCluster)
			require.Contains(t,
				errCluster.Error(),

				"expected key \"cluster\" to appear 1 times, but found it 3 times",
			)
		},
	)

	t.Run(
		"2. HasKeyWithValue",
		func(t *testing.T) {
			// 'app' is 'web' exactly 2 times
			require.NoError(t, logSet.HasKeyWithValue("app", "web", 2))

			// JSON numbers are float64
			require.NoError(t, logSet.HasKeyWithValue("status", 200.0, 2))

			// Error case: mismatch count
			require.Error(t, logSet.HasKeyWithValue("app", "db", 5))
		},
	)

	t.Run(
		"3. HasKeysWithValues",
		func(t *testing.T) {
			// Match multiple keys in a single record: app=web AND cluster=us-east (found once in record 1)
			require.NoError(t,
				logSet.HasKeysWithValues(1, "app", "web", "cluster", "us-east"),
			)

			// Validate argument count error
			errNoArguments := logSet.HasKeysWithValues(1, "app")
			require.Error(t, errNoArguments)
			require.Equal(t,
				"hasKeysWithValues requires an even number of kv arguments",
				errNoArguments.Error(),
			)
		},
	)

	t.Run(
		"4. FilterBy",
		func(t *testing.T) {
			filtered := logSet.FilterBy("app", "web")
			require.Len(t,
				filtered,
				2,
				"expected 2 records after filtering by app=web",
			)
		},
	)

	t.Run(
		"5. FirstAndLast",
		func(t *testing.T) {
			// Test First
			require.Equal(t, "1", logSet.First()[0].timestamp)

			// Test Last
			last := logSet.Last()[0]
			require.True(t,
				last.IsRAW(),
				"last record should be the non-json string",
			)
			require.Equal(t, "non-json-record", last.String())

			// Test Empty Set
			var emptySet LogSet

			require.Empty(t, emptySet.First().String())
			require.Empty(t, emptySet.Last().String())
		},
	)
}

func TestLogSet_SortByTimestamp(t *testing.T) {
	input := `{"ts": "2026-05-12T10:00:01Z", "msg": "middle"}
{"ts": "2026-05-12T10:00:00Z", "msg": "earliest"}
{"ts": "2026-05-12T10:00:02Z", "msg": "latest"}`

	logSet, errCr := NewLogset(input)
	require.NoError(t, errCr)
	require.Len(t, logSet, 3)

	t.Run(
		"1. Ascending",
		func(t *testing.T) {
			sorted := logSet.SortByTimestamp(false)

			require.Len(t, sorted, 3)
			require.Equal(t, "2026-05-12T10:00:00Z", sorted[0].timestamp)
			require.Equal(t, "2026-05-12T10:00:01Z", sorted[1].timestamp)
			require.Equal(t, "2026-05-12T10:00:02Z", sorted[2].timestamp)

			// Verify original logSet remains unmutated
			require.Equal(t,
				"2026-05-12T10:00:01Z",
				logSet[0].timestamp,
				"original logSet should not be mutated by SortByTimestamp",
			)
		},
	)

	t.Run(
		"2. Descending",
		func(t *testing.T) {
			sorted := logSet.SortByTimestamp(true)

			require.Len(t, sorted, 3)
			require.Equal(t, "2026-05-12T10:00:02Z", sorted[0].timestamp)
			require.Equal(t, "2026-05-12T10:00:01Z", sorted[1].timestamp)
			require.Equal(t, "2026-05-12T10:00:00Z", sorted[2].timestamp)
		},
	)
}
