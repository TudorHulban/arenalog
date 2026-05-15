package arenalog

import (
	"runtime"
	"testing"

	"log"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena/helpers"
)

// BenchmarkStandardLogger-16    	 5514883	       216.9 ns/op	       8 B/op	       0 allocs/op
func BenchmarkStandardLogger(b *testing.B) {
	writer := helpers.CountWriterNoBuffer{}

	log.SetOutput(&writer)
	log.SetFlags(log.LstdFlags)

	runtime.GC()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		log.Printf(
			`{"level":"info","msg":"user login","user_id":%d}`,
			i,
		)
	}

	require.NotZero(b,
		writer.TotalBytesWritten.Load(), // force writer to stay live
	)
}

// BenchmarkStandardLoggerParallel-16    	 8201892	       149.2 ns/op	       8 B/op	       0 allocs/op
func BenchmarkStandardLoggerParallel(b *testing.B) {
	writer := helpers.CountWriterNoBuffer{}

	log.SetOutput(&writer)
	log.SetFlags(log.LstdFlags)

	runtime.GC()

	b.ReportAllocs()
	b.ResetTimer()
	b.SetParallelism(1)

	b.RunParallel(
		func(pb *testing.PB) {
			i := 0

			for pb.Next() {
				log.Printf(
					`{"level":"info","msg":"user login","user_id":%d}`,
					i, // ← variable
				)

				i++
			}
		},
	)

	require.NotZero(b,
		writer.TotalBytesWritten.Load(), // force writer to stay live
	)
}
