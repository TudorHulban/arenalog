package arenalog

import "github.com/tudorhulban/arenalog/helpers"

func (l *Logger) labelDebug() string {
	if l.withColor {
		return colorDebug(logLevels[LevelDebug])
	}

	return logLevels[LevelDebug]
}

func (l *Logger) Debug(args ...any) {
	if Level(l.logLevel.Load()) > LevelDebug {
		return
	}

	l.logWithLabel(
		l.labelDebug(),
		helpers.GetEstimatedMessageSize("", args),
		args,
	)
}

func (l *Logger) Debugf(format string, args ...any) {
	if Level(l.logLevel.Load()) > LevelDebug {
		return
	}

	l.logfWithLabel(
		l.labelDebug(),
		format,
		helpers.GetEstimatedMessageSize(format, args),
		args,
	)
}

func (l *Logger) Debugw(msg string, keysAndValues ...any) {
	if Level(l.logLevel.Load()) > LevelDebug {
		return
	}

	if (len(keysAndValues) & 1) != 0 {
		keysAndValues = append(keysAndValues, "(MISSING)")

		msg = "LOG_ERR(odd_args): " + msg
	}

	l.logwWithLabel(
		l.labelDebug(),
		msg,
		uint32(len(msg))+helpers.GetEstimatedMessageSize("", keysAndValues),
		keysAndValues...,
	)
}
