package arenalog_test

import (
	"context"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/arenalog"
	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/bytearena/helpers"
)

// test produces
// {"ts":"2026-05-19T15:03:19.745Z","level":"INFO","msg":"created logger, level INFO"}
// {"ts":"2026-05-19T15:03:19.745Z","level":"INFO","msg":"logger ready"}
// {"ts":"2026-05-19T15:03:19.745Z","level":"INFO","service":"auth","area":"some area","msg":"benchmark test"}
// {"ts":"2026-05-19T15:03:19.745Z","level":"INFO","service":"auth","other key":"something else","msg":"other test"}

func TestArenalog_HowToUse(t *testing.T) {
	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		os.Stdout,

		helpers.TernaryWithValueIn(
			[]int{1},
			runtime.NumCPU(),
			nil,
			bytearena.WithCounterCoreCPU(),
		),
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

	entry1 := logContext.WithString("area", "some area")
	entry1.Info().Msg("benchmark test")

	entry2 := logContext.WithString("other key", "something else")
	entry2.Info().Msg("other test")

	cancel()
	<-chIngestionEnd
}
