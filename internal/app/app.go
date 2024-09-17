package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	grpcapp "github.com/babs-corp/babs-maps-auth/internal/app/grpc"
	restapp "github.com/babs-corp/babs-maps-auth/internal/app/rest"
	"github.com/babs-corp/babs-maps-auth/internal/services/auth"
	postgres "github.com/babs-corp/babs-maps-auth/internal/storage/pgx"
)

type App struct {
	GRPCSrv *grpcapp.App
	RestSrv *restapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	restPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	storage, err := postgres.New(context.TODO(), storagePath)
	if err != nil {
		panic(fmt.Errorf("cannot init storage: %w", err))
	}

	// TODO: refactor?
	authService := auth.New(log, storage, storage, storage, tokenTTL)

	restApp := restapp.New(log, authService, restPort)
	return &App{
		RestSrv: restApp,
	}
}
