package arenalog

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/arenalog/query"
	"github.com/tudorhulban/bytearena"
)

// test produces
// {"ts":"2026-05-12T15:51:49+03:00","level":"TRACE","msg":"created logger, level TRACE"}
// {"ts":"2026-05-12T15:51:49+03:00","level":"TRACE","msg":"xxx"}
// {"ts":"2026-05-12T15:51:49+03:00","level":"TRACE","msg":"LOG_ERR(odd_args): xxx","key":"(MISSING)"}

func TestErrorsPrint(t *testing.T) {
	var buf bytes.Buffer

	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		&buf,
	)
	require.NoError(t, errCrIngestor)
	require.NotNil(t, ingestor)

	l, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelTrace,

			WithFatalWriter: &buf,
			WithJSON:        true,
		},

		WithTimestampRFC3339UTC(t.Context()),
	)
	require.NoError(t, errCrLogger)
	require.NotNil(t, l)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	// 1. Valid call - No extra keys
	l.Printw("xxx")

	// 2. Invalid call - Odd number of keys (1 key, 0 values)
	// This should trigger the internal "panicw" log line
	l.Printw("xxx", "key")

	cancel()
	<-chIngestionEnd

	// Parse the results from the buffer
	logSet, errCr := query.NewLogset(buf.String())
	require.NoError(t, errCr)

	require.Len(t, logSet, 3)

	require.NoError(t, logSet.HasKeyWithValue("msg", "xxx", 1))
	require.NoError(t, logSet.HasKey("msg", 3))
	require.NoError(t, logSet.HasKeyWithValue("key", _Missing, 1))
	require.NoError(t, logSet.HasKey("level", 3))
}

// test produces
// {"level":"INFO","msg":"created logger, level INFO"}
// {"level":"INFO","msg":"benchmark test"}

func TestLoggerMsg_NoTimestamp(t *testing.T) {
	var bufLogs, bufFatal bytes.Buffer

	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		&bufLogs,
	)
	require.NoError(t, errCrIngestor)
	require.NotNil(t, ingestor)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelInfo,

			WithCaller:      true,
			WithFatalWriter: &bufFatal,
			WithJSON:        false,
		},
	)
	require.NoError(t, errCrLogger)

	logger.Msg("benchmark test")

	cancel()
	<-chIngestionEnd

	// --- Processing Lines ---
	require.Zero(t, bufFatal.Len())

	logSet, errParse := query.NewLogset(bufLogs.String())
	require.NoError(t, errParse)

	require.Len(t,
		logSet,
		2,

		logSet,
	)
}
