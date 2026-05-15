package timestamp

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTimestampRFC3339(t *testing.T) {
	chReady := StartRFC3339UTCCache(t.Context())

	<-chReady

	t.Run(
		"1. Test basic appending to an empty slice",
		func(t *testing.T) {
			got := TimestampRFC3339UTC(nil)

			// Validate format by attempting to parse it back
			_, errParse := time.Parse(time.RFC3339Nano, string(got))
			require.NoError(t,
				errParse,

				"Result is not a valid RFC3339 timestamp: %v",
				errParse,
			)
		},
	)

	t.Run(
		"2. Test appending to an existing slice (functional integrity)",
		func(t *testing.T) {
			prefix := []byte("log_prefix: ")
			got := TimestampRFC3339UTC(prefix)

			if !bytes.HasPrefix(got, prefix) {
				t.Fatalf("Expected prefix not found. Got: %s", string(got))
			}

			timestampPart := got[len(prefix):]

			_, errParse := time.Parse(time.RFC3339Nano, string(timestampPart))
			require.NoError(t,
				errParse,

				"Appended part is not a valid timestamp: %v",
				errParse,
			)
		},
	)

	t.Run(
		"3. Consistency check (values should be close to time.Now)",
		func(t *testing.T) {
			now := time.Now().UTC()
			got := TimestampRFC3339UTC(nil)

			parsed, _ := time.Parse(time.RFC3339Nano, string(got))

			// Check if the timestamp is within a reasonable drift (e.g., 1 second)
			// because the cache update might be slightly behind time.Now()
			diff := parsed.Sub(now)
			if diff < 0 {
				diff = -diff
			}

			if diff > time.Second {
				t.Errorf("Timestamp drift too high: %v", diff)
			}
		},
	)
}
