package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/tudorhulban/arenalog"
	"github.com/tudorhulban/bytearena"
)

type app struct {
	l *arenalog.Logger
}

// This is your actual logic you want to benchmark
func (a *app) HandleRequest(ctx context.Context) (string, error) {
	a.l.Info("request")

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
