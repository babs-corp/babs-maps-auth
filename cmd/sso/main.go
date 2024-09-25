package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/babs-corp/babs-maps-auth/internal/app"
	"github.com/babs-corp/babs-maps-auth/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting service")

	application := app.New(
		log,
		cfg.Grpc.Port,
		cfg.Rest.Port,
		cfg.StoragePath,
		cfg.TokenTTL,
		cfg.Secret,
	)

	// go application.GRPCSrv.MustRun()
	go application.RestSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	sign := <-stop
	log.Info("shutting down", slog.String("signal", sign.String()))
	application.RestSrv.Stop()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelInfo},
			),
		)
	}

	return log
}
