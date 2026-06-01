package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v3"
	fiberlog "github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/logger"

	"github.com/tudorhulban/arenalog"
	arenafiber "github.com/tudorhulban/arenalog/arena-fiber"
	"github.com/tudorhulban/bytearena"
)

func main() {
	// 1. Writer
	// os.Stdout

	// 2. Ingestor
	ingestor, err := bytearena.NewIngestor(
		bytearena.Size100K(),
		os.Stdout,
	)
	if err != nil {
		fmt.Println(err)

		os.Exit(1)
	}

	// 3. Logger
	l, err := arenalog.NewLogger(
		&arenalog.ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: arenalog.LevelTrace,

			WithFatalWriter: os.Stdout,
			WithCaller:      true,
			WithColor:       true,
			WithJSON:        true,
		},
	)
	if err != nil {
		fmt.Println(err)

		os.Exit(1)
	}

	fiberLogger := arenafiber.ALogger{
		L: l,
	}

	// 4. Start ingestion.
	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	// 5. Register the created logger with Fiber.
	fiberlog.SetLogger(&fiberLogger)

	// 6. Create Fiber app.
	app := fiber.New()

	// 7. Add the logger middleware
	app.Use(
		logger.New(
			logger.Config{
				LoggerFunc: func(c fiber.Ctx, data *logger.Data, cfg *logger.Config) error {
					fiberLogger.L.Info(
						fmt.Sprintf("%s %s %d %s",
							c.Method(),                // GET / POST etc
							c.OriginalURL(),           // full path with query
							data.Stop.Sub(data.Start), // latency
							c.IP(),                    // client IP
						),
					)

					return nil
				},
			},
		),
	)

	// 8. Routes
	app.Get(
		"/",
		func(c fiber.Ctx) error {
			return c.SendString("Hi!")
		},
	)

	// 9. Start server.
	if err := app.Listen(":3001"); err != nil {
		l.Fatal(err)
	}

	// 10. Cleanup
	cancel()
	<-chIngestionEnd
}
