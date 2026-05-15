# Arenalog — The most efficient single‑core logger

Arenalog is a high‑performance structured logger engineered for workloads where
logging must remain predictable, low‑latency, and never become a burden inside
constrained environments.

Benchmark numbers show that normal applications can rely on Arenalog for all their
logging without dedicating more than one core. When running on more than one core the possible slowdown reflects Arenalog’s design choice: it prioritizes single‑core determinism over multi‑core scaling.

Design goals:

- deterministic behavior under load
- minimal branching
- zero allocations in the hot path
- minimal GC interaction
- predictable latency
- optimized for 1‑core execution
- simple ingestion pipeline
- efficient timestamping and field handling

When to choose Arenalog:

- two threads containers where one thread shares the business logic with arenalog
and the other thread can run business logic 100%.
- CPU‑bound or latency‑sensitive systems
- constrained hardware
- environments where multi‑core scaling is irrelevant or undesirable

Examples:

- embedded systems
- edge devices
- proxies and gateways
- serverless cold starts
- real‑time telemetry
- game engines

## Positioning

Arenalog achieves exceptional single-core performance and remains faster than alternatives even with cross-core ingestion, at the cost of a fixed inter-core communication overhead.

## Architecture

Arenalog consists of a formatter and an ingestor.  
The ingestor is the bytearena at https://github.com/TudorHulban/bytearena.  
The formatter strives zero allocations along the way for efficient operations.

## Benchmark Conditions

runtime.GOMAXPROCS() values of 1,2,3 and 4 were used with b.SetParallelism(1) to provide conditions of hosts with limited resources.

Benchmarking loggers in Go is subtle because testing.B.RunParallel does not measure the latency of a single log call. By default, it measures aggregate throughput across all CPU cores.

On a machine with 16 logical CPUs, RunParallel will spawn 16 workers and distribute the work across them.
This makes any logger appear 10–16× faster, even though the logger itself did not improve.

To measure true per‑operation latency, not throughput, the benchmark must remove CPU‑level parallelism as much as possible and exercise the logger under realistic concurrency.

### runtime.GOMAXPROCS(1)

This forces the Go scheduler to run the benchmark on exactly one logical CPU.

All goroutines created by RunParallel will execute on the same logical CPU.
This eliminates the throughput illusion caused by multiple cores dividing the work.

This setting ensures that the benchmark measures:

- the real cost of a log call in a concurrency scenario
- the real timestamp cost
- the real JSON cost
- the real writer cost
- the real branch and pipeline behavior

In other words, it reveals true latency.

### b.SetParallelism(1)

This instructs the benchmark to spawn 1 multiplied by GOMAXPROCS worker goroutines, even though they all run on a single logical CPU.

This is important because it:

- keeps the CPU pipeline hot
- stabilizes branch prediction
- stabilizes timestamp generation
- stabilizes JSON formatting paths
- simulates realistic concurrent logging load

The result is a stable, low‑jitter measurement of the logger’s actual per‑operation cost.

### Combined Effect

Using both settings:

```go
runtime.GOMAXPROCS(1,2,3,4)
b.SetParallelism(1)
```

produces a benchmark configuration that:

- removes multi‑core throughput distortion
- preserves realistic concurrency
- reveals true per‑operation latency
- allows fair comparison between loggers

## Resources

```text
https://dave.cheney.net/2017/01/23/the-package-level-logger-anti-pattern
```
