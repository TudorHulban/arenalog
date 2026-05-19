package main

import (
	"context"
	"os"
	"runtime"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/tudorhulban/arenalog"
	"github.com/tudorhulban/bytearena"
)

// docker run --rm --memory="128m" -p 8080:8080 go-lambda-bench
// 18 May 2026 12:59:36,294 [INFO] (rapid) exec '/var/runtime/bootstrap' (cwd=/var/task, handler=)
// START RequestId: 8bc6364e-a7ed-4792-bb5f-b02ec4402732 Version: $LATEST
// 18 May 2026 12:59:41,476 [INFO] (rapid) INIT START(type: on-demand, phase: init)
// 18 May 2026 12:59:41,476 [INFO] (rapid) The extension's directory "/opt/extensions" does not exist, assuming no extensions to be loaded.
// 18 May 2026 12:59:41,476 [INFO] (rapid) Starting runtime without AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN , Expected?: false
// 18 May 2026 12:59:41,478 [INFO] (rapid) INIT RTDONE(status: success)
// 18 May 2026 12:59:41,478 [INFO] (rapid) INIT REPORT(durationMs: 2.472000)
// 18 May 2026 12:59:41,479 [INFO] (rapid) INVOKE START(requestId: 8bc6364e-a7ed-4792-bb5f-b02ec4402732)
// 18 May 2026 12:59:41,479 [INFO] (rapid) INVOKE RTDONE(status: success, produced bytes: 0, duration: 0.608000ms)
// END RequestId: 8bc6364e-a7ed-4792-bb5f-b02ec4402732
// REPORT RequestId: 8bc6364e-a7ed-4792-bb5f-b02ec4402732  Init Duration: 0.09 ms  Duration: 3.45 ms       Billed Duration: 4 ms   Memory Size: 3008 MB    Max Memory Used: 3008 MB
// {"ts":"2026-05-18T12:59:41.478Z","level":"INFO","msg":"created logger, level INFO"}
// {"ts":"2026-05-18T12:59:41.479Z","level":"INFO","msg":"request"}
// {"ts":"2026-05-18T12:59:41.479Z","level":"INFO","msg":"memory stats - Alloc MB:  0  Sys MB:  7"}

type app struct {
	l *arenalog.Logger
}

// This is your actual logic you want to benchmark
func (a *app) HandleRequest(ctx context.Context) (string, error) {
	a.l.Info("request")

	// Read live memory allocator statistics
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// m.Sys is the total memory obtained from the OS
	// m.Alloc is the bytes allocated and still in use on the heap
	a.l.Info("memory stats - Alloc MB: ", m.Alloc/1024/1024, " Sys MB: ", m.Sys/1024/1024)

	return "Success", nil
}

func main() {
	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		os.Stdout,
	)
	if errCrIngestor != nil {
		os.Exit(2)
	}

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	defer func() {
		cancel()
		<-chIngestionEnd
	}()

	logger, errCrLogger := arenalog.NewLogger(
		&arenalog.ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: arenalog.LevelInfo,

			WithFatalWriter: os.Stdout,
			WithJSON:        true,
		},

		arenalog.WithTimestampRFC3339UTC(context.Background()),
	)
	if errCrLogger != nil {
		os.Exit(3)
	}

	a := app{
		l: logger,
	}

	// This starts a local network listener that keeps the process alive
	lambda.Start(a.HandleRequest)
}
