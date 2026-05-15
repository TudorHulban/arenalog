package helpers

import (
	"testing"
)

func makeLongString(length int, includeEscapes bool) string {
	result := make([]byte, length)

	for ix := range length {
		if includeEscapes && ix%10 == 0 {
			result[ix] = '"'
		} else {
			result[ix] = 'a'
		}
	}

	return string(result)
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkAppendJSONString/ShortNoEscape-16         	96213393	        12.22 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAppendJSONString/ShortWithEscapes-16      	43420896	        27.38 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAppendJSONString/MediumNoEscape-16        	22103844	        53.55 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAppendJSONString/MediumWithEscapes-16     	30883579	        38.97 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAppendJSONString/LongNoEscape-16          	  650764	      1837 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAppendJSONString/LongWithEscapes-16       	  542268	      2151 ns/op	    1536 B/op	       1 allocs/op

// Benchmark to test performance
func BenchmarkAppendJSONString(b *testing.B) {
	benchmarks := []struct {
		name  string
		input string
	}{
		{"ShortNoEscape", "hello"},
		{"ShortWithEscapes", `hello "world"`},
		{"MediumNoEscape", "abcdefghijklmnopqrstuvwxyz"},
		{"MediumWithEscapes", "line1\nline2\t\"quote\"\\"},
		{"LongNoEscape", makeLongString(1000, false)},
		{"LongWithEscapes", makeLongString(1000, true)},
	}

	for _, bm := range benchmarks {
		b.Run(
			bm.name,
			func(b *testing.B) {
				buf := make([]byte, 0, 1024)

				b.ResetTimer()

				for b.Loop() {
					_ = AppendJSON(buf, []byte(bm.input))
				}
			},
		)
	}
}
