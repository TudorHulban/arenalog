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

// go test -run '^$' -bench '^BenchmarkArenalog_Msg_OneField$' -benchmem -memprofile=mem.prof -cpuprofile=cpu.prof
// go test -bench="^BenchmarkArenalog_Msg_OneField$" -run="^$" -count=1

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkArenalog_OneField/G1-16 	18887769	        61.41 ns/op	       0 B/op	       0 allocs/op
// BenchmarkArenalog_OneField/G2-16 	12352900	        96.42 ns/op	       5 B/op	       0 allocs/op
// BenchmarkArenalog_OneField/G3-16 	12475086	        95.51 ns/op	       5 B/op	       0 allocs/op
// BenchmarkArenalog_OneField/G4-16 	12511933	        96.82 ns/op	       5 B/op	       0 allocs/op

func BenchmarkArenalog_Msg_OneField(b *testing.B) {
	gomaxprocsValues := []int{1, 2, 3, 4}
	writer := helpers.CountWriterNoBuffer{}

	for _, g := range gomaxprocsValues {
		b.Run(
			fmt.Sprintf("G%d", g),
			func(b *testing.B) {
				prev := runtime.GOMAXPROCS(g)
				defer runtime.GOMAXPROCS(prev)

				defer writer.Reset()

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

// go test -run '^$' -bench '^BenchmarkArenalog_Msgs_OneField$' -benchmem -memprofile=mem.prof -cpuprofile=cpu.prof
// go tool pprof -alloc_objects arenalog.test mem.prof
// go test -bench="^BenchmarkArenalog_Msgs_OneField$" -run="^$" -count=1

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkArenalog_Msgs_OneField/G1-16    5491644               216.6 ns/op           240 B/op          1 allocs/op
// BenchmarkArenalog_Msgs_OneField/G2-16    5920932               196.5 ns/op           240 B/op          1 allocs/op
// BenchmarkArenalog_Msgs_OneField/G3-16    5952040               200.5 ns/op           240 B/op          1 allocs/op
// BenchmarkArenalog_Msgs_OneField/G4-16    5869352               202.9 ns/op           240 B/op          1 allocs/op

func BenchmarkArenalog_Msgs_OneField(b *testing.B) {
	gomaxprocsValues := []int{1, 2, 3, 4}
	writer := helpers.CountWriterNoBuffer{}

	for _, g := range gomaxprocsValues {
		b.Run(
			fmt.Sprintf("G%d", g),
			func(b *testing.B) {
				prev := runtime.GOMAXPROCS(g)
				defer runtime.GOMAXPROCS(prev)

				defer writer.Reset()

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
					entry.Info().Msgs("benchmark test")
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
// go tool pprof -alloc_objects mem.prof

// cpu: AMD Ryzen 5 5600U with Radeon Graphics
// BenchmarkContext_NoJSON_MultipleFields-12    	11057649	       110.6 ns/op	       4 B/op	       0 allocs/op

func BenchmarkContext_NoJSON_MultipleFields(b *testing.B) {
	var writer helpers.CountWriterNoBuffer

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

	require.NotZero(b,
		writer.TotalBytesWritten.Load(),
	)
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkContext_WithJSON_MultipleFields-16    	 9249074	       132.0 ns/op	       6 B/op	       0 allocs/op

func BenchmarkContext_WithJSON_MultipleFields(b *testing.B) {
	var writer helpers.CountWriterNoBuffer

	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		&writer,
	)
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

	require.NotZero(b,
		writer.TotalBytesWritten.Load(),
	)
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkArenalog_MultipleFields_Parallel/gomaxprocs=1-16         	15310951	        76.55 ns/op	       0 B/op	       0 allocs/op
// BenchmarkArenalog_MultipleFields_Parallel/gomaxprocs=2-16         	17851222	        66.67 ns/op	       0 B/op	       0 allocs/op
// BenchmarkArenalog_MultipleFields_Parallel/gomaxprocs=3-16         	16263550	        74.16 ns/op	       0 B/op	       0 allocs/op
// BenchmarkArenalog_MultipleFields_Parallel/gomaxprocs=4-16         	15345500	        76.67 ns/op	       0 B/op	       0 allocs/op
// BenchmarkArenalog_MultipleFields_Parallel/gomaxprocs=8-16         	16770168	        73.00 ns/op	       0 B/op	       0 allocs/op

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

		cancel()
		<-chIngestionEnd

		require.NotZero(b,
			writer.TotalBytesWritten.Load(),
			"writer must record bytes",
		)

		writer.Reset()
	}
}
