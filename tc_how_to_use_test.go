package arenalog_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/arenalog"
	"github.com/tudorhulban/bytearena"
)

// test produces
// {"ts":"2026-05-19T13:56:10.288Z","level":"INFO","msg":"created logger, level INFO"}
// {"ts":"2026-05-19T13:56:10.288Z","level":"INFO","msg":"logger ready"}
// {"ts":"2026-05-19T13:56:10.288Z","level":"INFO","service":"auth","area":"some area","msg":"benchmark test"}

func TestArenalog_HowToUse(t *testing.T) {
	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		os.Stdout,
	)
	require.NoError(t, errCrIngestor)
	require.NotNil(t, ingestor)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	logger, errCrLogger := arenalog.NewLogger(
		&arenalog.ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: arenalog.LevelInfo,

			WithFatalWriter: os.Stdout,
			WithJSON:        true,
		},

		arenalog.WithTimestampRFC3339UTC(t.Context()),
	)
	require.NoError(t, errCrLogger)

	logger.Info("logger ready")

	logContext := arenalog.NewLogContext(logger).
		WithRoot("service", "auth")

	entry := logContext.WithString("area", "some area")
	entry.Info().Msg("benchmark test")

	cancel()
	<-chIngestionEnd
}
