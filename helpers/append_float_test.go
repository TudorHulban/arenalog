package helpers

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppendFloat(t *testing.T) {
	tests := []struct { //nolint:govet
		description string
		value       float64
		prec        int
		expected    string
	}{
		// 1. Error-like cases
		{"nan", math.NaN(), 4, "nan"},
		{"pos_inf", math.Inf(+1), 4, "inf"},
		{"neg_inf", math.Inf(-1), 4, "-inf"},

		// 2. Normal cases (truncate, do not round)
		{"small_int", 42, 0, "42"},
		{"small_frac", 3.14159, 4, "3.1415"},
		{"large_int", 123456789012345.0, 0, "123456789012345"},
		{"large_frac", 987654321.123456, 6, "987654321.123456"},
		{"tiny", 0.00000012345, 10, "0.0000001234"},
		{"negative", -123.456, 3, "-123.456"},
	}

	for _, tc := range tests {
		t.Run(
			tc.description,
			func(t *testing.T) {
				var dst []byte

				out := AppendFloat(dst, tc.value, tc.prec)

				// Normalize to lowercase for safety
				actual := strings.ToLower(string(out))
				expected := strings.ToLower(tc.expected)

				// NaN requires special handling
				if math.IsNaN(tc.value) {
					require.Equal(t, expected, actual)
					return
				}

				require.Equal(t, expected, actual)
			},
		)
	}
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkAppendFloat/1_custom_1._nan-16         	 8278422	       145.4 ns/op	       8 B/op	       1 allocs/op
// BenchmarkAppendFloat/2_stdlib_1._nan-16         	 8229361	       143.4 ns/op	       8 B/op	       1 allocs/op

// BenchmarkAppendFloat/1_custom_2._pos_inf-16     	 8517171	       142.3 ns/op	       8 B/op	       1 allocs/op
// BenchmarkAppendFloat/2_stdlib_2._pos_inf-16     	 8217518	       143.7 ns/op	       8 B/op	       1 allocs/op

// BenchmarkAppendFloat/1_custom_3._neg_inf-16     	 8315826	       141.4 ns/op	       8 B/op	       1 allocs/op
// BenchmarkAppendFloat/2_stdlib_3._neg_inf-16     	 8195685	       144.3 ns/op	       8 B/op	       1 allocs/op

// BenchmarkAppendFloat/1_custom_4._small_int-16   	 7951039	       149.9 ns/op	       8 B/op	       1 allocs/op
// BenchmarkAppendFloat/2_stdlib_4._small_int-16   	 6827728	       177.1 ns/op	       8 B/op	       1 allocs/op

// BenchmarkAppendFloat/1_custom_5._small_frac-16  	 7545920	       158.1 ns/op	       8 B/op	       1 allocs/op
// BenchmarkAppendFloat/2_stdlib_5._small_frac-16  	 6271718	       188.8 ns/op	       8 B/op	       1 allocs/op

// BenchmarkAppendFloat/1_custom_6._large_int-16   	 7592638	       158.5 ns/op	      16 B/op	       1 allocs/op
// BenchmarkAppendFloat/2_stdlib_6._large_int-16   	 6490910	       186.0 ns/op	      16 B/op	       1 allocs/op

// BenchmarkAppendFloat/1_custom_7._large_frac-16  	 7291576	       164.7 ns/op	      16 B/op	       1 allocs/op
// BenchmarkAppendFloat/2_stdlib_7._large_frac-16  	 6132566	       196.2 ns/op	      16 B/op	       1 allocs/op

// BenchmarkAppendFloat/1_custom_8._tiny-16        	 6586474	       181.8 ns/op	      24 B/op	       2 allocs/op
// BenchmarkAppendFloat/2_stdlib_8._tiny-16        	 5570665	       214.9 ns/op	      24 B/op	       2 allocs/op

// BenchmarkAppendFloat/1_custom_9._negative-16    	 7618142	       156.9 ns/op	       8 B/op	       1 allocs/op
// BenchmarkAppendFloat/2_stdlib_9._negative-16    	 6682270	       178.5 ns/op	       8 B/op	       1 allocs/op

func BenchmarkAppendFloat(b *testing.B) {
	tests := []struct {
		description string
		value       float64
		prec        int
	}{
		// Error-like cases
		{"1. nan", math.NaN(), 4},
		{"2. pos_inf", math.Inf(+1), 4},
		{"3. neg_inf", math.Inf(-1), 4},

		// Normal cases
		{"4. small_int", 42, 0},
		{"5. small_frac", 3.14159, 4},
		{"6. large_int", 123456789012345.0, 0},
		{"7. large_frac", 987654321.123456, 6},
		{"8. tiny", 0.00000012345, 10},
		{"9. negative", -123.456, 3},
	}

	for _, tc := range tests {
		fmt.Println("")

		b.Run(
			"1_custom_"+tc.description,
			func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()

				var dst []byte

				for b.Loop() {
					dst = dst[:0]

					out := AppendFloat(dst, tc.value, tc.prec)
					require.NotNil(b, out)
				}
			},
		)

		b.Run(
			"2_stdlib_"+tc.description,
			func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()

				var dst []byte

				for b.Loop() {
					dst = dst[:0]

					out := strconv.AppendFloat(dst, tc.value, 'f', tc.prec, 64)
					require.NotNil(b, out)
				}
			},
		)
	}
}
