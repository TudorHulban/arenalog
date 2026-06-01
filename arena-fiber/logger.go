package arenafiber

import (
	"context"
	"io"

	fiberlog "github.com/gofiber/fiber/v3/log"
	"github.com/tudorhulban/arenalog"
)

// see https://docs.gofiber.io/api/log#global-log

var _ fiberlog.AllLogger[*arenalog.Logger] = (*Logger)(nil)

type Logger struct {
	L *arenalog.Logger
}

// Note: These methods are lightweight wrappers that directly call the underlying
// arenalog instance. The Go compiler automatically inlines these single-call
// passthrough methods (mid-stack inlining), eliminating function call overhead
// at compile time without requiring manual inline optimization.

// --- Trace ---
func (f *Logger) Trace(args ...any)                 { f.L.Trace(args...) }
func (f *Logger) Tracef(format string, args ...any) { f.L.Tracef(format, args...) }
func (f *Logger) Tracew(msg string, keysAndValues ...any) {
	f.L.Tracew(msg, keysAndValues...)
}

// --- Debug ---
func (f *Logger) Debug(args ...any)                 { f.L.Debug(args...) }
func (f *Logger) Debugf(format string, args ...any) { f.L.Debugf(format, args...) }
func (f *Logger) Debugw(msg string, keysAndValues ...any) {
	f.L.Debugw(msg, keysAndValues...)
}

// --- Info ---
func (f *Logger) Info(args ...any)                 { f.L.Info(args...) }
func (f *Logger) Infof(format string, args ...any) { f.L.Infof(format, args...) }
func (f *Logger) Infow(msg string, keysAndValues ...any) {
	f.L.Infow(msg, keysAndValues...)
}

// --- Warn ---
func (f *Logger) Warn(args ...any)                 { f.L.Warn(args...) }
func (f *Logger) Warnf(format string, args ...any) { f.L.Warnf(format, args...) }
func (f *Logger) Warnw(msg string, keysAndValues ...any) {
	f.L.Warnw(msg, keysAndValues...)
}

// --- Error ---
func (f *Logger) Error(args ...any)                 { f.L.Error(args...) }
func (f *Logger) Errorf(format string, args ...any) { f.L.Errorf(format, args...) }
func (f *Logger) Errorw(msg string, keysAndValues ...any) {
	f.L.Errorw(msg, keysAndValues...)
}

// --- Fatal ---
func (f *Logger) Fatal(args ...any)                 { f.L.Fatal(args...) }
func (f *Logger) Fatalf(format string, args ...any) { f.L.Fatalf(format, args...) }
func (f *Logger) Fatalw(msg string, keysAndValues ...any) {
	f.L.Fatalw(msg, keysAndValues...)
}

// --- Panic ---
func (f *Logger) Panic(args ...any)                 { f.L.Panic(args...) }
func (f *Logger) Panicf(format string, args ...any) { f.L.Panicf(format, args...) }
func (f *Logger) Panicw(msg string, keysAndValues ...any) {
	f.L.Panicw(msg, keysAndValues...)
}

// --- Print ---
func (f *Logger) Print(args ...any)                 { f.L.Print(args...) }
func (f *Logger) Printf(format string, args ...any) { f.L.Printf(format, args...) }
func (f *Logger) Printw(msg string, keysAndValues ...any) {
	f.L.Printw(msg, keysAndValues...)
}

// --- Fiber-required structured logging ---
func (f *Logger) With(args ...any) fiberlog.AllLogger[*arenalog.Logger] {
	// Your logger does not support structured fields natively.
	// No-op is acceptable.
	return f
}

func (f *Logger) WithGroup(name string) fiberlog.AllLogger[*arenalog.Logger] {
	// Same: no-op unless you want grouping.
	return f
}

func (f *Logger) Logger() *arenalog.Logger {
	return f.L
}

func (f *Logger) SetLevel(level fiberlog.Level) {
	switch level { //nolint:revive
	case fiberlog.LevelTrace:
		f.L.SetLogLevel(arenalog.LevelTrace)
	case fiberlog.LevelDebug:
		f.L.SetLogLevel(arenalog.LevelDebug)
	case fiberlog.LevelInfo:
		f.L.SetLogLevel(arenalog.LevelInfo)
	case fiberlog.LevelWarn:
		f.L.SetLogLevel(arenalog.LevelWarn)
	case fiberlog.LevelError:
		f.L.SetLogLevel(arenalog.LevelError)
	case fiberlog.LevelFatal:
		f.L.SetLogLevel(arenalog.LevelFatal)
	case fiberlog.LevelPanic:
		f.L.SetLogLevel(arenalog.LevelPanic)

	default:
		// fallback to Info
		f.L.SetLogLevel(arenalog.LevelInfo)
	}
}

func (*Logger) SetOutput(w io.Writer) {
	// Your logger does not expose SetOutput directly,
	// but the standard pattern is to redirect PrintRaw / PrintMessage
	// through a writer. If your logger has no writer concept,
	// you must store it and ignore it.
	//
	// Minimal no‑op implementation:
	_ = w
}

// should be below, requires changes to the ingestor.
// func (f *FiberLogger) SetOutput(w io.Writer) {
//     if f.L == nil || f.L.ingestor == nil {
//         return
//     }
//     f.L.ingestor.SetWriter(w)
// }

func (f *Logger) WithContext(ctx context.Context) fiberlog.CommonLogger {
	// Your logger does not use context, so this is a no‑op.
	return f
}
