package arenalog_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	log "github.com/tudorhulban/arenalog"
	"github.com/tudorhulban/bytearena"
)

// {"ts":"2026-05-15T13:50:44.712Z","level":"INFO","msg":"created logger, level INFO"}
// {"ts":"2026-05-15T13:50:44.712Z","level":"INFO","area":"some area","msg":"benchmark test"}

func TestArenalog_HowToUse(t *testing.T) {
	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		os.Stdout,
	)
	require.NoError(t, errCrIngestor)
	require.NotNil(t, ingestor)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	logger, errCrLogger := log.NewLogger(
		&log.ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: log.LevelInfo,

			WithFatalWriter: os.Stdout,
			WithJSON:        true,
		},

		log.WithTimestampRFC3339UTC(t.Context()),
	)
	require.NoError(t, errCrLogger)

	logContext := log.NewLogContext(logger)

	entry := logContext.WithString("area", "some area")
	entry.Info().Msg("benchmark test")

	cancel()
	<-chIngestionEnd
}
