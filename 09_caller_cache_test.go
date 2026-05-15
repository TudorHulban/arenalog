package arenalog

import (
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// callHelper provides a stable, repeatable call site for the cache test
func callHelper(l *Logger, level int) (string, int) {
	return l.getCallerData(level)
}

func TestGetCallerData(t *testing.T) {
	l := Logger{
		callerLevel: 2,
	}

	// 1. Verify File Suffix
	// We check that the file path returned contains the current test file name.
	callingFromFile, _ := l.getCallerData(0)
	require.True(t,
		strings.HasSuffix(callingFromFile, "09_caller_cache_test.go"),

		"expected file to end with logger_test.go, got %q",
		callingFromFile,
	)

	// 2. Test Accuracy
	// We get the line number of the line immediately before our call.
	_, _, expectedLine, _ := runtime.Caller(0)
	_, callingFromLine := l.getCallerData(0) // This is line 'expectedLine + 1'
	require.Equal(t,
		callingFromLine,
		expectedLine+1,

		"expected line %d, got %d",
		expectedLine+1,
		callingFromLine,
	)

	// 3. Test the Cache (Fast Path)
	// Running this in a loop ensures the atomic pointer swap and
	// subsequent loads from the table are stable.
	for ix := range 100 {
		callingFile, callingLine := callHelper(&l, 0)
		if !strings.HasSuffix(callingFile, "09_caller_cache_test.go") || callingLine == 0 {
			t.Fatalf(
				"Cache iteration %d failed: got %s:%d",
				ix,
				callingFile,
				callingLine,
			)
		}
	}
}

func TestCallerLevel(t *testing.T) {
	l := Logger{}

	// Helper function to wrap the call
	wrapper := func(depth int) (string, int) { //nolint:gocritic
		// depth 0 inside here = this anonymous function
		// depth 1 inside here = TestCallerLevel
		return l.getCallerData(depth)
	}

	// We want the anonymous func (depth 0 relative to wrapper's call)
	f0, _ := wrapper(0) // produces /mnt/tmpfs.ramdisk/log/09_caller_cache_test.go, line 68

	// We want TestCallerLevel (depth 1 relative to wrapper's call)
	f1, _ := wrapper(1) // produces /mnt/tmpfs.ramdisk/log/09_caller_cache_test.go, line 75

	require.True(t,
		strings.Contains(f0, "09_caller_cache_test.go"),

		"Level 0 (wrapper) failed, got %s",
		f0,
	)

	require.True(t,
		strings.Contains(f1, "09_caller_cache_test.go"),

		"Level 1 (TestCallerLevel) failed, got %s",
		f1,
	)
}
