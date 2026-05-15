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

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkLogger_Info/1.standard_timestamp-16         	27905121	        44.21 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_Info/2.yyyy-month_timestamp-16       	26776683	        44.49 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_Info/3.nano_timestamp-16             	18282640	        66.17 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_Info/4.json_-_rfc3339-16             	20330120	        60.16 ns/op	      10 B/op	       1 allocs/op
// BenchmarkLogger_Info/5.nano_timestamp_-_json-16      	15303798	        78.12 ns/op	      10 B/op	       1 allocs/op
// BenchmarkLogger_Info/6.nano_-_json_caller-16         	 4518417	       267.0 ns/op	      12 B/op	       1 allocs/op

func BenchmarkLogger_Info(b *testing.B) {
	tests := []struct {
		timestampOption Option
		description     string
		withJSON        bool
		withCaller      bool
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
			description:     "4.json - rfc3339",
			timestampOption: WithTimestampRFC3339UTC(b.Context()),
			withJSON:        true,
		},
		{
			description: "5.nano timestamp - json",
			withJSON:    true,
		},
		{
			description: "6.nano - json caller",
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
						LoggerLevel: LevelDebug,

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
					logger.Info(i)
				}

				cancel()
				<-chIngestionEnd
			},
		)
	}
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkLogger_Infof/1.standard_timestamp-16         	23847734	        51.67 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_Infof/2.yyyy-month_timestamp-16       	23209123	        51.95 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_Infof/3.nano_timestamp-16             	16714066	        72.70 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_Infof/4.json_-_rfc3339-16             	16421894	        74.62 ns/op	      10 B/op	       1 allocs/op
// BenchmarkLogger_Infof/5.nano_timestamp_-_json-16      	12550860	        96.64 ns/op	      10 B/op	       1 allocs/op
// BenchmarkLogger_Infof/6.nano_-_json_caller-16         	 4168384	       289.3 ns/op	      12 B/op	       1 allocs/op

func BenchmarkLogger_Infof(b *testing.B) {
	tests := []struct {
		timestampOption Option
		description     string
		withJSON        bool
		withCaller      bool
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
			description:     "4.json - rfc3339",
			timestampOption: WithTimestampRFC3339UTC(b.Context()),
			withJSON:        true,
		},
		{
			description: "5.nano timestamp - json",
			withJSON:    true,
		},
		{
			description: "6.nano - json caller",
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
						LoggerLevel: LevelDebug,

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
					logger.Infof(
						"%d",
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
// BenchmarkLogger_Infow/1.standard_timestamp-16         	34667320	        34.71 ns/op	       0 B/op	       0 allocs/op
// BenchmarkLogger_Infow/2.yyyy-month_timestamp-16       	34153743	        34.58 ns/op	       0 B/op	       0 allocs/op
// BenchmarkLogger_Infow/3.nano_timestamp-16             	24508124	        49.17 ns/op	       0 B/op	       0 allocs/op
// BenchmarkLogger_Infow/4.json_-_rfc3339-16             	23080941	        52.20 ns/op	       0 B/op	       0 allocs/op
// BenchmarkLogger_Infow/5.nano_timestamp_-_json-16      	19016328	        63.11 ns/op	       0 B/op	       0 allocs/op
// BenchmarkLogger_Infow/6.nano_-_json_caller-16         	 4697342	       255.8 ns/op	       0 B/op	       0 allocs/op

func BenchmarkLogger_Infow(b *testing.B) {
	tests := []struct {
		timestampOption Option
		description     string
		withJSON        bool
		withCaller      bool
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
			description:     "4.json - rfc3339",
			timestampOption: WithTimestampRFC3339UTC(b.Context()),
			withJSON:        true,
		},
		{
			description: "5.nano timestamp - json",
			withJSON:    true,
		},
		{
			description: "6.nano - json caller",
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
						LoggerLevel: LevelDebug,

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
					logger.Infow(_Payload, "key", "value")
				}

				cancel()
				<-chIngestionEnd
			},
		)
	}
}
