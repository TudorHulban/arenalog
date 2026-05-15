package arenalog

import (
	"os"
)

// terminal event, should not have dependencies like the ingestor.

func (l *Logger) Fatal(args ...any) {
	_, _ = l.fatalWriter.Write(
		[]byte(
			l.format(LevelFatal, args...),
		),
	)

	os.Exit(1)
}

func (l *Logger) Fatalf(format string, args ...any) {
	_, _ = l.fatalWriter.Write(
		[]byte(
			l.formatf(LevelFatal, format, args...),
		),
	)

	os.Exit(1)
}

func (l *Logger) Fatalw(msg string, keysAndValues ...any) {
	_, _ = l.fatalWriter.Write(
		[]byte(
			l.formatw(LevelFatal, msg, keysAndValues...),
		),
	)

	os.Exit(1)
}
