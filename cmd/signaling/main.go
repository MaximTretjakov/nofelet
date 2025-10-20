package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"nofelet/config"
	"nofelet/internal/app/signaling"
	"nofelet/internal/dependency"

	"nofelet/pkg/httpserver"
)

func main() {
	if err := config.New(); err != nil {
		panic(err)
	}
	cfg := config.Current()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	deps, err := dependency.New(&cfg, logger)
	if err != nil {
		log.Fatal(err)
	}

	if err := signaling.New(deps); err != nil {
		log.Fatal(err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	httpServer := httpserver.New(deps.Signaling.Routes,
		httpserver.WithAddress(":"+cfg.HTTP.Port),
		httpserver.WithReadTimeout(cfg.HTTP.ReadTimeout),
		httpserver.WithReadHeaderTimeout(cfg.HTTP.ReadHeaderTimeout),
		httpserver.WithWriteTimeout(cfg.HTTP.WriteTimeout),
		httpserver.WithShutdownTimeout(cfg.HTTP.ShutdownTimeout),
	)

	select {
	case s := <-interrupt:
		logger.Error("signal: " + s.String())
	case err := <-httpServer.Notify():
		logger.Error("httpServer.Notify: %s", err)
	}

	if err := httpServer.Shutdown(); err != nil {
		logger.Error("httpServer.Shutdown: %s", err)
	}
}
