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
// BenchmarkLogContextPrint/G1-16 	42220496	        28.11 ns/op	       0 B/op	       0 allocs/op
// BenchmarkLogContextPrint/G2-16 	14031519	        86.45 ns/op	       0 B/op	       0 allocs/op
// BenchmarkLogContextPrint/G3-16 	16144488	        76.03 ns/op	       0 B/op	       0 allocs/op
// BenchmarkLogContextPrint/G4-16 	16930374	        72.70 ns/op	       0 B/op	       0 allocs/op

func BenchmarkLogContextPrint(b *testing.B) {
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

				f := NewLogContext(l).
					WithRoot("service", "auth").
					SetInt("req_id", 12345).
					SetBool("cache_hit", true).
					SetString("root ends", "here")

				// warmup
				time.Sleep(10 * time.Millisecond)
				runtime.GC()

				b.SetParallelism(1)
				b.ReportAllocs()
				b.ResetTimer()

				b.RunParallel(
					func(pb *testing.PB) {
						for pb.Next() {
							f.Print(_Payload)
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
// BenchmarkLogContextPrints/G1-16 	 5843539	       203.3 ns/op	     256 B/op	       2 allocs/op
// BenchmarkLogContextPrints/G2-16 	 9113028	       131.3 ns/op	     256 B/op	       2 allocs/op
// BenchmarkLogContextPrints/G3-16 	11232903	       105.2 ns/op	     256 B/op	       2 allocs/op
// BenchmarkLogContextPrints/G4-16 	12061215	        93.68 ns/op	     256 B/op	       2 allocs/op

func BenchmarkLogContextPrints(b *testing.B) {
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

				f := NewLogContext(l).
					WithRoot("service", "auth").
					SetInt("req_id", 12345).
					SetBool("cache_hit", true).
					SetString("root ends", "here")

				// warmup
				time.Sleep(10 * time.Millisecond)
				runtime.GC()

				b.SetParallelism(1)
				b.ReportAllocs()
				b.ResetTimer()

				b.RunParallel(
					func(pb *testing.PB) {
						for pb.Next() {
							f.Prints(_Payload)
						}
					},
				)

				cancel()
				<-chIngestionEnd
			},
		)
	}
}
