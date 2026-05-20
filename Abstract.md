# Building an efficient logger in Go for hosts with restricted resources

*Arenalog: Async Go logging optimized for 1–2 vCPU containers and serverless runtimes.*

---

## The Problem: Current loggers are not efficient in hosts with up to two cores

Traditional Go loggers are engineered for macro-scale vertical
concurrency, accepting a baseline structural overhead to survive multi-core environments.  
But what if your service runs in a 1–2 vCPU container? Then aggregate throughput becomes irrelevant. What matters is: how many cycles does one log call take from your business logic?

## Design Philosophy

Arenalog completely flips this paradigm. It is purposely designed from the ground up
to be transparent in high-density, low-core cloud architectures: Kubernetes Sidecars,
Microservices, Serverless environments like AWS Lambda or Cloud Run.  

By hyper-optimizing for a single core, Arenalog achieves a paradigm shift:
Delivering 3-core generic logger throughput on a 1-core cloud budget. On single core, the runtime completely avoids cross-core cache-line bouncing and atomic synchronization thrashing helping the bytearena ingestor achieve its most efficient ingest conditions.  
Arenalog operates asynchronously as long as the underlying writer writes.

## Building blocks

Arenalog is a message formatter and ingestor writing to an io.Writer.  
The ingestor design was covered in https://dev.to/tudorhulban/building-a-log-ingestor-in-go-with-double-buffered-sharded-arenas-48n7.  
The ingestor is built using atomic types which hit peak performance when bound to 1 core because the CPU doesn't have to constantly broadcast and synchronize cache-line mutations across multiple physical cores (eliminating L1/L2 cache thrashing).

The formatter part strives to keep heap allocations at a minimum.

## How to use

First step is to initialize the ingestor.  
By default the ingestor is optimized for one core hosts (ex. `docker run --rm --memory="128m" --cpuset-cpus="0" -p 8080:8080 cname`). To use with multicore hosts the `bytearena.WithCounterCoreCPU()` should be used as in below:

```go
ingestor, errCrIngestor := bytearena.NewIngestor(
			bytearena.Size100K(),// ingestor uses two arenas of size
			&writer,

			helpers.TernaryWithValueIn(
				[]int{1},        // if NumCPU() == 1
				runtime.NumCPU(),// condition
				nil,             // use default
				bytearena.WithCounterCoreCPU(),
			),
		)
```

where the helper is:

```go
func TernaryWithValueIn[V comparable, T any](haystack []V, needle V, value1, value2 T) T {
	if slices.Contains(haystack, needle) {
		return value1
	}

	return value2
}
```

Ingestion should be started:

```go
ctx, cancel := context.WithCancel(context.Background())
chIngestionEnd := ingestor.StartIngestion(ctx)
```

We can use the ingestor in the logger creation:

```go
logger, errCrLogger := arenalog.NewLogger(
		&arenalog.ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: arenalog.LevelInfo,

			WithFatalWriter: os.Stdout,
			WithJSON:        true,
		},

		arenalog.WithTimestampRFC3339UTC(t.Context()),
	)
```

Once we have a logger we can start using it directly:

```go
logger.Info("logger ready")
```

Or create a context root that would be shared among entries:

```go
	logContext := arenalog.NewLogContext(logger).
		WithRoot("service", "auth")

    entry := logContext.WithString("area", "some area")
	entry.Info().Msg("benchmark test")
```

And that should produce:

```json
{"ts":"2026-05-19T13:56:10.288Z","level":"INFO","msg":"created logger, level INFO"}
{"ts":"2026-05-19T13:56:10.288Z","level":"INFO","msg":"logger ready"}
{"ts":"2026-05-19T13:56:10.288Z","level":"INFO","service":"auth","area":"some area","msg":"benchmark test"}
```

The first line is emitted automatically by NewLogger() as an initialization confirmation. Subsequent lines come from explicit logger.Info() calls.

## Measured Performance

A log context with multiple fields as

```go
logContext := NewLogContext(logger).
			WithRoot("service", "auth").
			SetInt("req_id", 12345).
			SetBool("cache_hit", true)

	entry := logContext.
								WithString("area", "some area").
								Info().
								WithString("user", "tudor").
								WithInt("attempt", i).
								WithFloat("some float", 1.1137).
								WithBool("success", true)

							entry.Msg("benchmark test")
```

Is scoring as below:

```text
// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkArenalog_MultipleFields_Parallel/gomaxprocs=1-16         	15310951	        76.55 ns/op	       0 B/op	       0 allocs/op
// BenchmarkArenalog_MultipleFields_Parallel/gomaxprocs=2-16         	17851222	        66.67 ns/op	       0 B/op	       0 allocs/op
// BenchmarkArenalog_MultipleFields_Parallel/gomaxprocs=3-16         	16263550	        74.16 ns/op	       0 B/op	       0 allocs/op
// BenchmarkArenalog_MultipleFields_Parallel/gomaxprocs=4-16         	15345500	        76.67 ns/op	       0 B/op	       0 allocs/op
// BenchmarkArenalog_MultipleFields_Parallel/gomaxprocs=8-16         	16770168	        73.00 ns/op	       0 B/op	       0 allocs/op
```

## When NOT to Use Arenalog

Arenalog prioritizes per-call latency on 1–2 core hosts. If you are running on 8+ core machines and need maximum aggregate throughput, a batched, multi-worker logger may still be more efficient.  
Note that logs may be dropped under backpressure. Backpressure could be measured for specific cases, but on modern CPUs, a safe harbour value should be above one million logs per second per core if the arena size corresponds to the payloads velocity.

## Source

The source is available on GitHub at [github.com/TudorHulban/arenalog](https://github.com/TudorHulban/arenalog). Benchmarks, tests and helpers are all included.
