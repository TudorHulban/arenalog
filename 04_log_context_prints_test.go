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
// {"ts":"2026-05-26T09:30:09.687Z","level":"TRACE","caller":"/mnt/tmpfs.ramdisk/arenalog/methods_07_print.go","line":68,"msg":"created logger, level TRACE"}
// {"ts":"2026-05-26T09:30:09.687Z","service":"auth","req_id":12345,"cache_hit":true,"root ends":"here","caller":"/mnt/tmpfs.ramdisk/arenalog/04_log_context_prints_test.go","line":51,"msg":"xxx1"}
// {"ts":"2026-05-26T09:30:09.687Z","service":"auth","req_id":12345,"cache_hit":true,"root ends":"here","area":"some area","caller":"/mnt/tmpfs.ramdisk/arenalog/04_log_context_prints_test.go","line":53,"msg":"login ok again"}

func TestContext_JSON_Prints(t *testing.T) {
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
	f.Prints("xxx1")
	f.SetString("area", "some area")
	f.Prints("login ok again")

	// Allow time for the goroutine and the ingestor tick
	time.Sleep(100 * time.Millisecond)

	cancel()
	<-chIngestionEnd

	// --- Processing Lines ---
	require.Zero(t, bufFatal.Len())

	fmt.Println(bufLogs.String())

	logSet, errParse := query.NewLogset(bufLogs.String())
	require.NoError(t, errParse)

	require.Len(t,
		logSet,
		3,

		logSet,
	)

	require.Len(t, logSet.WithTimestamp(), 3)
	require.NoError(t, logSet.HasKey("msg", 3))
	require.NoError(t, logSet.HasKey("req_id", 2))
	require.NoError(t, logSet.HasKey("level", 1))
	require.NoError(t, logSet.HasKey("caller", 3))
	require.NoError(t, logSet.HasKeyWithValue("service", "auth", 2))
	require.NoError(t, logSet.HasKeyWithValue("level", "TRACE", 1))
}

// test produces
// 2026-05-26T09:42:39.028Z /mnt/tmpfs.ramdisk/arenalog/methods_07_print.go Line68 TRACE: created logger, level TRACE
// 2026-05-26T09:42:39.028Z /mnt/tmpfs.ramdisk/arenalog/04_log_context_prints_test.go Line126 service=auth req_id=12345 cache_hit=true root ends=here msg=xxx1
// 2026-05-26T09:42:39.028Z /mnt/tmpfs.ramdisk/arenalog/04_log_context_prints_test.go Line128 service=auth req_id=12345 cache_hit=true root ends=here area=some area msg=login ok again

func TestContext_Raw_Prints(t *testing.T) {
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
	f.Prints("xxx1")
	f.SetString("area", "some area")
	f.Prints("login ok again")

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
		3,

		logSet,
	)

	require.Len(t, logSet.WithTimestamp(), 3)
	require.NoError(t, logSet.HasKey("msg", 2))
	require.NoError(t, logSet.HasKey("req_id", 2))
	require.NoError(t, logSet.Contains(1, "level"))
	require.NoError(t, logSet.HasKey("caller", 3))
	require.NoError(t, logSet.HasKeyWithValue("service", "auth", 2))
	require.NoError(t, logSet.ContainsLike(1, "TRACE"))
}
