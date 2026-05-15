package arenalog

import (
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/phuslu/log"
	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena/helpers"
)

func TestPhuslu_OneField(t *testing.T) {
	logger := log.Logger{
		Level:      log.InfoLevel,
		TimeFormat: time.RFC3339,
		Writer: &log.IOWriter{
			Writer: os.Stdout,
		},
	}

	// {"time":"2026-04-28T12:09:35+03:00","level":"info","area":"some area","message":"benchmark test"}

	logger.Info().
		Str("area", "some area").
		Msg("benchmark test")
}

// go test -run '^$' -bench '^BenchmarkPhuslu_OneField$' -benchmem

// cpu: AMD Ryzen 5 5600U with Radeon Graphics
// BenchmarkPhuslu_OneField/G1-12 	 8532501	       140.0 ns/op	       0 B/op	       0 allocs/op
// BenchmarkPhuslu_OneField/G2-12 	 8480228	       139.4 ns/op	       0 B/op	       0 allocs/op
// BenchmarkPhuslu_OneField/G3-12 	 8606548	       140.0 ns/op	       0 B/op	       0 allocs/op
// BenchmarkPhuslu_OneField/G4-12 	 8446911	       141.2 ns/op	       0 B/op	       0 allocs/op

func BenchmarkPhuslu_OneField(b *testing.B) {
	gomaxprocsValues := []int{1, 2, 3, 4}
	writer := helpers.CountWriterNoBuffer{}

	for _, g := range gomaxprocsValues {
		b.Run(
			fmt.Sprintf("G%d", g),
			func(b *testing.B) {
				prev := runtime.GOMAXPROCS(g)
				defer runtime.GOMAXPROCS(prev)

				logger := log.Logger{
					Level:      log.InfoLevel,
					TimeFormat: time.RFC3339,
					Writer:     &log.IOWriter{Writer: &writer},
				}

				b.ReportAllocs()
				b.ResetTimer()

				for b.Loop() {
					logger.Info().
						Str("area", "some area").
						Msg("benchmark test")
				}

				require.NotZero(b, writer.TotalBytesWritten.Load())
			},
		)
	}
}

// cpu: AMD Ryzen 5 5600U with Radeon Graphics
// BenchmarkPhuslu_Parallel_OneField/gomaxprocs=1-12         	 8586711	       139.1 ns/op	       0 B/op	       0 allocs/op
// BenchmarkPhuslu_Parallel_OneField/gomaxprocs=2-12         	15442574	        76.53 ns/op	       0 B/op	       0 allocs/op
// BenchmarkPhuslu_Parallel_OneField/gomaxprocs=3-12         	23493219	        51.19 ns/op	       0 B/op	       0 allocs/op
// BenchmarkPhuslu_Parallel_OneField/gomaxprocs=4-12         	29966448	        39.28 ns/op	       0 B/op	       0 allocs/op

func BenchmarkPhuslu_Parallel_OneField(b *testing.B) {
	gomaxprocsValues := []int{1, 2, 3, 4}

	writer := helpers.CountWriterNoBuffer{}

	logger := log.Logger{
		Level:      log.InfoLevel,
		TimeFormat: time.RFC3339,
		Writer:     &log.IOWriter{Writer: &writer},
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.SetParallelism(16)

	for _, g := range gomaxprocsValues {
		inner := g

		b.Run(
			fmt.Sprintf("gomaxprocs=%d", g),
			func(b *testing.B) {
				prev := runtime.GOMAXPROCS(inner)
				defer runtime.GOMAXPROCS(prev)

				b.RunParallel(
					func(pb *testing.PB) {
						for pb.Next() {
							logger.Info().
								Str("area", "some area").
								Msg("benchmark test")
						}
					},
				)
			},
		)
	}

	require.NotZero(
		b,
		writer.TotalBytesWritten.Load(),
		"1. writer must record bytes",
	)
}

// cpu: AMD Ryzen 5 5600U with Radeon Graphics
// BenchmarkPhuslu_WithFields-12    	 4801941	       250.9 ns/op	       0 B/op	       0 allocs/op

func BenchmarkPhuslu_WithFields(b *testing.B) {
	var writer helpers.NoopWriter

	logger := log.Logger{
		Level:      log.InfoLevel,
		TimeField:  "ts",
		TimeFormat: time.RFC3339,
		Writer:     &log.IOWriter{Writer: &writer},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		logger.Info().
			Str("service", "auth").
			Int("req_id", 12345).
			Str("area", "some area").
			Str("user", "tudor").
			Int("attempt", i).
			Float64("some float", 1.1137).
			Bool("success", true).
			Msg("benchmark test")
	}
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkPhuslu_WithFields_Parallel/gomaxprocs=1-16         	 4837551	       244.7 ns/op	       0 B/op	       0 allocs/op
// BenchmarkPhuslu_WithFields_Parallel/gomaxprocs=2-16         	 8898666	       136.2 ns/op	       0 B/op	       0 allocs/op
// BenchmarkPhuslu_WithFields_Parallel/gomaxprocs=3-16         	12984573	        92.34 ns/op	       0 B/op	       0 allocs/op
// BenchmarkPhuslu_WithFields_Parallel/gomaxprocs=4-16         	17022438	        70.46 ns/op	       0 B/op	       0 allocs/op
func BenchmarkPhuslu_WithFields_Parallel(b *testing.B) {
	gomaxprocsValues := []int{1, 2, 3, 4}

	writer := helpers.CountWriterNoBuffer{}

	logger := log.Logger{
		Level:      log.InfoLevel,
		TimeField:  "ts",
		TimeFormat: time.RFC3339,
		Writer:     &log.IOWriter{Writer: &writer},
	}

	b.SetParallelism(16)

	for _, g := range gomaxprocsValues {
		inner := g

		b.Run(
			fmt.Sprintf("gomaxprocs=%d", inner),
			func(b *testing.B) {
				prev := runtime.GOMAXPROCS(inner)
				defer runtime.GOMAXPROCS(prev)

				b.ReportAllocs()
				b.ResetTimer()

				b.RunParallel(
					func(pb *testing.PB) {
						i := 0
						for pb.Next() {
							logger.Info().
								Str("service", "auth").
								Int("req_id", 12345).
								Str("area", "some area").
								Str("user", "tudor").
								Int("attempt", i).
								Float64("some float", 1.1137).
								Bool("success", true).
								Msg("benchmark test")

							i++
						}
					},
				)
			},
		)
	}

	require.NotZero(
		b,
		writer.TotalBytesWritten.Load(),
		"1. writer must record bytes",
	)
}
