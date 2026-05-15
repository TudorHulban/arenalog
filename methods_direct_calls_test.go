package arenalog

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/arenalog/query"
	"github.com/tudorhulban/bytearena"
)

// test produces:
// {"ts":"2026-05-10T18:30:22+03:00","level":"INFO","caller":"/mnt/tmpfs.ramdisk/log/methods_02_info.go","line":18,"msg":"333"}
// {"ts":"2026-05-10T18:30:22+03:00","level":"ERROR","caller":"/mnt/tmpfs.ramdisk/log/methods_04_error.go","line":43,"msg":"msg-error","key1":1}
// {"ts":"2026-05-10T18:30:22+03:00","level":"TRACE","caller":"/mnt/tmpfs.ramdisk/log/methods_07_print.go","line":81,"msg":"msg-print","key1":1}
// {"ts":"2026-05-10T18:30:22+03:00","level":"TRACE","caller":"/mnt/tmpfs.ramdisk/log/methods_07_print.go","line":68,"msg":"created logger, level TRACE"}
// {"ts":"2026-05-10T18:30:22+03:00","level":"INFO","caller":"/mnt/tmpfs.ramdisk/log/methods_02_info.go","line":30,"msg":"1"}
// {"ts":"2026-05-10T18:30:22+03:00","level":"TRACE","msg":"666"}
// {"ts":"2026-05-10T18:30:22+03:00","level":"TRACE","caller":"/mnt/tmpfs.ramdisk/log/methods_00_trace.go","line":18,"msg":"111"}
// {"ts":"2026-05-10T18:30:22+03:00","level":"INFO","caller":"/mnt/tmpfs.ramdisk/log/methods_02_info.go","line":47,"msg":"msg-info","key1":1}
// {"ts":"2026-05-10T18:30:22+03:00","level":"TRACE","caller":"/mnt/tmpfs.ramdisk/log/methods_07_print.go","line":60,"msg":"777"}
// {"ts":"2026-05-10T18:30:22+03:00","level":"TRACE","caller":"/mnt/tmpfs.ramdisk/log/methods_00_trace.go","line":30,"msg":"1"}
// {"ts":"2026-05-10T18:30:22+03:00","level":"WARN","caller":"/mnt/tmpfs.ramdisk/log/methods_03_warn.go","line":18,"msg":"444"}
// {"ts":"2026-05-10T18:30:22+03:00","level":"TRACE","caller":"/mnt/tmpfs.ramdisk/log/methods_07_print.go","line":60,"msg":"777"}
// {"ts":"2026-05-10T18:30:22+03:00","level":"TRACE","caller":"/mnt/tmpfs.ramdisk/log/methods_00_trace.go","line":47,"msg":"msg-trace","key1":1}
// {"ts":"2026-05-10T18:30:22+03:00","level":"WARN","caller":"/mnt/tmpfs.ramdisk/log/methods_03_warn.go","line":30,"msg":"1"}
// RAW777
// {"ts":"2026-05-10T18:30:22+03:00","level":"DEBUG","caller":"/mnt/tmpfs.ramdisk/log/methods_01_debug.go","line":18,"msg":"222"}
// {"ts":"2026-05-10T18:30:22+03:00","level":"WARN","caller":"/mnt/tmpfs.ramdisk/log/methods_03_warn.go","line":43,"msg":"msg-warn","key1":1}
// {"ts":"2026-05-10T18:30:22+03:00","level":"TRACE","caller":"/mnt/tmpfs.ramdisk/log/methods_07_print.go","line":68,"msg":"1"}
// {"ts":"2026-05-10T18:30:22+03:00","level":"DEBUG","caller":"/mnt/tmpfs.ramdisk/log/methods_01_debug.go","line":30,"msg":"1"}
// {"ts":"2026-05-10T18:30:22+03:00","level":"ERROR","caller":"/mnt/tmpfs.ramdisk/log/methods_04_error.go","line":18,"msg":"555"}
// {"ts":"2026-05-10T18:30:22+03:00","level":"TRACE","caller":"/mnt/tmpfs.ramdisk/log/methods_07_print.go","line":68,"msg":"1"}
// {"ts":"2026-05-10T18:30:22+03:00","level":"DEBUG","caller":"/mnt/tmpfs.ramdisk/log/methods_01_debug.go","line":47,"msg":"msg-debug","key1":1}
// {"ts":"2026-05-10T18:30:22+03:00","level":"ERROR","caller":"/mnt/tmpfs.ramdisk/log/methods_04_error.go","line":30,"msg":"1"}
// {"ts":"2026-05-10T18:30:22+03:00","level":"TRACE","caller":"/mnt/tmpfs.ramdisk/log/methods_07_print.go","line":81,"msg":"msg-print","key1":1}

func TestDirectCalls(t *testing.T) {
	var writer bytes.Buffer

	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		&writer,
	)
	require.NoError(t, errCrIngestor)
	require.NotNil(t, ingestor)

	l, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelTrace,

			WithFatalWriter: &writer,
			WithCaller:      true,
			WithColor:       false,
			WithJSON:        true,
		},

		WithTimestampRFC3339UTC(t.Context()),
	)
	require.NoError(t, errCrLogger)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	const (
		value    = 1
		keyValue = 1
	)

	l.Trace("111")
	l.Tracef("%d", value)
	l.Tracew("msg-trace", "key1", keyValue)

	l.Debug("222")
	l.Debugf("%d", value)
	l.Debugw("msg-debug", "key1", keyValue)

	l.Info("333")
	l.Infof("%d", value)
	l.Infow("msg-info", "key1", keyValue)

	l.Warn("444")
	l.Warnf("%d", value)
	l.Warnw("msg-warn", "key1", keyValue)

	l.Error("555")
	l.Errorf("%d", value)
	l.Errorw("msg-error", "key1", keyValue)

	l.Msg("666")
	l.Print("777")
	l.Print("777")
	l.PrintRaw([]byte("RAW777\n")) // This line is skipped by your DSL (no '{')
	l.Printf("%d", value)
	l.Printf("%d", value)
	l.Printw("msg-print", "key1", keyValue)
	l.Printw("msg-print", "key1", keyValue)

	cancel()
	<-chIngestionEnd

	logSet, errParse := query.NewLogset(writer.String())
	require.NoError(t, errParse)
	require.NotEmpty(t, logSet)

	// --- Assertions based on updated "produces" data ---
	require.Len(t, logSet, 24)

	// Verify Timestamp
	require.EqualValues(t,
		len(logSet.WithTimestamp()),
		len(logSet)-1,
	)

	// Counts by Level
	require.NoError(t, logSet.HasKeyWithValue("level", "TRACE", 11))
	require.NoError(t, logSet.HasKeyWithValue("level", "DEBUG", 3))
	require.NoError(t, logSet.HasKeyWithValue("level", "INFO", 3))
	require.NoError(t, logSet.HasKeyWithValue("level", "WARN", 3))
	require.NoError(t, logSet.HasKeyWithValue("level", "ERROR", 3))

	// Verify Specific Messages
	require.NoError(t, logSet.HasKeyWithValue("msg", "111", 1))
	require.NoError(t, logSet.HasKeyWithValue("msg", "777", 2))
	require.NoError(t, logSet.HasKeyWithValue("msg", "1", 7)) // 5 from formatted calls + 2 from Printf
	require.NoError(t, logSet.HasKeyWithValue("msg", "msg-print", 2))

	// Verify Structured Key-Values
	// key1: 1 appears in Tracew, Debugw, Infow, Warnw, Errorw, and 2x Printw = 7 total
	require.NoError(t, logSet.HasKeyWithValue("key1", keyValue, 7))

	// Verify Complex Matches
	// Ensure Printw produced TRACE level logs with the correct message
	require.NoError(t,
		logSet.HasKeysWithValues(2,
			"level", "TRACE",
			"msg", "msg-print",
			"key1", keyValue,
		),
	)

	// Verify Caller Presence
	// Total lines is 23 (24 minus the RAW777 line which newTestEntries skips).
	// One line (msg: 666) is missing a caller in your "produces" output.
	require.NoError(t, logSet.HasKey("caller", 22))
}
