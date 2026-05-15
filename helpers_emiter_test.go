package arenalog

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena"
)

func TestCreateEmitData(t *testing.T) {
	var totals uint32 = 10

	state, errCrState := createEmitData(
		totals,
		map[Level]*uint32{},
		new(1),
	)
	require.NoError(t, errCrState)
	require.NotNil(t, state)

	require.Len(t, state.dictionary, int(totals))
}

func TestEmitData(t *testing.T) {
	var totals uint32 = 300

	data, errCrData := createEmitData(
		totals,
		map[Level]*uint32{},
		new(0),
	)
	require.NoError(t, errCrData)
	require.NotNil(t, data)

	safeWriter := os.Stdout // io.Discard

	writer := newTrackingWriter(
		safeWriter,
		data,
	)

	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		writer,
	)
	require.NoError(t, errCrIngestor)
	require.NotNil(t, ingestor)

	l, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelDebug,

			WithFatalWriter: safeWriter,
			WithJSON:        true,
		},

		WithTimestampRFC3339UTC(t.Context()),
	)
	require.NoError(t, errCrLogger)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	time.Sleep(10 * time.Millisecond) // warm up

	processErrors := emitData(data, l)
	require.Empty(t,
		processErrors,

		"number errors: %d",
		len(processErrors),
	)

	cancel()
	<-chIngestionEnd

	require.Zero(t,
		writer.UnreceivedCount(),
	)
}
