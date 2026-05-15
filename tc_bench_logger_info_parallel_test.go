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
// BenchmarkLogger_Parallel_Info/1.standard_timestamp-16         	21400766	        58.16 ns/op	       0 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Info/2.yyyy-month_timestamp-16       	21109534	        58.25 ns/op	       0 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Info/3.nano_timestamp-16             	21019228	        58.18 ns/op	       0 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Info/4.nano_timestamp_-_json-16      	20786533	        58.28 ns/op	       0 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Info/5.nano_-_json,_caller-16        	19577006	        62.34 ns/op	       0 B/op	       0 allocs/op

func BenchmarkLogger_Parallel_Info(b *testing.B) {
	runtime.GOMAXPROCS(1)

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
			description: "4.nano timestamp - json",
			withJSON:    true,
		},
		{
			description: "5.nano - json, caller",
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
					&helpers.CountWriterNoBuffer{},
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
						for pb.Next() {
							logger.Info(_Payload)
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
// BenchmarkLogger_Parallel_Infof/1.standard_timestamp-16         	21388173	        58.69 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Infof/2.yyyy-month_timestamp-16       	20876204	        60.12 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Infof/3.nano_timestamp-16             	20556169	        59.31 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Infof/4.nano_timestamp_-_json-16      	20568652	        60.14 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Infof/5.nano_-_json,_caller-16        	19268158	        64.02 ns/op	       7 B/op	       0 allocs/op

func BenchmarkLogger_Parallel_Infof(b *testing.B) {
	runtime.GOMAXPROCS(1)

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
			description: "4.nano timestamp - json",
			withJSON:    true,
		},
		{
			description: "5.nano - json, caller",
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
					&helpers.CountWriterNoBuffer{},
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

				b.SetParallelism(16)
				b.ReportAllocs()
				b.ResetTimer()

				b.RunParallel(
					func(pb *testing.PB) {
						i := 0

						for pb.Next() {
							logger.Infof(
								"%d",
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
// BenchmarkLogger_Parallel_Infow/1.standard_timestamp-16         	19586179	        61.92 ns/op	      32 B/op	       1 allocs/op
// BenchmarkLogger_Parallel_Infow/2.yyyy-month_timestamp-16       	19557171	        61.68 ns/op	      32 B/op	       1 allocs/op
// BenchmarkLogger_Parallel_Infow/3.nano_timestamp-16             	19596718	        62.57 ns/op	      32 B/op	       1 allocs/op
// BenchmarkLogger_Parallel_Infow/4.nano_timestamp_-_json-16      	18916557	        63.68 ns/op	      32 B/op	       1 allocs/op
// BenchmarkLogger_Parallel_Infow/5.nano_-_json,_caller-16        	16438960	        71.58 ns/op	      32 B/op	       1 allocs/op

func BenchmarkLogger_Parallel_Infow(b *testing.B) {
	runtime.GOMAXPROCS(1)

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
			description: "4.nano timestamp - json",
			withJSON:    true,
		},
		{
			description: "5.nano - json, caller",
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
					&helpers.CountWriterNoBuffer{},
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

				b.SetParallelism(16)
				b.ReportAllocs()
				b.ResetTimer()

				b.RunParallel(
					func(pb *testing.PB) {
						for pb.Next() {
							logger.Infow("key", _Payload)
						}
					},
				)

				cancel()
				<-chIngestionEnd
			},
		)
	}
}
