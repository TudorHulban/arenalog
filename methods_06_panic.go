package arenalog

// terminal event, should not have dependencies like the ingestor.

func (l *Logger) Panic(args ...any) {
	msg := l.format(LevelPanic, args...)

	// We write to the writer before panicking so the log
	// is captured even if the panic is recovered elsewhere.
	_, _ = l.fatalWriter.Write([]byte(msg))

	panic(msg)
}

func (l *Logger) Panicf(format string, args ...any) {
	msg := l.formatf(LevelPanic, format, args...)

	// We write to the writer before panicking so the log
	// is captured even if the panic is recovered elsewhere.
	_, _ = l.fatalWriter.Write([]byte(msg))

	panic(msg)
}

func (l *Logger) Panicw(msg string, keysAndValues ...any) {
	out := l.formatw(LevelPanic, msg, keysAndValues...)

	// We write to the writer before panicking so the log
	// is captured even if the panic is recovered elsewhere.
	_, _ = l.fatalWriter.Write([]byte(msg))

	panic(out)
}
