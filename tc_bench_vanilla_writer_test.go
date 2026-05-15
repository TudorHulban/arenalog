package arenalog

import (
	"bytes"
	"io"
	"strconv"
	"sync"
	"testing"
	"time"
)

// BenchmarkLoggerZeroAlloc-16    	11793486	       103.0 ns/op	       0 B/op	       0 allocs/op
func BenchmarkLoggerZeroAlloc(b *testing.B) {
	b.ReportAllocs()

	writerDiscard := io.Discard

	pool := sync.Pool{
		New: func() any {
			buf := &bytes.Buffer{}
			buf.Grow(256)

			return buf
		},
	}

	b.ResetTimer()

	for ix := 0; b.Loop(); ix++ {
		buf, _ := pool.Get().(*bytes.Buffer) //nolint:revive
		buf.Reset()

		buf.WriteString(`{"level":"info","ts":"`)
		buf.Write(time.Now().AppendFormat(buf.AvailableBuffer(), time.RFC3339)) // no alloc
		buf.WriteString(`","msg":"user login","user_id":`)
		buf.Write(strconv.AppendInt(buf.AvailableBuffer(), int64(ix), 10)) // no alloc
		buf.WriteByte('}')
		buf.WriteByte('\n')

		_, _ = writerDiscard.Write(buf.Bytes())
		pool.Put(buf)
	}
}
