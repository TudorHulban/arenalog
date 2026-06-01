package arenafiber

import (
	"context"
	"io"

	fiberlog "github.com/gofiber/fiber/v3/log"
	"github.com/tudorhulban/arenalog"
)

// see https://docs.gofiber.io/api/log#global-log

var _ fiberlog.AllLogger[*ALogger] = (*ALogger)(nil)

type ALogger struct {
	L *arenalog.Logger
}

// Note: These methods are lightweight wrappers that directly call the underlying
// arenalog instance. The Go compiler should automatically inline these single-call
// passthrough methods (mid-stack inlining), eliminating function call overhead
// at compile time without requiring manual inline optimization.

// --- Trace ---
func (f *ALogger) Trace(args ...any)                 { f.L.Trace(args...) }
func (f *ALogger) Tracef(format string, args ...any) { f.L.Tracef(format, args...) }
func (f *ALogger) Tracew(msg string, keysAndValues ...any) {
	f.L.Tracew(msg, keysAndValues...)
}

// --- Debug ---
func (f *ALogger) Debug(args ...any)                 { f.L.Debug(args...) }
func (f *ALogger) Debugf(format string, args ...any) { f.L.Debugf(format, args...) }
func (f *ALogger) Debugw(msg string, keysAndValues ...any) {
	f.L.Debugw(msg, keysAndValues...)
}

// --- Info ---
func (f *ALogger) Info(args ...any)                 { f.L.Info(args...) }
func (f *ALogger) Infof(format string, args ...any) { f.L.Infof(format, args...) }
func (f *ALogger) Infow(msg string, keysAndValues ...any) {
	f.L.Infow(msg, keysAndValues...)
}

// --- Warn ---
func (f *ALogger) Warn(args ...any)                 { f.L.Warn(args...) }
func (f *ALogger) Warnf(format string, args ...any) { f.L.Warnf(format, args...) }
func (f *ALogger) Warnw(msg string, keysAndValues ...any) {
	f.L.Warnw(msg, keysAndValues...)
}

// --- Error ---
func (f *ALogger) Error(args ...any)                 { f.L.Error(args...) }
func (f *ALogger) Errorf(format string, args ...any) { f.L.Errorf(format, args...) }
func (f *ALogger) Errorw(msg string, keysAndValues ...any) {
	f.L.Errorw(msg, keysAndValues...)
}

// --- Fatal ---
func (f *ALogger) Fatal(args ...any)                 { f.L.Fatal(args...) }
func (f *ALogger) Fatalf(format string, args ...any) { f.L.Fatalf(format, args...) }
func (f *ALogger) Fatalw(msg string, keysAndValues ...any) {
	f.L.Fatalw(msg, keysAndValues...)
}

// --- Panic ---
func (f *ALogger) Panic(args ...any)                 { f.L.Panic(args...) }
func (f *ALogger) Panicf(format string, args ...any) { f.L.Panicf(format, args...) }
func (f *ALogger) Panicw(msg string, keysAndValues ...any) {
	f.L.Panicw(msg, keysAndValues...)
}

// --- Print ---
func (f *ALogger) Print(args ...any)                 { f.L.Print(args...) }
func (f *ALogger) Printf(format string, args ...any) { f.L.Printf(format, args...) }
func (f *ALogger) Printw(msg string, keysAndValues ...any) {
	f.L.Printw(msg, keysAndValues...)
}

// --- Fiber-required structured logging ---
func (f *ALogger) With(args ...any) fiberlog.AllLogger[*ALogger] {
	return f
}

func (f *ALogger) WithGroup(name string) fiberlog.AllLogger[*ALogger] {
	return f
}

func (f *ALogger) Logger() *ALogger {
	return f
}

func (f *ALogger) SetLevel(level fiberlog.Level) {
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

func (*ALogger) SetOutput(w io.Writer) {
	// Logger does not expose SetOutput directly,
	// but the standard pattern is to redirect PrintRaw / PrintMessage
	// through a writer.
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

func (f *ALogger) WithContext(ctx context.Context) fiberlog.CommonLogger {
	// Logger does not use context, so this is a no‑op.
	return f
}
