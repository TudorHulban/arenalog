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
)

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// Benchmark_Debug/G1-16 	23994546	        45.81 ns/op	       7 B/op	       0 allocs/op
// Benchmark_Debug/G2-16 	10402052	       111.5 ns/op	      21 B/op	       0 allocs/op
// Benchmark_Debug/G3-16 	11623293	       101.4 ns/op	      20 B/op	       0 allocs/op
// Benchmark_Debug/G4-16 	13196308	        89.85 ns/op	      18 B/op	       0 allocs/op

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
// Benchmark_Debugf/G1-16 	23244255	        51.84 ns/op	       7 B/op	       0 allocs/op
// Benchmark_Debugf/G2-16 	24547304	        47.95 ns/op	       8 B/op	       0 allocs/op
// Benchmark_Debugf/G3-16 	12575702	        98.49 ns/op	      17 B/op	       0 allocs/op
// Benchmark_Debugf/G4-16 	13308579	        89.87 ns/op	      16 B/op	       0 allocs/op

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
// Benchmark_Debugw/G1-16 	26162787	        45.03 ns/op	       5 B/op	       0 allocs/op
// Benchmark_Debugw/G2-16 	11207133	       109.2 ns/op	      15 B/op	       0 allocs/op
// Benchmark_Debugw/G3-16 	12244240	        98.87 ns/op	      15 B/op	       0 allocs/op
// Benchmark_Debugw/G4-16 	13454937	        87.53 ns/op	      14 B/op	       0 allocs/op

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
