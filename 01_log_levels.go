package arenalog

import "fmt"

/*
LOG LEVEL TRUTH TABLE
---------------------

Numeric ordering (lowest → highest severity):
    TRACE(0) < DEBUG(1) < INFO(2) < WARN(3) < ERROR(4) < FATAL(5) < PANIC(6)

Filtering rule:
    A log entry is emitted if entry.level >= logger.threshold

PRINT is outside the severity hierarchy and always emitted.

+----------------+-------+-------+-------+-------+-------+-------+-------+
| Threshold ↓    | TRACE | DEBUG | INFO  | WARN  | ERROR | FATAL | PANIC |
+----------------+-------+-------+-------+-------+-------+-------+-------+
| TRACE (0)      |  YES  |  YES  |  YES  |  YES  |  YES  |  YES  |  YES  |
| DEBUG (1)      |   NO  |  YES  |  YES  |  YES  |  YES  |  YES  |  YES  |
| INFO  (2)      |   NO  |   NO  |  YES  |  YES  |  YES  |  YES  |  YES  |
| WARN  (3)      |   NO  |   NO  |   NO  |  YES  |  YES  |  YES  |  YES  |
| ERROR (4)      |   NO  |   NO  |   NO  |   NO  |  YES  |  YES  |  YES  |
| FATAL (5)      |   NO  |   NO  |   NO  |   NO  |   NO  |  YES  |  YES  |
| PANIC (6)      |   NO  |   NO  |   NO  |   NO  |   NO  |   NO  |  YES  |
+----------------+-------+-------+-------+-------+-------+-------+-------+

Legend:
    YES → log entry is emitted
    NO  → log entry is suppressed
*/

type Level uint8

const (
	LevelTrace Level = 0 // most verbose: all severity levels emitted
	LevelDebug Level = 1 // suppress TRACE
	LevelInfo  Level = 2 // suppress TRACE, DEBUG
	LevelWarn  Level = 3 // suppress TRACE, DEBUG, INFO
	LevelError Level = 4 // suppress TRACE, DEBUG, INFO, WARN
	LevelFatal Level = 5 // suppress TRACE, DEBUG, INFO, WARN, ERROR
	LevelPanic Level = 6 // only PANIC emitted (PRINT always emitted)
)

func (l Level) String() string {
	switch l {
	case LevelTrace:
		return "TRACE"
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	case LevelPanic:
		return "PANIC"

	default:
		return fmt.Sprintf("Level(%d)", l)
	}
}

func (l Level) StringColored() string {
	switch l {
	case LevelTrace:
		return colorTrace("TRACE")
	case LevelDebug:
		return colorDebug("DEBUG")
	case LevelInfo:
		return colorInfo("INFO")
	case LevelWarn:
		return colorWarn("WARN")
	case LevelError:
		return colorError("ERROR")
	case LevelFatal:
		return "FATAL"
	case LevelPanic:
		return "PANIC"

	default:
		return fmt.Sprintf("Level(%d)", l)
	}
}

var logLevels = [7]string{
	"TRACE", // 0
	"DEBUG", // 1
	"INFO",  // 2
	"WARN",  // 3
	"ERROR", // 4
	"FATAL", // 5
	"PANIC", // 6
}

// convertLevel clamps an input level to the valid range [LevelTrace, LevelPanic]
func convertLevel(level Level) Level {
	if level < LevelTrace {
		return LevelTrace
	}

	if level > LevelPanic {
		return LevelPanic
	}

	return level
}
