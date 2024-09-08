package app

import (
	"fmt"
	"log/slog"
	"time"

	grpcapp "github.com/babs-corp/babs-maps-auth/internal/app/grpc"
	"github.com/babs-corp/babs-maps-auth/internal/services/auth"
	"github.com/babs-corp/babs-maps-auth/internal/storage/sqlite"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(fmt.Errorf("cannot init storage: %w", err))
	}

	// TODO: refactor?
	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)
	return &App{
		GRPCSrv: grpcApp,
	}
}
