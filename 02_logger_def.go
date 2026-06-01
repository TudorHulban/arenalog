package arenalog

import (
	"errors"
	"io"
	"sync/atomic"

	"github.com/tudorhulban/arenalog/timestamp"
	"github.com/tudorhulban/bytearena"
)

type Logger struct { //nolint:govet
	// Cache line 0 (hot, bytes 0..63):
	// - atomic level and per-level estimated sizes are read on every log call.
	// - booleans and callerLevel are also hot (formatting + caller decision).
	// - ingestor and fnTimestamp are hot-ish and placed here so a single cache load
	//   brings the common hot state into L1.
	logLevel atomic.Uint32 // 4

	estimatedMessageSizeOverall uint32 // 4
	callerLevel                 uint8  // 1
	withCaller                  bool   // 1
	withColor                   bool   // 1
	withJSON                    bool   // 1

	ingestor    *bytearena.Ingestor // 8
	fnTimestamp timestamp.Timestamp // 8

	// pad to end of 64‑byte cache line
	_pad0 [32]byte //nolint:unused

	// Cache line 1 (cold/rare, bytes 64..127):
	// - fatalWriter is rare (fatal path). Keep it separate so it doesn't evict hot state.
	fatalWriter io.Writer // 16 (interface: type+data)
}

type ParamsNewLogger struct {
	Ingestor        *bytearena.Ingestor
	WithFatalWriter io.Writer

	EstimatedMessageSizeOverall uint32
	LoggerLevel                 Level

	CallerLevel uint8

	WithCaller bool
	WithColor  bool
	WithJSON   bool
}

// NewLogger comes by default with nano logging.
//
// Use an option for the desired timestamp pattern.
func NewLogger(params *ParamsNewLogger, options ...Option) (*Logger, error) {
	if params.Ingestor == nil {
		return nil,
			errors.New("nil ingestor")
	}

	if params.WithFatalWriter == nil {
		return nil,
			errors.New("nil fatal writer")
	}

	result := Logger{
		withCaller:  params.WithCaller,
		callerLevel: params.CallerLevel,

		fnTimestamp: timestamp.TimestampNano,
		withColor:   params.WithColor,
		withJSON:    params.WithJSON,

		ingestor: params.Ingestor,

		estimatedMessageSizeOverall: params.EstimatedMessageSizeOverall,
		fatalWriter:                 params.WithFatalWriter,
	}

	result.SetLogLevel(params.LoggerLevel)

	if result.estimatedMessageSizeOverall == 0 {
		result.estimatedMessageSizeOverall = MessageSmallSize
	}

	if result.callerLevel == 0 {
		result.callerLevel = 2
	}

	for _, option := range options {
		if option == nil {
			continue
		}

		option(&result)
	}

	result.Printf(
		"created logger, level %v",
		logLevels[params.LoggerLevel],
	)

	return &result,
		nil
}
