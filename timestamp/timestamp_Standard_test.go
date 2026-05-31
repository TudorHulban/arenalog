package timestamp

import (
	"bytes"
	"testing"
	"time"
)

func TestTimestampStandard(t *testing.T) {
	const layout = "2006/01/02 15:04:05.000"

	chReady := StartStandardCache(t.Context())

	<-chReady

	t.Run(
		"1. Test basic appending and format validation",
		func(t *testing.T) {
			got := TimestampStandard(nil)

			// Validate length (YYYY/MM/DD HH:MM:SS.mmm is 23 characters)
			if len(got) != 23 {
				t.Errorf("Expected length 23, got %d (%s)", len(got), string(got))
			}

			// Validate format by parsing with the custom layout
			_, errParse := time.Parse(layout, string(got))
			if errParse != nil {
				t.Errorf(
					"Result does not match standard layout: %v",
					errParse,
				)
			}
		},
	)

	t.Run(
		"2. Test appending to existing slice",
		func(t *testing.T) {
			prefix := []byte("DATA|")
			got := TimestampStandard(prefix)

			if !bytes.HasPrefix(got, prefix) {
				t.Fatalf("Prefix lost. Got: %s", string(got))
			}

			timestampPart := got[len(prefix):]

			_, errParse := time.Parse(layout, string(timestampPart))
			if errParse != nil {
				t.Errorf(
					"Appended timestamp part is invalid: %v",
					errParse,
				)
			}
		},
	)

	t.Run(
		"3. Cache Refresh (Sequential calls)",
		func(t *testing.T) {
			// First capture
			t1 := string(TimestampStandard(nil))

			noMsSleep := 10

			// Sleep to bypass the 1ms gateStandard.CompareAndSwap
			time.Sleep(time.Duration(noMsSleep) * time.Millisecond)

			// Second capture
			t2 := string(TimestampStandard(nil))

			if t1 == t2 {
				t.Errorf(
					"Cache failed to update after %dms sleep.\nT1: %s\nT2: %s",
					noMsSleep,
					t1,
					t2,
				)
			}
		},
	)
}
