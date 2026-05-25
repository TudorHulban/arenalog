package arenalog

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/arenalog/query"
	"github.com/tudorhulban/bytearena"
)

// test produces
// {"ts":"2026-05-12T13:34:48+03:00","level":"TRACE","caller":"/mnt/tmpfs.ramdisk/log/methods_07_print.go","line":68,"msg":"created logger, level TRACE"}
// {"ts":"2026-05-12T13:34:48+03:00","service":"auth","req_id":12345,"cache_hit":true,"root ends":"here","caller":"/mnt/tmpfs.ramdisk/log/04_log_context_test.go","line":59,"msg":"xxx1"}
// {"ts":"2026-05-12T13:34:48+03:00","service":"auth","req_id":12345,"cache_hit":true,"root ends":"here","area":"some area","caller":"/mnt/tmpfs.ramdisk/log/04_log_context_test.go","line":61,"msg":"login ok again"}
// {"ts":"2026-05-12T13:34:48+03:00","level":"INFO","service":"auth","req_id":12345,"cache_hit":true,"root ends":"here","area":"some area","zzzz":4.299999999999,"g":8,"msg":"finished"}
// {"ts":"2026-05-12T13:34:48+03:00","level":"ERROR","service":"auth","req_id":12345,"cache_hit":true,"root ends":"here","area":"some area","xxxxxxxxxxxxx":"2","g":10,"msg":"some error"}

func TestContext_JSON_Print(t *testing.T) {
	var bufLogs, bufFatal bytes.Buffer

	writer := &bufLogs

	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		writer,
		bytearena.WithTickIfDataMilliseconds(10),
	)
	require.NoError(t, errCrIngestor)

	serviceLogging, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelTrace,

			WithFatalWriter: &bufFatal,
			WithJSON:        true,
			WithCaller:      true,
		},

		WithTimestampRFC3339UTC(t.Context()),
	)
	require.NoError(t, errCrLogger)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	f := NewLogContext(serviceLogging).
		WithRoot("service", "auth").
		SetInt("req_id", 12345).
		SetBool("cache_hit", true).
		SetString("root ends", "here")

	// Execution
	f.Print("xxx1")
	f.SetString("area", "some area")
	f.Print("login ok again")

	go func() {
		f.With("xxxxxxxxxxxxx", "2").
			Error().
			WithGoroutineID().
			Msg("some error")
	}()

	f.With("zzzz", 4.3).
		Info().
		WithGoroutineID().
		Msg("finished")

	// Allow time for the goroutine and the ingestor tick
	time.Sleep(100 * time.Millisecond)

	cancel()
	<-chIngestionEnd

	// --- Processing Lines ---
	require.Zero(t, bufFatal.Len())

	logSet, errParse := query.NewLogset(bufLogs.String())
	require.NoError(t, errParse)

	require.Len(t,
		logSet,
		5,

		logSet,
	)

	require.Len(t, logSet.WithTimestamp(), 5)
	require.NoError(t, logSet.HasKey("msg", 5))
	require.NoError(t, logSet.HasKey("req_id", 4))
	require.NoError(t, logSet.HasKey("level", 3))
	require.NoError(t, logSet.HasKey("caller", 3))
	require.NoError(t, logSet.HasKeyWithValue("service", "auth", 4))
	require.NoError(t, logSet.HasKeyWithValue("level", "TRACE", 1))
	require.NoError(t,
		logSet.HasKeysWithValues(1,
			"req_id", 12345,
			"level", "INFO",
			"zzzz", 4.299999999999,
		),
	)
	require.NoError(t,
		logSet.HasKeysWithValues(1,
			"req_id", 12345,
			"level", "ERROR",
		),
	)
}

// test produces
// 2026-05-12T14:06:06+03:00 /mnt/tmpfs.ramdisk/log/methods_07_print.go Line68 TRACE: created logger, level TRACE
// 2026-05-12T14:06:06+03:00 /mnt/tmpfs.ramdisk/log/04_log_context_test.go Line152 service=auth req_id=12345 cache_hit=true root ends=here msg=xxx1
// 2026-05-12T14:06:06+03:00 /mnt/tmpfs.ramdisk/log/04_log_context_test.go Line154 service=auth req_id=12345 cache_hit=true root ends=here area=some area msg=login ok again
// 2026-05-12T14:06:06+03:00 level=INFO service=auth req_id=12345 cache_hit=true root ends=here area=some area zzzz=4.299999999999 g=8 msg=finished
// 2026-05-12T14:06:06+03:00 level=ERROR service=auth req_id=12345 cache_hit=true root ends=here area=some area xxxxxxxxxxxxx=2 g=10 msg=some error

func TestContext_Raw_Print(t *testing.T) {
	var bufLogs, bufFatal bytes.Buffer

	writer := &bufLogs

	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		writer,
		bytearena.WithTickIfDataMilliseconds(10),
	)
	require.NoError(t, errCrIngestor)

	serviceLogging, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelTrace,

			WithFatalWriter: &bufFatal,
			WithJSON:        false,
			WithCaller:      true,
		},

		WithTimestampRFC3339UTC(t.Context()),
	)
	require.NoError(t, errCrLogger)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	f := NewLogContext(serviceLogging).
		WithRoot("service", "auth").
		SetInt("req_id", 12345).
		SetBool("cache_hit", true).
		SetString("root ends", "here")

	// Execution
	f.Print("xxx1")
	f.SetString("area", "some area")
	f.Print("login ok again")
	f.Prints("hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh")

	go func() {
		f.With("xxxxxxxxxxxxx", "2").
			Error().
			WithGoroutineID().
			Msg("some error")
	}()

	f.With("zzzz", 4.3).
		Info().
		WithGoroutineID().
		Msg("finished")

	// Allow time for the goroutine and the ingestor tick
	time.Sleep(100 * time.Millisecond)

	cancel()
	<-chIngestionEnd

	// --- Processing Lines ---
	require.Zero(t, bufFatal.Len())

	fmt.Println(bufLogs.String())

	logSet, errParse := query.NewRawset(bufLogs.String())
	require.NoError(t, errParse)

	require.Len(t,
		logSet,
		5,

		logSet,
	)

	require.Len(t, logSet.WithTimestamp(), 5)
	require.NoError(t, logSet.HasKey("msg", 4))
	require.NoError(t, logSet.HasKey("req_id", 4))

	require.NoError(t, logSet.HasKey("caller", 3))
	require.NoError(t, logSet.HasKeyWithValue("service", "auth", 4))

	require.NoError(t, logSet.HasKey("level", 2))
	require.NoError(t, logSet.HasKeyWithValue("level", "INFO", 1))
	require.NoError(t, logSet.HasKeyWithValue("level", "ERROR", 1))
	require.NoError(t, logSet.HasKeyWithValue("level", "TRACE", 0))
	require.NoError(t,
		logSet.HasKeysWithValues(1,
			"req_id", 12345,
			"level", "INFO",
			"zzzz", 4.299999999999,
		),
	)
	require.NoError(t,
		logSet.HasKeysWithValues(1,
			"req_id", 12345,
			"level", "ERROR",
		),
	)
}
