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
// {"level":"INFO","msg":"created logger, level INFO"}
// {"level":"INFO","area":"some area","msg":"benchmark test"}

func TestEntryMsg_NoTimestamp(t *testing.T) {
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

	logSet, errParse := query.NewLogset(bufLogs.String())
	require.NoError(t, errParse)

	require.Len(t,
		logSet,
		2,

		logSet,
	)

	require.NoError(t, logSet.HasKey("level", 2))
	require.NoError(t, logSet.HasKeyWithValue("level", "INFO", 2))
}
