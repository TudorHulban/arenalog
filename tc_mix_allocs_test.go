package arenalog

import (
	"context"
	"io"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena"
)

// go test -run='^TestAllocations$' -memprofile=mem.out
// go tool pprof -alloc_objects mem.out

func TestAllocations(t *testing.T) {
	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		io.Discard,
	)
	require.NoError(t, errCrIngestor)
	require.NotNil(t, ingestor)

	l, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelTrace,

			WithFatalWriter: os.Stdout,
			WithJSON:        true,
		},

		WithTimestampRFC3339UTC(t.Context()),
	)
	require.NoError(t, errCrLogger)
	require.NotNil(t, l)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	defer func() {
		cancel()
		<-chIngestionEnd
	}()

	// warmup
	time.Sleep(10 * time.Millisecond)
	runtime.GC()

	// one allocation from time cache

	require.EqualValues(t,
		1,
		testing.AllocsPerRun(
			100,
			func() {
				l.Trace("1")
			},
		),

		"trace",
	)

	require.EqualValues(t,
		1,
		testing.AllocsPerRun(
			100,
			func() {
				l.Debug("1")
			},
		),

		"debug",
	)

	require.EqualValues(t,
		1,
		testing.AllocsPerRun(
			100,
			func() {
				l.Info("1")
			},
		),

		"info",
	)

	require.EqualValues(t,
		1,
		testing.AllocsPerRun(
			100,
			func() {
				l.Warn("1")
			},
		),

		"warn",
	)

	require.EqualValues(t,
		1,
		testing.AllocsPerRun(
			100,
			func() {
				l.Error("1")
			},
		),

		"error",
	)
}
