package arenalog

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/arenalog/query"
	"github.com/tudorhulban/bytearena"
)

// test produces
// {"ts":"1779980759609781478","level":"INFO","msg":"created logger, level INFO"}
// {"ts":"1779980759609784153","level":"INFO","area":"some area","msg":"benchmark test"}

func TestEntryMsg_DefaultTimestamp(t *testing.T) {
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

			WithFatalWriter: &bufFatal,
			WithJSON:        true,
		},
	)
	require.NoError(t, errCrLogger)

	logContext := NewLogContext(logger)

	entry := logContext.WithString("area", "some area")
	entry.Info().Msg("benchmark test")

	cancel()
	<-chIngestionEnd

	// --- Processing Lines ---
	require.Zero(t, bufFatal.Len())

	// fmt.Println(bufLogs.String())

	logSet, errParse := query.NewLogset(bufLogs.String())
	require.NoError(t, errParse)

	require.Len(t,
		logSet,
		2,

		logSet,
	)

	require.Len(t, logSet.WithTimestamp(), 2)
	require.NoError(t, logSet.HasKey("level", 2))
	require.NoError(t, logSet.HasKeyWithValue("level", "INFO", 2))
}

func TestEntryMsgs_DefaultTimestamp(t *testing.T) {
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

			WithFatalWriter: &bufFatal,
			WithJSON:        true,
		},
	)
	require.NoError(t, errCrLogger)

	logContext := NewLogContext(logger)

	entry := logContext.WithString("area", "some area")
	entry.Info().Msg("benchmark test")

	cancel()
	<-chIngestionEnd

	// --- Processing Lines ---
	require.Zero(t, bufFatal.Len())

	fmt.Println(bufLogs.String())

	logSet, errParse := query.NewLogset(bufLogs.String())
	require.NoError(t, errParse)

	require.Len(t,
		logSet,
		2,

		logSet,
	)

	require.Len(t, logSet.WithTimestamp(), 2)
	require.NoError(t, logSet.HasKey("level", 2))
	require.NoError(t, logSet.HasKeyWithValue("level", "INFO", 2))
}
