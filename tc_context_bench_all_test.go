package arenalog

import (
	"testing"
)

// cpu: AMD Ryzen 5 5600U with Radeon Graphics
// BenchmarkAll_OneField/Phuslu_OneField/G1-12      	 8492310	       141.7 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Phuslu_OneField/G2-12      	 8619159	       139.5 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Phuslu_OneField/G3-12      	 8662359	       138.6 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Phuslu_OneField/G4-12      	 8657698	       138.3 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Zerolog_OneField/G1-12     	 7613612	       156.3 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Zerolog_OneField/G2-12     	 7613492	       157.3 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Zerolog_OneField/G3-12     	 7650607	       157.1 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Zerolog_OneField/G4-12     	 7614378	       157.7 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Arenalog_OneField/G1-12    	18319458	        65.20 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Arenalog_OneField/G2-12    	11892158	       100.7 ns/op	       6 B/op	       0 allocs/op
// BenchmarkAll_OneField/Arenalog_OneField/G3-12    	11452639	       103.6 ns/op	       6 B/op	       0 allocs/op
// BenchmarkAll_OneField/Arenalog_OneField/G4-12    	11748916	       102.1 ns/op	       6 B/op	       0 allocs/op

func BenchmarkAll_OneField(b *testing.B) {
	b.Run("Phuslu_OneField", BenchmarkPhuslu_OneField)
	b.Run("Zerolog_OneField", BenchmarkZerolog_OneField)
	b.Run("Arenalog_OneField", BenchmarkArenalog_OneField)
}

// cpu: AMD Ryzen 5 5600U with Radeon Graphics
// BenchmarkAll_SeveralFields/Phuslu_SeveralFields/gomaxprocs=1-12         	 4683182	       254.5 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Phuslu_SeveralFields/gomaxprocs=2-12         	 8703964	       138.6 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Phuslu_SeveralFields/gomaxprocs=3-12         	12884524	        93.39 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Phuslu_SeveralFields/gomaxprocs=4-12         	16642962	        71.29 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Zerolog_SeveralFields/gomaxprocs=1-12        	 4143967	       287.0 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Zerolog_SeveralFields/gomaxprocs=2-12        	 8196825	       147.9 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Zerolog_SeveralFields/gomaxprocs=3-12        	11779630	       101.5 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Zerolog_SeveralFields/gomaxprocs=4-12        	15411469	        78.31 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Arenalog_SeveralFields/gomaxprocs=1-12       	15539617	        80.51 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Arenalog_SeveralFields/gomaxprocs=2-12       	17391261	        67.75 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Arenalog_SeveralFields/gomaxprocs=3-12       	14691753	        79.73 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Arenalog_SeveralFields/gomaxprocs=4-12       	14577597	        85.60 ns/op	       0 B/op	       0 allocs/op

func BenchmarkAll_SeveralFields(b *testing.B) {
	b.Run("Phuslu_SeveralFields", BenchmarkPhuslu_WithFields_Parallel)
	b.Run("Zerolog_SeveralFields", BenchmarkZerolog_WithFields_Parallel)
	b.Run("Arenalog_SeveralFields", BenchmarkArenalog_MultipleFields_Parallel)
}
