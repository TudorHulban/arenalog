package helpers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSprintfInt(t *testing.T) {
	tests := []struct {
		description string
		format      string
		want        string

		args []int
	}{
		{
			description: "01: no substitutions",
			format:      "abc",
			args:        nil,
			want:        "abc",
		},
		{
			description: "02: missing args",
			format:      "x=%d y=%d",
			args:        []int{5},
			want:        "x=5 y=%d",
		},
		{
			description: "03: zero value",
			format:      "v=%d",
			args:        []int{0},
			want:        "v=0",
		},
		{
			description: "04: negative value",
			format:      "n=%d",
			args:        []int{-42},
			want:        "n=-42",
		},
		{
			description: "05: multiple ints",
			format:      "a=%d b=%d c=%d",
			args:        []int{1, 22, 333},
			want:        "a=1 b=22 c=333",
		},
		{
			description: "06: mixed literal and ints",
			format:      "x=%d end",
			args:        []int{99},
			want:        "x=99 end",
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.description,
			func(t *testing.T) {
				got := SprintfInt(tc.format, tc.args...)
				require.Equal(t, tc.want, got)
			},
		)
	}
}

// BenchmarkHelper_SprintfInt-16    	27305144	        42.92 ns/op	      24 B/op	       1 allocs/op
func BenchmarkHelper_SprintfInt(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		_ = SprintfInt(
			`<label>%d:</label>`,

			7,
		)
	}
}

// BenchmarkLibrary_SprintfInt-16    	18936204	        64.20 ns/op	      24 B/op	       1 allocs/op
func BenchmarkLibrary_SprintfInt(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		_ = fmt.Sprintf(
			`<label>%d:</label>`,

			7,
		)
	}
}
