package timestamp

import (
	"testing"
	"time"
)

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkTimestampYYYYMonth-16    	309270427	         3.902 ns/op	       0 B/op	       0 allocs/op
func BenchmarkTimestampYYYYMonth(b *testing.B) {
	chReady := StartYYYYMonthCache(b.Context())

	<-chReady

	var scratch [64]byte

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		TimestampYYYYMonth(scratch[:0])
	}
}

// BenchmarkTimestampRFC3339-16    	309137002	         3.882 ns/op	       0 B/op	       0 allocs/op
func BenchmarkTimestampRFC3339(b *testing.B) {
	chReady := StartRFC3339UTCCache(b.Context())

	<-chReady

	var scratch [64]byte

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		TimestampRFC3339UTC(scratch[:0])
	}
}

// BenchmarkTimestampRFC3339Bucharest-16    	309029247	         3.880 ns/op	       0 B/op	       0 allocs/op
func BenchmarkTimestampRFC3339Bucharest(b *testing.B) {
	chReady := StartRFC3339BucharestCache(b.Context())

	<-chReady

	var scratch [64]byte

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		TimestampRFC3339Bucharest(scratch[:0])
	}
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkTimenow-16    	28263483	        41.83 ns/op	       0 B/op	       0 allocs/op
func BenchmarkTimenow(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		time.Now().UnixNano()
	}
}
