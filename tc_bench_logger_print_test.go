package arenalog

import (
	"context"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/bytearena/helpers"
)

// go test -run '^$' -bench '^BenchmarkLogger_Print$' -benchmem

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkLogger_Print-16    	31454732	        38.31 ns/op	       0 B/op	       0 allocs/op
func BenchmarkLogger_Print(b *testing.B) {
	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		&helpers.NoopWriter{},
	)
	require.NoError(b, errCrIngestor)
	require.NotNil(b, ingestor)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

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

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		logger.Print("hi", 123, "world")
	}

	cancel()
	<-chIngestionEnd
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkLogger_Printf/1.standard_timestamp-16         	12163123	        98.82 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_Printf/2.yyyy-month_timestamp-16       	11890149	        99.23 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_Printf/3.nano_timestamp-16             	 9798445	       123.3 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_Printf/4.nano_timestamp_-_json-16      	 6019629	       201.3 ns/op	      54 B/op	       2 allocs/op

func BenchmarkLogger_Printf(b *testing.B) {
	tests := []struct {
		timestampOption Option
		description     string
		withJSON        bool
		withCaller      bool
	}{
		{
			description:     "1.standard timestamp",
			timestampOption: WithTimestampRFC3339UTC(b.Context()),
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

				logger, errCrLogger := NewLogger(
					&ParamsNewLogger{
						Ingestor:    ingestor,
						LoggerLevel: LevelInfo,

						WithFatalWriter: os.Stdout,

						WithJSON:   tc.withJSON,
						WithCaller: tc.withCaller,
					},

					tc.timestampOption,
				)
				require.NoError(b, errCrLogger)
				require.NotNil(b, logger)

				runtime.GC()

				b.ReportAllocs()
				b.ResetTimer()

				for i := 0; b.Loop(); i++ {
					logger.Printf(
						`{"level":"info","msg":"user login","user_id":%d}`,
						i,
					)
				}

				cancel()
				<-chIngestionEnd
			},
		)
	}
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkLogger_Printw-16    	23154374	        51.98 ns/op	       0 B/op	       0 allocs/op

func BenchmarkLogger_Printw(b *testing.B) {
	var writer helpers.NoopWriter

	ingestor, errCrIngestor := bytearena.NewIngestor(bytearena.Size100K(), &writer)
	require.NoError(b, errCrIngestor)
	require.NotNil(b, ingestor)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

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

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		logger.Printw("hi", 123, "world")
	}

	cancel()
	<-chIngestionEnd
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkLogger_PrintRaw-16    	68537496	        17.74 ns/op	       0 B/op	       0 allocs/op

func BenchmarkLogger_PrintRaw(b *testing.B) {
	var writer helpers.NoopWriter

	ingestor, errCrIngestor := bytearena.NewIngestor(bytearena.Size100K(), &writer)
	require.NoError(b, errCrIngestor)
	require.NotNil(b, ingestor)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

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

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		logger.PrintRaw(
			[]byte(_Payload),
		)
	}

	cancel()
	<-chIngestionEnd
}
