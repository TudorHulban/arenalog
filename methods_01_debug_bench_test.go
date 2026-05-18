package arenalog

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/bytearena/helpers"
)

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// Benchmark_Debug/G1-16 	34433020	        34.52 ns/op	       0 B/op	       0 allocs/op
// Benchmark_Debug/G2-16 	15059820	        80.17 ns/op	       0 B/op	       0 allocs/op
// Benchmark_Debug/G3-16 	18178716	        66.70 ns/op	       0 B/op	       0 allocs/op
// Benchmark_Debug/G4-16 	17501660	        71.16 ns/op	       0 B/op	       0 allocs/op

func Benchmark_Debug(b *testing.B) {
	gomaxprocsValues := []int{1, 2, 3, 4}

	for _, g := range gomaxprocsValues {
		b.Run(
			fmt.Sprintf("G%d", g),
			func(b *testing.B) {
				prev := runtime.GOMAXPROCS(g)
				defer runtime.GOMAXPROCS(prev)

				ingestor, errCrIngestor := bytearena.NewIngestor(
					bytearena.Size100K(),
					io.Discard,

					helpers.TernaryWithValueIn(
						[]int{1},
						g,
						nil,
						bytearena.WithCounterCoreCPU(),
					),
				)
				require.NoError(b, errCrIngestor)
				require.NotNil(b, ingestor)

				l, errCrLogger := NewLogger(
					&ParamsNewLogger{
						Ingestor:    ingestor,
						LoggerLevel: LevelDebug,

						WithFatalWriter: os.Stdout,
						WithJSON:        true,
					},

					WithTimestampRFC3339UTC(b.Context()),
				)
				require.NoError(b, errCrLogger)
				require.NotNil(b, l)

				ctx, cancel := context.WithCancel(context.Background())
				chIngestionEnd := ingestor.StartIngestion(ctx)

				// warmup
				time.Sleep(10 * time.Millisecond)
				runtime.GC()

				b.SetParallelism(1)
				b.ReportAllocs()
				b.ResetTimer()

				b.RunParallel(
					func(pb *testing.PB) {
						for pb.Next() {
							l.Debug("1")
						}
					},
				)

				cancel()
				<-chIngestionEnd
			},
		)
	}
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// Benchmark_Debugf/G1-16 	30545702	        39.43 ns/op	       0 B/op	       0 allocs/op
// Benchmark_Debugf/G2-16 	26972127	        45.02 ns/op	       0 B/op	       0 allocs/op
// Benchmark_Debugf/G3-16 	16123359	        72.93 ns/op	       0 B/op	       0 allocs/op
// Benchmark_Debugf/G4-16 	16760646	        74.27 ns/op	       0 B/op	       0 allocs/op

func Benchmark_Debugf(b *testing.B) {
	gomaxprocsValues := []int{1, 2, 3, 4}

	for _, g := range gomaxprocsValues {
		b.Run(
			fmt.Sprintf("G%d", g),
			func(b *testing.B) {
				prev := runtime.GOMAXPROCS(g)
				defer runtime.GOMAXPROCS(prev)

				ingestor, errCrIngestor := bytearena.NewIngestor(
					bytearena.Size100K(),
					io.Discard,

					helpers.TernaryWithValueIn(
						[]int{1},
						g,
						nil,
						bytearena.WithCounterCoreCPU(),
					),
				)
				require.NoError(b, errCrIngestor)
				require.NotNil(b, ingestor)

				l, errCrLogger := NewLogger(
					&ParamsNewLogger{
						Ingestor:    ingestor,
						LoggerLevel: LevelDebug,

						WithFatalWriter: os.Stdout,
						WithJSON:        true,
					},

					WithTimestampRFC3339UTC(b.Context()),
				)
				require.NoError(b, errCrLogger)

				ctx, cancel := context.WithCancel(context.Background())
				chIngestionEnd := ingestor.StartIngestion(ctx)

				// warmup
				time.Sleep(10 * time.Millisecond)
				runtime.GC()

				b.SetParallelism(1)
				b.ReportAllocs()
				b.ResetTimer()

				b.RunParallel(
					func(pb *testing.PB) {
						for pb.Next() {
							l.Debugf("x=%d", 1)
						}
					},
				)

				cancel()
				<-chIngestionEnd
			},
		)
	}
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// Benchmark_Debugw/G1-16 	31837712	        36.43 ns/op	       0 B/op	       0 allocs/op
// Benchmark_Debugw/G2-16 	15062788	        79.61 ns/op	       0 B/op	       0 allocs/op
// Benchmark_Debugw/G3-16 	17563960	        69.25 ns/op	       0 B/op	       0 allocs/op
// Benchmark_Debugw/G4-16 	17123467	        71.49 ns/op	       0 B/op	       0 allocs/op

func Benchmark_Debugw(b *testing.B) {
	gomaxprocsValues := []int{1, 2, 3, 4}

	for _, g := range gomaxprocsValues {
		b.Run(
			fmt.Sprintf("G%d", g),
			func(b *testing.B) {
				prev := runtime.GOMAXPROCS(g)
				defer runtime.GOMAXPROCS(prev)

				ingestor, errCrIngestor := bytearena.NewIngestor(
					bytearena.Size100K(),
					io.Discard,

					helpers.TernaryWithValueIn(
						[]int{1},
						g,
						nil,
						bytearena.WithCounterCoreCPU(),
					),
				)
				require.NoError(b, errCrIngestor)
				require.NotNil(b, ingestor)

				l, errCrLogger := NewLogger(
					&ParamsNewLogger{
						Ingestor:    ingestor,
						LoggerLevel: LevelDebug,

						WithFatalWriter: io.Discard,
						WithJSON:        true,
					},

					WithTimestampRFC3339UTC(b.Context()),
				)
				require.NoError(b, errCrLogger)

				ctx, cancel := context.WithCancel(context.Background())
				chIngestionEnd := ingestor.StartIngestion(ctx)

				// warmup
				time.Sleep(10 * time.Millisecond)
				runtime.GC()

				b.ReportAllocs()
				b.ResetTimer()
				b.SetParallelism(1)

				b.RunParallel(
					func(pb *testing.PB) {
						for pb.Next() {
							l.Debugw("1some message", "some key", "some value")
						}
					},
				)

				cancel()
				<-chIngestionEnd
			},
		)
	}
}
