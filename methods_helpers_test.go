package arenalog

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/bytearena/helpers"
)

func Test_GetLogLevel(t *testing.T) {
	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		&helpers.NoopWriter{},
	)
	require.NoError(t, errCrIngestor)
	require.NotNil(t, ingestor)

	l, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelDebug,

			WithFatalWriter: os.Stdout,
		},
	)
	require.NoError(t, errCrLogger)
	require.NotNil(t, l)

	require.EqualValues(t,
		LevelDebug,
		l.GetLogLevel(),
	)
}

func Test_SetLogLevel(t *testing.T) {
	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		&helpers.NoopWriter{},
	)
	require.NoError(t, errCrIngestor)
	require.NotNil(t, ingestor)

	l, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor: ingestor,

			WithFatalWriter: os.Stdout,
		},
	)
	require.NoError(t, errCrLogger)
	require.NotNil(t, l)

	l.SetLogLevel(LevelInfo)

	require.EqualValues(t,
		LevelInfo,
		l.GetLogLevel(),
	)
}
