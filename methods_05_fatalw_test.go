package arenalog

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/arenalog/query"
	"github.com/tudorhulban/bytearena"
)

// test produces
// {"ts":"1778663357411073623","level":"FATAL","msg":"LOG_ERR(odd_args): failure","service":"(MISSING)"}

func TestLogger_JSON_Fatalw(t *testing.T) {
	// This is the main test that spawns the sub-process
	if os.Getenv("BE_CRASHING") == "1" {
		ingestor, _ := bytearena.NewIngestor(
			bytearena.Size100K(),
			io.Discard,
		)

		// Use os.Stdout so the parent process can capture it in the buffer
		l, _ := NewLogger(
			&ParamsNewLogger{
				Ingestor:    ingestor,
				LoggerLevel: LevelTrace,

				WithFatalWriter: os.Stdout,
				WithJSON:        true,
			},
		)
		if l == nil {
			os.Exit(99)
		}

		_ = ingestor.StartIngestion(context.Background())

		l.Fatalw("failure", "service")
	}

	// --- Parent Process Logic (The Actual Test) ---
	var buf bytes.Buffer

	// We capture the output of the sub-process
	cmd := exec.CommandContext(
		context.Background(),
		os.Args[0],
		"-test.run=^TestLogger_JSON_Fatalw$",
	)

	cmd.Env = append(os.Environ(), "BE_CRASHING=1")
	cmd.Stdout = &buf
	cmd.Stderr = &buf // Fatal might write to stderr depending on the writer

	err := cmd.Run()

	// 1. Verify that the process actually exited with an error code (1)
	var exitErr *exec.ExitError
	require.ErrorAs(t, err, &exitErr)
	require.Equal(t, 1, exitErr.ExitCode(), "Expected os.Exit(1)")

	// 2. Use your DSL to verify the content of the output
	logSet, errCr := query.NewLogset(buf.String())
	require.NoError(t, errCr)
	assert.Len(t, logSet, 1)

	require.NotZero(t,
		logSet.WithTimestamp(),
		"timestamp should be present",
	)
	require.NoError(t, logSet.First().HasKeyWithValue("level", "FATAL", 1))
	require.NoError(t, logSet.First().HasKeyWithValueLike("msg", "odd_args", 1))
}

// test produces
// 1778663478385819188 [FATAL] LOG_ERR(odd_args): failure service=(MISSING)

func TestLogger_Raw_Fatalw(t *testing.T) {
	// This is the main test that spawns the sub-process
	if os.Getenv("BE_CRASHING") == "1" {
		ingestor, _ := bytearena.NewIngestor(
			bytearena.Size100K(),
			io.Discard,
		)

		// Use os.Stdout so the parent process can capture it in the buffer
		l, _ := NewLogger(
			&ParamsNewLogger{
				Ingestor:    ingestor,
				LoggerLevel: LevelTrace,

				WithFatalWriter: os.Stdout,
				WithJSON:        false,
			},
		)
		if l == nil {
			os.Exit(99)
		}

		_ = ingestor.StartIngestion(context.Background())

		l.Fatalw("failure", "service")
	}

	// --- Parent Process Logic (The Actual Test) ---
	var buf bytes.Buffer

	// We capture the output of the sub-process
	cmd := exec.CommandContext(
		context.Background(),
		os.Args[0],
		"-test.run=^TestLogger_Raw_Fatalw$",
	)

	cmd.Env = append(os.Environ(), "BE_CRASHING=1")
	cmd.Stdout = &buf
	cmd.Stderr = &buf // Fatal might write to stderr depending on the writer

	err := cmd.Run()

	// 1. Verify that the process actually exited with an error code (1)
	var exitErr *exec.ExitError
	require.ErrorAs(t, err, &exitErr)
	require.Equal(t, 1, exitErr.ExitCode(), "Expected os.Exit(1)")

	// 2. Use your DSL to verify the content of the output
	logSet, errCr := query.NewLogset(buf.String())
	require.NoError(t, errCr)
	assert.Len(t, logSet, 1)

	require.NoError(t, logSet.First().Contains(1, "[FATAL]"))
	require.NoError(t, logSet.First().Contains(1, "odd_args"))
}

// test produces
// 1778663721904231440 [FATAL] failure service=infra

func TestLogger_Raw_Fatalw_Even(t *testing.T) {
	// This is the main test that spawns the sub-process
	if os.Getenv("BE_CRASHING") == "1" {
		ingestor, _ := bytearena.NewIngestor(
			bytearena.Size100K(),
			io.Discard,
		)

		// Use os.Stdout so the parent process can capture it in the buffer
		l, _ := NewLogger(
			&ParamsNewLogger{
				Ingestor:    ingestor,
				LoggerLevel: LevelTrace,

				WithFatalWriter: os.Stdout,
				WithJSON:        false,
			},
		)
		if l == nil {
			os.Exit(99)
		}

		_ = ingestor.StartIngestion(context.Background())

		l.Fatalw("failure", "service", "infra")
	}

	// --- Parent Process Logic (The Actual Test) ---
	var buf bytes.Buffer

	// We capture the output of the sub-process
	cmd := exec.CommandContext(
		context.Background(),
		os.Args[0],
		"-test.run=^TestLogger_Raw_Fatalw_Even$",
	)

	cmd.Env = append(os.Environ(), "BE_CRASHING=1")
	cmd.Stdout = &buf
	cmd.Stderr = &buf // Fatal might write to stderr depending on the writer

	err := cmd.Run()

	// 1. Verify that the process actually exited with an error code (1)
	var exitErr *exec.ExitError
	require.ErrorAs(t, err, &exitErr)
	require.Equal(t, 1, exitErr.ExitCode(), "Expected os.Exit(1)")

	// 2. Use your DSL to verify the content of the output
	logSet, errCr := query.NewLogset(buf.String())
	require.NoError(t, errCr)
	assert.Len(t, logSet, 1)

	require.NoError(t, logSet.First().Contains(1, "[FATAL]"))
	require.NoError(t, logSet.First().Contains(1, "service=infra"))
}
