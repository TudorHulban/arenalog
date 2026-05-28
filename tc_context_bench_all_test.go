package arenalog

import (
	"testing"
)

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkAll_OneField/Phuslu_OneField/G1-16      	 9260145	       129.6 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Phuslu_OneField/G2-16      	 9109618	       130.6 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Phuslu_OneField/G3-16      	 9085010	       130.4 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Phuslu_OneField/G4-16      	 9106308	       131.9 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Zerolog_OneField/G1-16     	 7665572	       155.5 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Zerolog_OneField/G2-16     	 7519861	       157.5 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Zerolog_OneField/G3-16     	 7683258	       155.7 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Zerolog_OneField/G4-16     	 7644064	       155.7 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Arenalog_OneField/G1-16    	19064499	        61.81 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Arenalog_OneField/G2-16    	12054794	       101.7 ns/op	       6 B/op	       0 allocs/op
// BenchmarkAll_OneField/Arenalog_OneField/G3-16    	11919625	       100.0 ns/op	       5 B/op	       0 allocs/op
// BenchmarkAll_OneField/Arenalog_OneField/G4-16    	11704431	       100.8 ns/op	       5 B/op	       0 allocs/op

func BenchmarkAll_OneField(b *testing.B) {
	b.Run("Phuslu_OneField", BenchmarkPhuslu_OneField)
	b.Run("Zerolog_OneField", BenchmarkZerolog_OneField)
	b.Run("Arenalog_OneField", BenchmarkArenalog_Msg_OneField)
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkAll_SeveralFields/Phuslu_SeveralFields/gomaxprocs=1-16         	 4880181	       251.9 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Phuslu_SeveralFields/gomaxprocs=2-16         	 9006177	       133.9 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Phuslu_SeveralFields/gomaxprocs=3-16         	13132957	        91.48 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Phuslu_SeveralFields/gomaxprocs=4-16         	17136068	        70.17 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Phuslu_SeveralFields/gomaxprocs=8-16         	28374472	        43.35 ns/op	       0 B/op	       0 allocs/op

// BenchmarkAll_SeveralFields/Zerolog_SeveralFields/gomaxprocs=1-16        	 4400286	       267.7 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Zerolog_SeveralFields/gomaxprocs=2-16        	 8288901	       146.2 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Zerolog_SeveralFields/gomaxprocs=3-16        	12138207	        99.80 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Zerolog_SeveralFields/gomaxprocs=4-16        	16041536	        76.67 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Zerolog_SeveralFields/gomaxprocs=8-16        	27730740	        42.03 ns/op	       0 B/op	       0 allocs/op

// BenchmarkAll_SeveralFields/Arenalog_SeveralFields/gomaxprocs=1-16       	15903200	        74.98 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Arenalog_SeveralFields/gomaxprocs=2-16       	18381885	        73.89 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Arenalog_SeveralFields/gomaxprocs=3-16       	17373678	        72.36 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Arenalog_SeveralFields/gomaxprocs=4-16       	15811662	        76.26 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Arenalog_SeveralFields/gomaxprocs=8-16       	16928438	        73.56 ns/op	       0 B/op	       0 allocs/op

func BenchmarkAll_SeveralFields(b *testing.B) {
	b.Run("Phuslu_SeveralFields", BenchmarkPhuslu_WithFields_Parallel)
	b.Run("Zerolog_SeveralFields", BenchmarkZerolog_WithFields_Parallel)
	b.Run("Arenalog_SeveralFields", BenchmarkArenalog_MultipleFields_Parallel)
}
