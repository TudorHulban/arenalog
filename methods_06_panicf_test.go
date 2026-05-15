package arenalog

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/arenalog/query"
	"github.com/tudorhulban/bytearena"
)

// test produces
// panic: {"ts":"1778686100197545957","level":"PANIC","msg":"core system panic"} [recovered, repanicked]

func TestLogger_JSON_Panicf(t *testing.T) {
	// This is the main test that spawns the sub-process
	if os.Getenv("BE_PANICKING") == "1" {
		ingestor, _ := bytearena.NewIngestor(
			bytearena.Size100K(),
			io.Discard,
		)

		l, _ := NewLogger(
			&ParamsNewLogger{
				Ingestor:    ingestor,
				LoggerLevel: LevelTrace,

				WithFatalWriter: os.Stdout, // Capture the formatted panic log here
				WithJSON:        true,
			},
		)
		if l == nil {
			os.Exit(99)
		}

		_ = ingestor.StartIngestion(context.Background())

		// This call will trigger a panic, which causes the process to crash
		l.Panicf("%s system panic", "core")
	}

	// --- Parent Process Logic ---
	var buf bytes.Buffer

	cmd := exec.CommandContext(
		context.Background(),
		os.Args[0],
		"-test.run=^TestLogger_JSON_Panicf$",
	)

	cmd.Env = append(os.Environ(), "BE_PANICKING=1")
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	err := cmd.Run()

	// 1. Verify that the process exited with an error code.
	// Go panics usually exit with code 2.
	var exitErr *exec.ExitError
	require.ErrorAs(t, err, &exitErr)
	require.True(t, exitErr.ExitCode() != 0, "Expected non-zero exit code from panic")

	// 2. Parse the output and verify the log entry captured from fatalWriter
	logSet, errCr := query.NewLogset(buf.String())
	require.NoError(t, errCr)

	cleanedLine, errClean := query.ExtractLogFromPanic(
		logSet.Skip(1).First().String(),
	)
	require.NoError(t, errClean)

	cleanLogset, errCrClean := query.NewLogset(cleanedLine)
	require.NoError(t, errCrClean)

	// We look for the first log entry in the set
	require.NotEmpty(t,
		cleanLogset,
		"Should have captured at least one log line",
	)

	require.NotZero(t, cleanLogset.WithTimestamp(), "timestamp should be present")
	require.NoError(t, cleanLogset.HasKeyWithValue("level", "PANIC", 1))
	require.NoError(t, cleanLogset.HasKeyWithValue("msg", "core system panic", 1))
}

// test produces
// panic: 1778686169669376532 [PANIC] core system panic [recovered, repanicked]

func TestLogger_Raw_Panicf(t *testing.T) {
	// This is the main test that spawns the sub-process
	if os.Getenv("BE_PANICKING") == "1" {
		ingestor, _ := bytearena.NewIngestor(
			bytearena.Size100K(),
			io.Discard,
		)

		l, _ := NewLogger(
			&ParamsNewLogger{
				Ingestor:    ingestor,
				LoggerLevel: LevelTrace,

				WithFatalWriter: os.Stdout, // Capture the formatted panic log here
				WithJSON:        false,
			},
		)
		if l == nil {
			os.Exit(99)
		}

		_ = ingestor.StartIngestion(context.Background())

		// This call will trigger a panic, which causes the process to crash
		l.Panicf("%s system panic", "core")
	}

	// --- Parent Process Logic ---
	var buf bytes.Buffer

	cmd := exec.CommandContext(
		context.Background(),
		os.Args[0],
		"-test.run=^TestLogger_Raw_Panicf$",
	)

	cmd.Env = append(os.Environ(), "BE_PANICKING=1")
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	err := cmd.Run()

	// 1. Verify that the process exited with an error code.
	// Go panics usually exit with code 2.
	var exitErr *exec.ExitError
	require.ErrorAs(t, err, &exitErr)
	require.True(t, exitErr.ExitCode() != 0, "Expected non-zero exit code from panic")

	// 2. Parse the output and verify the log entry captured from fatalWriter
	logSet, errCr := query.NewRawset(buf.String())
	require.NoError(t, errCr)

	// We look for the first log entry in the set
	require.NotEmpty(t,
		logSet,
		"Should have captured at least one log line",
	)

	require.NoError(t, logSet.Skip(1).First().Contains(1, "PANIC"))
	require.NoError(t, logSet.Skip(1).First().Contains(1, "core system panic"))
}
