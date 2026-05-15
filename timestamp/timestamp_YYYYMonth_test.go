package timestamp

import (
	"bytes"
	"testing"
	"time"
)

func TestTimestampYYYYMonth(t *testing.T) {
	chReady := StartYYYYMonthCache(t.Context())

	<-chReady

	// The format produced is: YYYYMM DD HH:MM:SS.mmm
	// Example: 202405 21 15:04:05.123
	const layout = "200601 02 15:04:05.000"

	t.Run(
		"1. Basic format and length validation",
		func(t *testing.T) {
			got := TimestampYYYYMonth(nil)

			// YYYYMM(6) + space(1) + DD(2) + space(1) + HH:MM:SS.mmm(12) = 22 characters
			if len(got) != 22 {
				t.Errorf("Expected length 22, got %d (%s)", len(got), string(got))
			}

			// Validate structure via parsing
			_, errParse := time.Parse(layout, string(got))
			if errParse != nil {
				t.Errorf("Result does not match YYYYMonth layout: %v", errParse)
			}
		},
	)

	t.Run(
		"2. Functional integrity with prefix",
		func(t *testing.T) {
			prefix := []byte("LOG_START|")
			got := TimestampYYYYMonth(prefix)

			if !bytes.HasPrefix(got, prefix) {
				t.Fatalf("Prefix corrupted. Got: %s", string(got))
			}

			timestampPart := got[len(prefix):]
			if len(timestampPart) != 22 {
				t.Errorf("Timestamp part length mismatch. Got %d", len(timestampPart))
			}

			_, errParse := time.Parse(layout, string(timestampPart))
			if errParse != nil {
				t.Errorf("Appended timestamp is invalid: %v", errParse)
			}
		},
	)

	t.Run(
		"3. Cache TTL and Gate verification",
		func(t *testing.T) {
			// First capture
			t1 := string(TimestampYYYYMonth(nil))

			// Sleep to bypass the 1ms gateStandard.CompareAndSwap
			time.Sleep(2 * time.Millisecond)

			// Second capture
			t2 := string(TimestampYYYYMonth(nil))

			if t1 == t2 {
				t.Errorf(
					"Cache failed to update after 2ms sleep.\nT1: %s\nT2: %s",
					t1,
					t2,
				)
			}
		},
	)
}
