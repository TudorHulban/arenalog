package arenalog

import "github.com/tudorhulban/arenalog/helpers"

func (l *Logger) labelTrace() string {
	if l.withColor {
		return colorTrace(logLevels[LevelTrace])
	}

	return logLevels[LevelTrace]
}

func (l *Logger) Trace(args ...any) {
	if Level(l.logLevel.Load()) > LevelTrace {
		return
	}

	l.logWithLabel(
		l.labelTrace(),
		helpers.GetEstimatedMessageSize("", args),
		args,
	)
}

func (l *Logger) Tracef(format string, args ...any) {
	if Level(l.logLevel.Load()) > LevelTrace {
		return
	}

	l.logfWithLabel(
		l.labelTrace(),
		format,
		helpers.GetEstimatedMessageSize(format, args),
		args,
	)
}

func (l *Logger) Tracew(msg string, keysAndValues ...any) {
	if Level(l.logLevel.Load()) > LevelTrace {
		return
	}

	if (len(keysAndValues) & 1) != 0 {
		keysAndValues = append(keysAndValues, "(MISSING)")

		msg = "LOG_ERR(odd_args): " + msg
	}

	l.logwWithLabel(
		l.labelTrace(),
		msg,
		uint32(len(msg))+helpers.GetEstimatedMessageSize("", keysAndValues),
		keysAndValues...,
	)
}
