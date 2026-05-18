package arenalog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/arenalog/query"
	"github.com/tudorhulban/bytearena"
)

// test produces
// {"ts":"1779110935558730814","level":"TRACE","msg":"created logger, level TRACE"}
// {"ts":"1779110935558733739","level":"TRACE","service":"auth","req_id":12345,"cache_hit":true,"area":"trace-area","user":"arena-trace","msg":"trace message"}
// {"ts":"1779110935558734882","level":"DEBUG","service":"auth","req_id":12345,"cache_hit":true,"area":"debug-area","user":"arena-debug","attempt":1,"msg":"debug message"}
// {"ts":"1779110935558735593","level":"INFO","service":"auth","req_id":12345,"cache_hit":true,"area":"info-area","user":"arena-info","some_float":1.113699999999,"success":true,"msg":"info message"}
// {"ts":"1779110935558736615","level":"ERROR","service":"auth","req_id":12345,"cache_hit":true,"area":"error-area","user":"arena-error","error_detail":"something failed","msg":"error message"}

func TestArenalog_MultipleFields_AllLevels(t *testing.T) {
	// 1. Create a buffer to capture output
	var bufLogs, bufFatal bytes.Buffer

	// Define constant for req_id
	const expectedReqID = 12345

	// 2. Setup Ingestor with the buffer as the writer
	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		&bufLogs,
	)
	require.NoError(t, errCrIngestor)
	require.NotNil(t, ingestor)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	// 3. Setup Logger
	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelTrace,

			WithFatalWriter: &bufFatal,
			WithJSON:        true,
		},
	)
	require.NoError(t, errCrLogger)

	// 4. Create Context with Base Fields
	logContext := NewLogContext(logger).
		WithRoot("service", "auth").
		SetInt("req_id", expectedReqID).
		SetBool("cache_hit", true)

	// 5. Generate Logs for All Levels

	// --- TRACE ---
	logContext.
		WithString("area", "trace-area").
		Trace().
		WithString("user", "arena-trace").
		Msg("trace message")

	// --- DEBUG ---
	logContext.
		WithString("area", "debug-area").
		Debug().
		WithString("user", "arena-debug").
		WithInt("attempt", 1).
		Msg("debug message")

	// --- INFO ---
	logContext.
		WithString("area", "info-area").
		Info().
		WithString("user", "arena-info").
		WithFloat("some_float", 1.1137).
		WithBool("success", true).
		Msg("info message")

	// --- ERROR ---
	logContext.
		WithString("area", "error-area").
		Error().
		WithString("user", "arena-error").
		WithString("error_detail", "something failed").
		Msg("error message")

	// 6. Stop Ingestor and Wait for flush
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
	require.NoError(t, logSet.HasKey("level", 5))
	require.NoError(t, logSet.HasKey("msg", 5))
	require.NoError(t, logSet.HasKey("user", 4))
	require.NoError(t, logSet.HasKey("area", 4))
	require.NoError(t, logSet.HasKeyWithValue("service", "auth", 4))
	require.NoError(t, logSet.HasKeyWithValue("level", "TRACE", 2))
	require.NoError(t, logSet.ContainsEach(1, "DEBUG", "INFO", "ERROR"))
	require.NoError(t, logSet.HasKeyWithValueInDelta("some_float", 1.1137, 0.0001, 1))
}

// test produces
// {"ts":"1779112304402035665","level":"TRACE","msg":"created logger, level TRACE"}
// {"ts":"1779112304402037749","level":"TRACE","entry start":"trace-area","component":"scanner","msg":"minimal trace"}
// {"ts":"1779112304402038841","level":"DEBUG","entry start":"debug-area","component":"scanner","msg":"minimal debug"}
// {"ts":"1779112304402039212","level":"INFO","entry start":"info-area","component":"scanner","msg":"minimal info"}
// {"ts":"1779112304402039653","level":"ERROR","entry start":"error-area","code":500,"msg":"minimal error"}

func TestArenalog_NoRootFields(t *testing.T) {
	// 1. Create a buffer to capture output
	var buf bytes.Buffer

	// 2. Setup Ingestor
	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		&buf,
	)
	require.NoError(t, errCrIngestor)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	// 3. Setup Logger
	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelTrace,

			WithFatalWriter: &buf,
			WithJSON:        true,
		},
	)
	require.NoError(t, errCrLogger)

	// 4. Create Context WITHOUT Root Fields
	logContext := NewLogContext(logger)

	// 5. Generate Logs

	// --- TRACE (Entry info only) ---
	logContext.
		WithString("entry start", "trace-area").
		Trace().
		WithString("component", "scanner").
		Msg("minimal trace")

		// --- DEBUG (Entry info only) ---
	logContext.
		WithString("entry start", "debug-area").
		Debug().
		WithString("component", "scanner").
		Msg("minimal debug")

		// --- INFO (Entry info only) ---
	logContext.
		WithString("entry start", "info-area").
		Info().
		WithString("component", "scanner").
		Msg("minimal info")

		// --- ERROR (Entry info only) ---
	logContext.
		WithString("entry start", "error-area").
		Error().
		WithInt("code", 500).
		Msg("minimal error")

	// 6. Stop Ingestor
	cancel()
	<-chIngestionEnd

	// 7. Parse and Assert
	output := buf.String()
	linesRaw := strings.Split(output, "\n")

	fmt.Println(output)

	var linesJSON []string

	for _, line := range linesRaw {
		trimmed := strings.TrimSpace(line)

		if strings.Contains(trimmed, "{") {
			linesJSON = append(linesJSON, trimmed[strings.Index(trimmed, "{"):]) //nolint:gocritic
		}
	}

	// Expect 5 lines: 1 Init + 4 Log messages
	require.Equal(t, 5, len(linesJSON))

	parseLog := func(line string) map[string]any {
		var m map[string]any
		require.NoError(t, json.Unmarshal([]byte(line), &m))

		return m
	}

	// Helper to ensure root fields from previous tests ARE NOT present
	assertNoRootFields := func(logData map[string]any) {
		require.Nil(t, logData["service"], "Root field 'service' should not exist")
		require.Nil(t, logData["req_id"], "Root field 'req_id' should not exist")
		require.Nil(t, logData["cache_hit"], "Root field 'cache_hit' should not exist")
	}

	// --- verify TRACE ---
	traceLog := parseLog(linesJSON[1])
	assertNoRootFields(traceLog)
	require.Equal(t, "TRACE", traceLog["level"])
	require.Equal(t, "scanner", traceLog["component"]) // Entry Info
	require.Equal(t, "minimal trace", traceLog["msg"])

	// --- verify DEBUG ---
	debugLog := parseLog(linesJSON[2])
	assertNoRootFields(debugLog)
	require.Equal(t, "DEBUG", debugLog["level"])
	require.Equal(t, "scanner", debugLog["component"]) // Entry Info
	require.Equal(t, "minimal debug", debugLog["msg"])

	// --- verify INFO ---
	infoLog := parseLog(linesJSON[3])
	assertNoRootFields(infoLog)
	require.Equal(t, "INFO", infoLog["level"])
	require.Equal(t, "scanner", infoLog["component"]) // Entry Info
	require.Equal(t, "minimal info", infoLog["msg"])

	// --- verify ERROR ---
	errorLog := parseLog(linesJSON[4])
	assertNoRootFields(errorLog)
	require.Equal(t, "ERROR", errorLog["level"])
	require.Equal(t, float64(500), errorLog["code"]) // Entry Info
	require.Equal(t, "minimal error", errorLog["msg"])
}

// test produces
// 2026-05-14T13:21:18+03:00 /mnt/tmpfs.ramdisk/log/methods_07_print.go Line68 TRACE: created logger, level TRACE
// 2026-05-14T13:21:18+03:00 /mnt/tmpfs.ramdisk/log/methods_00_trace.go Line18 TRACE: some text
// 2026-05-14T13:21:18+03:00 /mnt/tmpfs.ramdisk/log/methods_01_debug.go Line18 DEBUG: some text
// 2026-05-14T13:21:18+03:00 /mnt/tmpfs.ramdisk/log/methods_02_info.go Line18 INFO: some text
// 2026-05-14T13:21:18+03:00 /mnt/tmpfs.ramdisk/log/methods_03_warn.go Line18 WARN: some text
// 2026-05-14T13:21:18+03:00 /mnt/tmpfs.ramdisk/log/methods_04_error.go Line18 ERROR: some text
// 2026-05-14T13:21:18+03:00 /mnt/tmpfs.ramdisk/log/methods_07_print.go Line60 TRACE: some text

func TestMissing_Raw_Mix(t *testing.T) {
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

	l.Trace("some text")
	l.Debug("some text")
	l.Info("some text")
	l.Warn("some text")
	l.Error("some text")
	l.Print("some text")

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	cancel()
	<-chIngestionEnd

	out := writer.String()
	require.NotEmpty(t, out)

	logSet, errCr := query.NewLogset(out)
	require.NoError(t, errCr)

	require.NoError(t, logSet.ContainsLike(7, "Line"))
	require.NoError(t, logSet.Contains(6, "some text"))
	require.NoError(t, logSet.Contains(3, "TRACE"))
	require.NoError(t, logSet.ContainsEach(1, "DEBUG", "INFO", "WARN", "ERROR"))
}
