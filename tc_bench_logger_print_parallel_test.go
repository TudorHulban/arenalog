package arenalog

import (
	"context"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/bytearena/helpers"
)

// cpu: AMD Ryzen 5 5600U with Radeon Graphics
// BenchmarkLogger_Parallel_PrintRaw-12    	41441428	        28.67 ns/op	       0 B/op	       0 allocs/op

func BenchmarkLogger_Parallel_PrintRaw(b *testing.B) {
	prev := runtime.GOMAXPROCS(1)
	defer runtime.GOMAXPROCS(prev)

	writer := helpers.CountWriterNoBuffer{}

	ingestor, errCrIngestor := bytearena.NewIngestor(bytearena.Size100K(), &writer)
	require.NoError(b, errCrIngestor)
	require.NotNil(b, ingestor)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	time.Sleep(10 * time.Millisecond) // warmup

	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelInfo,

			WithFatalWriter: os.Stdout,
		},

		WithTimestampStandardLocal(b.Context()),
	)
	require.NoError(b, errCrLogger)
	require.NotNil(b, logger)

	runtime.GC()

	b.SetParallelism(16)
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(
		func(pb *testing.PB) {
			for pb.Next() {
				logger.PrintRaw([]byte(_Payload))
			}
		},
	)

	cancel()
	<-chIngestionEnd

	require.NotZero(b,
		writer.TotalBytesWritten.Load(),
	)
}

// go test -run '^$' -bench '^BenchmarkLogger_Parallel_Printf$' -benchmem -race

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkLogger_Parallel_Printf/1.standard_timestamp-16         	21054336	        59.81 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Printf/2.yyyy-month_timestamp-16       	19240754	        59.71 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Printf/3.nano_timestamp-16             	20278218	        60.05 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Printf/4.nano_timestamp_-_json-16      	20240864	        60.27 ns/op	       8 B/op	       1 allocs/op

func BenchmarkLogger_Parallel_Printf(b *testing.B) {
	runtime.GOMAXPROCS(1)

	tests := []struct {
		timestampOption Option
		description     string
		withJSON        bool
	}{
		{
			description:     "1.standard timestamp",
			timestampOption: WithTimestampStandardLocal(b.Context()),
		},
		{
			description:     "2.yyyy-month timestamp",
			timestampOption: WithTimestampYYYYMonthLocal(b.Context()),
		},
		{
			description: "3.nano timestamp",
		},
		{
			description: "4.nano timestamp - json",
			withJSON:    true,
		},
	}

	for _, tc := range tests {
		b.Run(
			tc.description,
			func(b *testing.B) {
				ingestor, errCrIngestor := bytearena.NewIngestor(
					bytearena.Size100K(),
					&helpers.NoopWriter{},
				)
				require.NoError(b, errCrIngestor)
				require.NotNil(b, ingestor)

				ctx, cancel := context.WithCancel(context.Background())
				chIngestionEnd := ingestor.StartIngestion(ctx)

				time.Sleep(10 * time.Millisecond) // warmup

				logger, errCrLogger := NewLogger(
					&ParamsNewLogger{
						Ingestor:    ingestor,
						LoggerLevel: LevelInfo,

						WithFatalWriter: os.Stdout,
						WithJSON:        tc.withJSON,
					},

					tc.timestampOption,
				)
				require.NoError(b, errCrLogger)
				require.NotNil(b, logger)

				runtime.GC()

				b.SetParallelism(16)
				b.ReportAllocs()
				b.ResetTimer()

				b.RunParallel(
					func(pb *testing.PB) {
						i := 0

						for pb.Next() {
							logger.Printf(
								`{"level":"info","msg":"user login","user_id":%d}`,
								i,
							)

							i++
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
// BenchmarkLogger_Parallel_Printw/1._standard_timestamp-16         	21100912	        59.47 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Printw/2._yyyy-month_timestamp-16       	20666488	        60.01 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Printw/3._nano_timestamp-16             	20684240	        59.53 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Printw/4._nano_timestamp_-_json-16      	20210739	        60.24 ns/op	       8 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Printw/5._nano_-_json,_caller-16        	18659180	        65.32 ns/op	       8 B/op	       0 allocs/op

func BenchmarkLogger_Parallel_Printw(b *testing.B) {
	runtime.GOMAXPROCS(1) // TODO: add multiple values

	tests := []struct {
		timestampOption Option
		description     string

		withJSON   bool
		withCaller bool
	}{
		{
			description:     "1. standard timestamp",
			timestampOption: WithTimestampStandardLocal(b.Context()),
		},
		{
			description:     "2. yyyy-month timestamp",
			timestampOption: WithTimestampYYYYMonthLocal(b.Context()),
		},
		{
			description: "3. nano timestamp",
		},
		{
			description: "4. nano timestamp - json",
			withJSON:    true,
		},
		{
			description: "5. nano - json, caller",
			withJSON:    true,
			withCaller:  true,
		},
	}

	for _, tc := range tests {
		b.Run(
			tc.description,
			func(b *testing.B) {
				ingestor, errCrIngestor := bytearena.NewIngestor(
					bytearena.Size100K(),
					&helpers.NoopWriter{},
				)
				require.NoError(b, errCrIngestor)
				require.NotNil(b, ingestor)

				ctx, cancel := context.WithCancel(context.Background())
				chIngestionEnd := ingestor.StartIngestion(ctx)

				time.Sleep(10 * time.Millisecond) // warmup

				logger, errCrLogger := NewLogger(
					&ParamsNewLogger{
						Ingestor:    ingestor,
						LoggerLevel: LevelInfo,

						WithFatalWriter: os.Stdout,
						WithJSON:        tc.withJSON,
						WithCaller:      tc.withCaller,
					},

					tc.timestampOption,
				)
				require.NoError(b, errCrLogger)
				require.NotNil(b, logger)

				runtime.GC()

				b.SetParallelism(16)
				b.ReportAllocs()
				b.ResetTimer()

				b.RunParallel(
					func(pb *testing.PB) {
						i := 0

						for pb.Next() {
							logger.Printf(
								`{"level":"info","msg":"user login","user_id":%d}`,
								i,
							)
							i++
						}
					},
				)

				cancel()
				<-chIngestionEnd
			},
		)
	}
}
