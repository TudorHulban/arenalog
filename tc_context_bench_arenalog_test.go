package arenalog

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/arenalog/timestamp"
	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/bytearena/helpers"
)

// test produces
// {"ts":"2026-05-11T09:33:49.624Z","level":"INFO","msg":"created logger, level INFO"}
// {"ts":"2026-05-11T09:33:49.624Z","level":"INFO","area":"some area","msg":"benchmark test"}

func TestArenalog_OneField(t *testing.T) {
	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		os.Stdout,
	)
	require.NoError(t, errCrIngestor)
	require.NotNil(t, ingestor)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelInfo,

			WithFatalWriter: os.Stdout,
			WithJSON:        true,
		},

		WithTimestampRFC3339UTC(t.Context()),
	)
	require.NoError(t, errCrLogger)

	logContext := NewLogContext(logger)

	entry := logContext.WithString("area", "some area")
	entry.Info().Msg("benchmark test")

	cancel()
	<-chIngestionEnd
}

// go test -run '^$' -bench '^BenchmarkArenalog_OneField$' -benchmem -memprofile=mem.prof -cpuprofile=cpu.prof

// cpu: AMD Ryzen 5 5600U with Radeon Graphics
// BenchmarkArenalog_OneField/G1-12 	18744702	        64.20 ns/op	       0 B/op	       0 allocs/op
// BenchmarkArenalog_OneField/G2-12 	11996622	       100.5 ns/op	       6 B/op	       0 allocs/op
// BenchmarkArenalog_OneField/G3-12 	11811764	       100.9 ns/op	       6 B/op	       0 allocs/op
// BenchmarkArenalog_OneField/G4-12 	11488641	       101.7 ns/op	       6 B/op	       0 allocs/op

func BenchmarkArenalog_OneField(b *testing.B) {
	gomaxprocsValues := []int{1, 2, 3, 4}
	writer := helpers.CountWriterNoBuffer{}

	for _, g := range gomaxprocsValues {
		b.Run(
			fmt.Sprintf("G%d", g),
			func(b *testing.B) {
				prev := runtime.GOMAXPROCS(g)
				defer runtime.GOMAXPROCS(prev)

				ingestor, errCrIngestor := bytearena.NewIngestor(
					bytearena.Size100K(),
					&writer,
				)
				require.NoError(b, errCrIngestor)
				require.NotNil(b, ingestor)

				ctx, cancel := context.WithCancel(context.Background())
				chIngestionEnd := ingestor.StartIngestion(ctx)

				logger, errCrLogger := NewLogger(
					&ParamsNewLogger{
						Ingestor:    ingestor,
						LoggerLevel: LevelDebug,

						WithFatalWriter: os.Stdout,
						WithJSON:        true,
					},

					WithTimestampRFC3339UTC(b.Context()),
				)
				require.NoError(b, errCrLogger)

				logContext := NewLogContext(logger).
					WithRoot("service", "auth")

				runtime.GC()

				// warm up the pool
				for i := 0; i < runtime.GOMAXPROCS(0)*4; i++ {
					e, _ := entryPool.Get().(*Entry) //nolint:revive
					entryPool.Put(e)
				}

				var warmupBuffer []byte
				timestamp.TimestampRFC3339UTC(warmupBuffer)

				b.ReportAllocs()
				b.ResetTimer()

				for b.Loop() {
					entry := logContext.WithString("area", "some area")
					entry.Info().Msg("benchmark test")
				}

				cancel()
				<-chIngestionEnd

				require.NotZero(b,
					writer.TotalBytesWritten.Load(),
				)
			},
		)
	}
}

// test produces
// {"ts":"2026-05-04T14:29:35Z","level":"INFO","msg":"created logger, level INFO"}
// {"ts":"2026-05-04T14:29:35Z","level":"INFO","service":"auth","req_id":12345,"cache_hit":true,"area":"some area","user":"tudor","attempt":1,"some float":1.113699999999,"success":true,"message":"benchmark test"}

func TestArenalog_MultipleFields(t *testing.T) {
	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		os.Stdout,
	)
	require.NoError(t, errCrIngestor)
	require.NotNil(t, ingestor)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelInfo,

			WithFatalWriter: os.Stdout,
			// WithJSON:        true,
		},

		WithTimestampRFC3339UTC(t.Context()),
	)
	require.NoError(t, errCrLogger)

	logContext := NewLogContext(logger).
		WithRoot("service", "auth").
		SetInt("req_id", 12345).
		SetBool("cache_hit", true)

	entry := logContext.
		WithString("area", "some area").
		Info().
		WithString("user", "tudor").
		WithInt("attempt", 1).
		WithFloat("some float", 1.1137).
		WithBool("success", true)

	entry.Msg("benchmark test")

	cancel()
	<-chIngestionEnd
}

// go test -run '^$' -bench '^BenchmarkContext_NoJSON_MultipleFields$' -benchmem -memprofile=mem.prof -cpuprofile=cpu.prof
// go tool pprof -alloc_objects mem.out

// cpu: AMD Ryzen 5 5600U with Radeon Graphics
// BenchmarkContext_NoJSON_MultipleFields-12    	11057649	       110.6 ns/op	       4 B/op	       0 allocs/op

func BenchmarkContext_NoJSON_MultipleFields(b *testing.B) {
	var writer helpers.NoopWriter

	ingestor, errCrIngestor := bytearena.NewIngestor(bytearena.Size100K(), &writer)
	require.NoError(b, errCrIngestor)
	require.NotNil(b, ingestor)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	defer func() {
		cancel()
		<-chIngestionEnd
	}()

	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelDebug,

			WithFatalWriter: os.Stdout,
		},

		WithTimestampRFC3339UTC(b.Context()),
	)
	require.NoError(b, errCrLogger)

	logContext := NewLogContext(logger).
		WithRoot("service", "auth").
		SetInt("req_id", 12345).
		SetBool("cache_hit", true)

	runtime.GC()

	// warm up the pool
	for i := 0; i < runtime.GOMAXPROCS(0)*4; i++ {
		e, _ := entryPool.Get().(*Entry) //nolint:revive
		entryPool.Put(e)
	}

	var warmupBuffer []byte
	timestamp.TimestampRFC3339UTC(warmupBuffer)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		// 1. Create request with several attributes
		entry := logContext.
			WithString("area", "some area").
			Info().
			WithString("user", "tudor").
			WithInt("attempt", int64(i)).
			WithFloat("some float", 1.1137).
			WithBool("success", true)

		// 2. Print
		entry.Msg("benchmark test")
	}
}

// cpu: AMD Ryzen 5 5600U with Radeon Graphics
// BenchmarkContext_WithJSON_MultipleFields-12    	 8336259	       146.6 ns/op	       6 B/op	       0 allocs/op

func BenchmarkContext_WithJSON_MultipleFields(b *testing.B) {
	var writer helpers.NoopWriter

	ingestor, errCrIngestor := bytearena.NewIngestor(bytearena.Size100K(), &writer)
	require.NoError(b, errCrIngestor)
	require.NotNil(b, ingestor)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	defer func() {
		cancel()
		<-chIngestionEnd
	}()

	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelDebug,

			WithFatalWriter: os.Stdout,
			WithJSON:        true,
		},

		WithTimestampRFC3339UTC(b.Context()),
	)
	require.NoError(b, errCrLogger)

	logContext := NewLogContext(logger).
		WithRoot("service", "auth").
		SetInt("req_id", 12345).
		SetBool("cache_hit", true)

	runtime.GC()

	// warm up the pool
	for i := 0; i < runtime.GOMAXPROCS(0)*4; i++ {
		e, _ := entryPool.Get().(*Entry) //nolint:revive
		entryPool.Put(e)
	}

	var warmupBuffer []byte
	timestamp.TimestampRFC3339UTC(warmupBuffer)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		// 1. Create request with several attributes
		entry := logContext.
			WithString("area", "some area").
			Info().
			WithString("user", "tudor").
			WithInt("attempt", int64(i)).
			WithFloat("some float", 1.1137).
			WithBool("success", true)

		// 2. Print
		entry.Msg("benchmark test")
	}
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkArenalog_MultipleFields_Parallel/gomaxprocs=1-16         	16182062	        75.53 ns/op	       0 B/op	       0 allocs/op
// BenchmarkArenalog_MultipleFields_Parallel/gomaxprocs=2-16         	18008524	        72.07 ns/op	       0 B/op	       0 allocs/op
// BenchmarkArenalog_MultipleFields_Parallel/gomaxprocs=3-16         	17262482	        70.27 ns/op	       0 B/op	       0 allocs/op
// BenchmarkArenalog_MultipleFields_Parallel/gomaxprocs=4-16         	15265447	        76.54 ns/op	       0 B/op	       0 allocs/op
// BenchmarkArenalog_MultipleFields_Parallel/gomaxprocs=8-16         	17359221	        73.10 ns/op	       0 B/op	       0 allocs/op

func BenchmarkArenalog_MultipleFields_Parallel(b *testing.B) {
	gomaxprocsValues := []int{1, 2, 3, 4, 8}
	writer := helpers.CountWriterNoBuffer{}

	b.SetParallelism(1)

	for _, g := range gomaxprocsValues {
		ingestor, errCrIngestor := bytearena.NewIngestor(
			bytearena.Size100K(),
			&writer,

			helpers.TernaryWithValueIn(
				[]int{1},
				g,
				nil,
				bytearena.WithCounterCoreCPU(),
			),
		)
		require.NoError(b, errCrIngestor)
		require.NotNil(b, ingestor)

		ctx, cancel := context.WithCancel(context.Background())
		chIngestionEnd := ingestor.StartIngestion(ctx)

		// keep ingestion running for the whole benchmark
		defer func() { //nolint:revive
			cancel()
			<-chIngestionEnd
		}()

		logger, errCrLogger := NewLogger(
			&ParamsNewLogger{
				Ingestor:        ingestor,
				LoggerLevel:     LevelDebug,
				WithFatalWriter: os.Stdout,
				WithJSON:        true,
			},

			WithTimestampRFC3339UTC(b.Context()),
		)
		require.NoError(b, errCrLogger)

		logContext := NewLogContext(logger).
			WithRoot("service", "auth").
			SetInt("req_id", 12345).
			SetBool("cache_hit", true)

		runtime.GC()

		// warm up the pool
		for i := 0; i < runtime.GOMAXPROCS(0)*4; i++ {
			e, _ := entryPool.Get().(*Entry) //nolint:revive
			entryPool.Put(e)
		}

		var warmupBuffer []byte
		timestamp.TimestampRFC3339UTC(warmupBuffer)

		time.Sleep(10 * time.Millisecond)

		b.Run(
			fmt.Sprintf("gomaxprocs=%d", g),
			func(b *testing.B) {
				prev := runtime.GOMAXPROCS(g)
				defer runtime.GOMAXPROCS(prev)

				b.ReportAllocs()
				b.ResetTimer()

				b.RunParallel(
					func(pb *testing.PB) {
						i := int64(0)
						for pb.Next() {
							entry := logContext.
								WithString("area", "some area").
								Info().
								WithString("user", "tudor").
								WithInt("attempt", i).
								WithFloat("some float", 1.1137).
								WithBool("success", true)

							entry.Msg("benchmark test")

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
