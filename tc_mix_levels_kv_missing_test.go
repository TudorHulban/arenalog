package arenalog

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/arenalog/query"
	"github.com/tudorhulban/bytearena"
)

// test produces
// {"ts":"2026-05-14T12:02:48+03:00","level":"TRACE","caller":"/mnt/tmpfs.ramdisk/log/methods_07_print.go","line":68,"msg":"created logger, level TRACE"}
// {"ts":"2026-05-14T12:02:48+03:00","level":"TRACE","caller":"/mnt/tmpfs.ramdisk/log/methods_00_trace.go","line":49,"msg":"LOG_ERR(odd_args): trc","key":"(MISSING)"}
// {"ts":"2026-05-14T12:02:48+03:00","level":"DEBUG","caller":"/mnt/tmpfs.ramdisk/log/methods_01_debug.go","line":49,"msg":"LOG_ERR(odd_args): dbg","key":"(MISSING)"}
// {"ts":"2026-05-14T12:02:48+03:00","level":"INFO","caller":"/mnt/tmpfs.ramdisk/log/methods_02_info.go","line":49,"msg":"LOG_ERR(odd_args): inf","key":"(MISSING)"}
// {"ts":"2026-05-14T12:02:48+03:00","level":"WARN","caller":"/mnt/tmpfs.ramdisk/log/methods_03_warn.go","line":49,"msg":"LOG_ERR(odd_args): wrn","key":"(MISSING)"}
// {"ts":"2026-05-14T12:02:48+03:00","level":"ERROR","caller":"/mnt/tmpfs.ramdisk/log/methods_04_error.go","line":49,"msg":"LOG_ERR(odd_args): err","key":"(MISSING)"}
// {"ts":"2026-05-14T12:02:48+03:00","level":"TRACE","caller":"/mnt/tmpfs.ramdisk/log/methods_07_print.go","line":83,"msg":"LOG_ERR(odd_args): prt","key":"(MISSING)"}

func TestKVMissing_JSON_Mix(t *testing.T) {
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

			WithFatalWriter: os.Stdout,
			WithCaller:      true,
			WithColor:       false,
			WithJSON:        true,
		},

		WithTimestampRFC3339UTC(t.Context()),
	)
	require.NoError(t, errCrLogger)

	l.Tracew("trc", "key")
	l.Debugw("dbg", "key")
	l.Infow("inf", "key")
	l.Warnw("wrn", "key")
	l.Errorw("err", "key")
	l.Printw("prt", "key")

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	cancel()
	<-chIngestionEnd

	out := writer.String()
	require.NotEmpty(t, out)

	logSet, errCr := query.NewLogset(out)
	require.NoError(t, errCr)

	require.Len(t, logSet.WithTimestamp(), 7)
	require.NoError(t, logSet.HasKey("caller", 7))
	require.NoError(t, logSet.HasKeyWithValue("key", "(MISSING)", 6))
	require.NoError(t, logSet.HasKeyWithValueLike("msg", "odd_args", 6))
	require.NoError(t, logSet.HasKeyWithValue("level", "TRACE", 3))
	require.NoError(t, logSet.HasKeyWithValue("level", "DEBUG", 1))
	require.NoError(t, logSet.HasKeyWithValue("level", "INFO", 1))
	require.NoError(t, logSet.HasKeyWithValue("level", "WARN", 1))
	require.NoError(t, logSet.HasKeyWithValue("level", "ERROR", 1))
}

func TestKVMissing_Raw_Mix(t *testing.T) {
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

			WithFatalWriter: os.Stdout,
			WithCaller:      true,
			WithColor:       false,
			WithJSON:        false,
		},

		WithTimestampRFC3339UTC(t.Context()),
	)
	require.NoError(t, errCrLogger)

	l.Tracew("trc", "key")
	l.Debugw("dbg", "key")
	l.Infow("inf", "key")
	l.Warnw("wrn", "key")
	l.Errorw("err", "key")
	l.Printw("prt", "key")

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	cancel()
	<-chIngestionEnd

	out := writer.String()
	require.NotEmpty(t, out)

	logSet, errCr := query.NewLogset(out)
	require.NoError(t, errCr)

	require.NoError(t, logSet.ContainsLike(7, "Line"))
	require.NoError(t, logSet.Contains(6, "key=(MISSING)"))
	require.NoError(t, logSet.ContainsEach(1, "trc", "dbg", "inf", "wrn", "err", "prt"))
}
