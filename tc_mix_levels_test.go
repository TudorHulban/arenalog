package arenalog

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena"
)

// PRINT is not a severity level. It is always emitted regardless of the
// configured logLevel, but its *label* is inherited from the current
// threshold. This means:
//
//   logLevel = TRACE → PRINT entries use "TRACE"
//   logLevel = DEBUG → PRINT entries use "DEBUG"
//   logLevel = INFO  → PRINT entries use "INFO"
//   logLevel = WARN  → PRINT entries use "WARN"
//   logLevel = ERROR → PRINT entries use "ERROR"
//   logLevel = FATAL → PRINT entries use "FATAL"
//   logLevel = PANIC → PRINT entries use "PANIC"
//
// PRINT therefore bypasses filtering, but does not have its own label.
// Tests must assert the inherited label, not "PRINT".

func TestLevelsMatrix(t *testing.T) {
	type tc struct {
		description string
		shouldSee   map[string]string // msg → expected level label in JSON
		shouldSkip  []string          // msgs that must not appear
		level       Level
	}

	tests := []tc{
		{
			description: "1. TRACE threshold → all safe levels emitted",
			level:       LevelTrace,
			shouldSee: map[string]string{
				"trc": `"level":"TRACE"`,
				"dbg": `"level":"DEBUG"`,
				"inf": `"level":"INFO"`,
				"wrn": `"level":"WARN"`,
				"err": `"level":"ERROR"`,
				"prt": `"level":"` + logLevels[LevelTrace] + `"`,
			},
			shouldSkip: nil,
		},
		{
			description: "2. DEBUG threshold → DEBUG+ emitted; TRACE suppressed",
			level:       LevelDebug,
			shouldSee: map[string]string{
				"dbg": `"level":"DEBUG"`,
				"inf": `"level":"INFO"`,
				"wrn": `"level":"WARN"`,
				"err": `"level":"ERROR"`,
				"prt": `"level":"` + logLevels[LevelDebug] + `"`,
			},
			shouldSkip: []string{"trc"},
		},
		{
			description: "3. INFO threshold → INFO+ emitted; TRACE, DEBUG suppressed",
			level:       LevelInfo,
			shouldSee: map[string]string{
				"inf": `"level":"INFO"`,
				"wrn": `"level":"WARN"`,
				"err": `"level":"ERROR"`,
				"prt": `"level":"` + logLevels[LevelInfo] + `"`,
			},
			shouldSkip: []string{"trc", "dbg"},
		},
		{
			description: "4. WARN threshold → WARN+ emitted; TRACE-INFO suppressed",
			level:       LevelWarn,
			shouldSee: map[string]string{
				"wrn": `"level":"WARN"`,
				"err": `"level":"ERROR"`,
				"prt": `"level":"` + logLevels[LevelWarn] + `"`,
			},
			shouldSkip: []string{"trc", "dbg", "inf"},
		},
		{
			description: "5. ERROR threshold → ERROR+ emitted; TRACE-WARN suppressed",
			level:       LevelError,
			shouldSee: map[string]string{
				"err": `"level":"ERROR"`,
				"prt": `"level":"` + logLevels[LevelError] + `"`,
			},
			shouldSkip: []string{"trc", "dbg", "inf", "wrn"},
		},
		{
			description: "6. FATAL threshold → only PRINT emitted (safe levels suppressed)",
			level:       LevelFatal,
			shouldSee: map[string]string{
				"prt": `"level":"` + logLevels[LevelFatal] + `"`,
			},
			shouldSkip: []string{"trc", "dbg", "inf", "wrn", "err"},
		},
		{
			description: "7. PANIC threshold → only PRINT emitted (all safe levels suppressed)",
			level:       LevelPanic,
			shouldSee: map[string]string{
				"prt": `"level":"` + logLevels[LevelPanic] + `"`,
			},
			shouldSkip: []string{"trc", "dbg", "inf", "wrn", "err"},
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.description,
			func(t *testing.T) {
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
						LoggerLevel: tc.level,

						WithFatalWriter: os.Stdout,
						WithCaller:      true,
						WithColor:       false,
						WithJSON:        true,
					},
				)
				require.NoError(t, errCrLogger)

				ctx, cancel := context.WithCancel(context.Background())
				chIngestionEnd := ingestor.StartIngestion(ctx)

				// Emit only non-terminating log levels.
				// NOTE: Fatal() and Panic() are excluded from this matrix test
				// because they may call os.Exit() or panic(), which would
				// terminate the test process. Their filtering logic is identical
				// to other levels (entry.level >= threshold), so coverage is
				// maintained by verifying suppression of lower levels.
				l.Trace("trc")
				l.Debug("dbg")
				l.Info("inf")
				l.Warn("wrn")
				l.Error("err")
				l.Print("prt")

				cancel()
				<-chIngestionEnd

				out := writer.String()
				require.NotEmpty(t, out)

				lines := strings.Split(strings.TrimSpace(out), "\n")
				require.NotEmpty(t, lines)

				// Verify expected entries appear with correct level label
				for msg, expectedLevelJSON := range tc.shouldSee {
					found := false

					for _, ln := range lines {
						if strings.Contains(ln, `"msg":"`+msg+`"`) &&
							strings.Contains(ln, expectedLevelJSON) &&
							strings.Contains(ln, `"ts":`) &&
							strings.Contains(ln, `"caller":`) {
							found = true

							break
						}
					}

					require.True(t, found,
						"expected to see msg=%q with level=%q in output:\n%s",
						msg, expectedLevelJSON, out)
				}

				// Verify suppressed entries do not appear
				for _, msg := range tc.shouldSkip {
					for _, ln := range lines {
						require.NotContains(t,
							ln,
							`"msg":"`+msg+`"`,

							"msg=%q must be suppressed at threshold=%s, output:\n%s",
							msg,
							tc.level.String(),
							out,
						)
					}
				}
			},
		)
	}
}

func TestLevelsfMatrix(t *testing.T) {
	type tc struct {
		description string
		shouldSee   map[string]string // msg → expected level label in JSON
		shouldSkip  []string          // msgs that must not appear
		level       Level
	}

	tests := []tc{
		{
			description: "1. TRACE threshold → all safe levels emitted",
			level:       LevelTrace,
			shouldSee: map[string]string{
				"trc": `"level":"TRACE"`,
				"dbg": `"level":"DEBUG"`,
				"inf": `"level":"INFO"`,
				"wrn": `"level":"WARN"`,
				"err": `"level":"ERROR"`,
				"prt": `"level":"` + logLevels[LevelTrace] + `"`,
			},
			shouldSkip: nil,
		},
		{
			description: "2. DEBUG threshold → DEBUG+ emitted; TRACE suppressed",
			level:       LevelDebug,
			shouldSee: map[string]string{
				"dbg": `"level":"DEBUG"`,
				"inf": `"level":"INFO"`,
				"wrn": `"level":"WARN"`,
				"err": `"level":"ERROR"`,
				"prt": `"level":"` + logLevels[LevelDebug] + `"`,
			},
			shouldSkip: []string{"trc"},
		},
		{
			description: "3. INFO threshold → INFO+ emitted; TRACE, DEBUG suppressed",
			level:       LevelInfo,
			shouldSee: map[string]string{
				"inf": `"level":"INFO"`,
				"wrn": `"level":"WARN"`,
				"err": `"level":"ERROR"`,
				"prt": `"level":"` + logLevels[LevelInfo] + `"`,
			},
			shouldSkip: []string{"trc", "dbg"},
		},
		{
			description: "4. WARN threshold → WARN+ emitted; TRACE-INFO suppressed",
			level:       LevelWarn,
			shouldSee: map[string]string{
				"wrn": `"level":"WARN"`,
				"err": `"level":"ERROR"`,
				"prt": `"level":"` + logLevels[LevelWarn] + `"`,
			},
			shouldSkip: []string{"trc", "dbg", "inf"},
		},
		{
			description: "5. ERROR threshold → ERROR+ emitted; TRACE-WARN suppressed",
			level:       LevelError,
			shouldSee: map[string]string{
				"err": `"level":"ERROR"`,
				"prt": `"level":"` + logLevels[LevelError] + `"`,
			},
			shouldSkip: []string{"trc", "dbg", "inf", "wrn"},
		},
		{
			description: "6. FATAL threshold → only PRINT emitted (safe levels suppressed)",
			level:       LevelFatal,
			shouldSee: map[string]string{
				"prt": `"level":"` + logLevels[LevelFatal] + `"`,
			},
			shouldSkip: []string{"trc", "dbg", "inf", "wrn", "err"},
		},
		{
			description: "7. PANIC threshold → only PRINT emitted (all safe levels suppressed)",
			level:       LevelPanic,
			shouldSee: map[string]string{
				"prt": `"level":"` + logLevels[LevelPanic] + `"`,
			},
			shouldSkip: []string{"trc", "dbg", "inf", "wrn", "err"},
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.description,
			func(t *testing.T) {
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
						LoggerLevel: tc.level,

						WithFatalWriter: os.Stdout,
						WithCaller:      true,
						WithColor:       false,
						WithJSON:        true,
					},
				)
				require.NoError(t, errCrLogger)

				ctx, cancel := context.WithCancel(context.Background())
				chIngestionEnd := ingestor.StartIngestion(ctx)

				// Emit only non-terminating log levels.
				// NOTE: Fatal() and Panic() are excluded from this matrix test
				// because they may call os.Exit() or panic(), which would
				// terminate the test process. Their filtering logic is identical
				// to other levels (entry.level >= threshold), so coverage is
				// maintained by verifying suppression of lower levels.
				l.Tracef("%s", "trc")
				l.Debugf("%s", "dbg")
				l.Infof("%s", "inf")
				l.Warnf("%s", "wrn")
				l.Errorf("%s", "err")
				l.Printf("%s", "prt")

				cancel()
				<-chIngestionEnd

				out := writer.String()
				require.NotEmpty(t, out)

				lines := strings.Split(strings.TrimSpace(out), "\n")
				require.NotEmpty(t, lines)

				// Verify expected entries appear with correct level label
				for msg, expectedLevelJSON := range tc.shouldSee {
					found := false

					for _, ln := range lines {
						if strings.Contains(ln, `"msg":"`+msg+`"`) &&
							strings.Contains(ln, expectedLevelJSON) &&
							strings.Contains(ln, `"ts":`) &&
							strings.Contains(ln, `"caller":`) {
							found = true

							break
						}
					}

					require.True(t, found,
						"expected to see msg=%q with level=%q in output:\n%s",
						msg, expectedLevelJSON, out)
				}

				// Verify suppressed entries do not appear
				for _, msg := range tc.shouldSkip {
					for _, ln := range lines {
						require.NotContains(t,
							ln,
							`"msg":"`+msg+`"`,

							"msg=%q must be suppressed at threshold=%s, output:\n%s",
							msg,
							tc.level.String(),
							out,
						)
					}
				}
			},
		)
	}
}

func TestLevelswMatrix(t *testing.T) {
	type tc struct {
		description string
		shouldSee   map[string]string // msg → expected level label in JSON
		shouldSkip  []string          // msgs that must not appear
		level       Level
	}

	tests := []tc{
		{
			description: "1. TRACE threshold → all safe levels emitted",
			level:       LevelTrace,
			shouldSee: map[string]string{
				"trc": `"level":"TRACE"`,
				"dbg": `"level":"DEBUG"`,
				"inf": `"level":"INFO"`,
				"wrn": `"level":"WARN"`,
				"err": `"level":"ERROR"`,
				"prt": `"level":"` + logLevels[LevelTrace] + `"`,
			},
			shouldSkip: nil,
		},
		{
			description: "2. DEBUG threshold → DEBUG+ emitted; TRACE suppressed",
			level:       LevelDebug,
			shouldSee: map[string]string{
				"dbg": `"level":"DEBUG"`,
				"inf": `"level":"INFO"`,
				"wrn": `"level":"WARN"`,
				"err": `"level":"ERROR"`,
				"prt": `"level":"` + logLevels[LevelDebug] + `"`,
			},
			shouldSkip: []string{"trc"},
		},
		{
			description: "3. INFO threshold → INFO+ emitted; TRACE, DEBUG suppressed",
			level:       LevelInfo,
			shouldSee: map[string]string{
				"inf": `"level":"INFO"`,
				"wrn": `"level":"WARN"`,
				"err": `"level":"ERROR"`,
				"prt": `"level":"` + logLevels[LevelInfo] + `"`,
			},
			shouldSkip: []string{"trc", "dbg"},
		},
		{
			description: "4. WARN threshold → WARN+ emitted; TRACE-INFO suppressed",
			level:       LevelWarn,
			shouldSee: map[string]string{
				"wrn": `"level":"WARN"`,
				"err": `"level":"ERROR"`,
				"prt": `"level":"` + logLevels[LevelWarn] + `"`,
			},
			shouldSkip: []string{"trc", "dbg", "inf"},
		},
		{
			description: "5. ERROR threshold → ERROR+ emitted; TRACE-WARN suppressed",
			level:       LevelError,
			shouldSee: map[string]string{
				"err": `"level":"ERROR"`,
				"prt": `"level":"` + logLevels[LevelError] + `"`,
			},
			shouldSkip: []string{"trc", "dbg", "inf", "wrn"},
		},
		{
			description: "6. FATAL threshold → only PRINT emitted (safe levels suppressed)",
			level:       LevelFatal,
			shouldSee: map[string]string{
				"prt": `"level":"` + logLevels[LevelFatal] + `"`,
			},
			shouldSkip: []string{"trc", "dbg", "inf", "wrn", "err"},
		},
		{
			description: "7. PANIC threshold → only PRINT emitted (all safe levels suppressed)",
			level:       LevelPanic,
			shouldSee: map[string]string{
				"prt": `"level":"` + logLevels[LevelPanic] + `"`,
			},
			shouldSkip: []string{"trc", "dbg", "inf", "wrn", "err"},
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.description,
			func(t *testing.T) {
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
						LoggerLevel: tc.level,

						WithFatalWriter: os.Stdout,
						WithCaller:      true,
						WithColor:       false,
						WithJSON:        true,
					},
				)
				require.NoError(t, errCrLogger)

				ctx, cancel := context.WithCancel(context.Background())
				chIngestionEnd := ingestor.StartIngestion(ctx)

				// Emit only non-terminating log levels.
				// NOTE: Fatal() and Panic() are excluded from this matrix test
				// because they may call os.Exit() or panic(), which would
				// terminate the test process. Their filtering logic is identical
				// to other levels (entry.level >= threshold), so coverage is
				// maintained by verifying suppression of lower levels.
				l.Tracew("trc")
				l.Debugw("dbg")
				l.Infow("inf")
				l.Warnw("wrn")
				l.Errorw("err")
				l.Printw("prt")

				cancel()
				<-chIngestionEnd

				out := writer.String()
				require.NotEmpty(t, out)

				lines := strings.Split(strings.TrimSpace(out), "\n")
				require.NotEmpty(t, lines)

				// Verify expected entries appear with correct level label
				for msg, expectedLevelJSON := range tc.shouldSee {
					found := false

					for _, ln := range lines {
						if strings.Contains(ln, `"msg":"`+msg+`"`) &&
							strings.Contains(ln, expectedLevelJSON) &&
							strings.Contains(ln, `"ts":`) &&
							strings.Contains(ln, `"caller":`) {
							found = true

							break
						}
					}

					require.True(t, found,
						"expected to see msg=%q with level=%q in output:\n%s",
						msg, expectedLevelJSON, out)
				}

				// Verify suppressed entries do not appear
				for _, msg := range tc.shouldSkip {
					for _, ln := range lines {
						require.NotContains(t,
							ln,
							`"msg":"`+msg+`"`,

							"msg=%q must be suppressed at threshold=%s, output:\n%s",
							msg,
							tc.level.String(),
							out,
						)
					}
				}
			},
		)
	}
}
