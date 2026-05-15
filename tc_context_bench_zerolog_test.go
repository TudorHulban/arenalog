package arenalog

import (
	"fmt"
	"os"
	"runtime"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena/helpers"
)

func TestZerolog_OneField(t *testing.T) {
	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Logger()

	// {"level":"info","area":"some area","time":"2026-04-28T17:33:32+03:00","message":"benchmark test"}

	logger.Info().
		Str("area", "some area").
		Msg("benchmark test")
}

// cpu: AMD Ryzen 5 5600U with Radeon Graphics
// BenchmarkZerolog_OneField/G1-12 	 7749836	       154.0 ns/op	       0 B/op	       0 allocs/op
// BenchmarkZerolog_OneField/G2-12 	 7865106	       153.2 ns/op	       0 B/op	       0 allocs/op
// BenchmarkZerolog_OneField/G3-12 	 7710064	       157.0 ns/op	       0 B/op	       0 allocs/op
// BenchmarkZerolog_OneField/G4-12 	 7761795	       154.6 ns/op	       0 B/op	       0 allocs/op

func BenchmarkZerolog_OneField(b *testing.B) {
	gomaxprocsValues := []int{1, 2, 3, 4}
	writer := helpers.CountWriterNoBuffer{}

	for _, g := range gomaxprocsValues {
		b.Run(
			fmt.Sprintf("G%d", g),
			func(b *testing.B) {
				prev := runtime.GOMAXPROCS(g)
				defer runtime.GOMAXPROCS(prev)

				logger := zerolog.New(&writer).With().
					Timestamp().
					Logger()

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

// BenchmarkZerolog_Parallel_OneField-16    	 7567926	       156.4 ns/op	       0 B/op	       0 allocs/op
func BenchmarkZerolog_Parallel_OneField(b *testing.B) {
	runtime.GOMAXPROCS(1)

	writer := helpers.CountWriterNoBuffer{}

	logger := zerolog.New(&writer).With().
		Timestamp().
		Logger()

	b.ReportAllocs()
	b.ResetTimer()
	b.SetParallelism(16)

	b.RunParallel(
		func(pb *testing.PB) {
			for pb.Next() {
				logger.Info().
					Str("area", "some area").
					Msg("benchmark test")
			}
		},
	)

	require.NotZero(b,
		writer.TotalBytesWritten.Load(), // force writer to stay live
	)
}

// BenchmarkZerolog_WithFields-16    	 4206484	       285.4 ns/op	       0 B/op	       0 allocs/op
func BenchmarkZerolog_WithFields(b *testing.B) {
	var writer helpers.NoopWriter

	logger := zerolog.New(&writer).With().
		Timestamp().
		Str("service", "auth").
		Int("req_id", 12345).
		Bool("cache_hit", true).
		Logger()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		logger.Info().
			Str("area", "some area").
			Str("user", "tudor").
			Int("attempt", i).
			Float64("some float", 1.1137).
			Bool("success", true).
			Msg("benchmark test")
	}
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkZerolog_WithFields_Parallel/gomaxprocs=1-16         	 4390494	       272.3 ns/op	       0 B/op	       0 allocs/op
// BenchmarkZerolog_WithFields_Parallel/gomaxprocs=2-16         	 8307152	       145.7 ns/op	       0 B/op	       0 allocs/op
// BenchmarkZerolog_WithFields_Parallel/gomaxprocs=3-16         	12156778	       100.5 ns/op	       0 B/op	       0 allocs/op
// BenchmarkZerolog_WithFields_Parallel/gomaxprocs=4-16         	15548028	        76.10 ns/op	       0 B/op	       0 allocs/op
func BenchmarkZerolog_WithFields_Parallel(b *testing.B) {
	gomaxprocsValues := []int{1, 2, 3, 4}

	writer := helpers.CountWriterNoBuffer{}

	logger := zerolog.New(&writer).With().
		Timestamp().
		Str("service", "auth").
		Int("req_id", 12345).
		Bool("cache_hit", true).
		Logger()

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
