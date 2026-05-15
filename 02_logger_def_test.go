package arenalog

import (
	"testing"
	"unsafe"
)

func TestLoggerCacheAlignment(t *testing.T) {
	if unsafe.Offsetof(Logger{}.fatalWriter) != 64 {
		t.Errorf(
			"fatalWriter not aligned to cache line 1: offset=%d",
			unsafe.Offsetof(Logger{}.fatalWriter),
		)
	}
}
