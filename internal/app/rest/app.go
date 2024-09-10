package rest

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/babs-corp/babs-maps-auth/internal/lib/logger/handlers/sl"
	"github.com/babs-corp/babs-maps-auth/internal/rest"
	"github.com/go-chi/chi/v5"
)

type App struct {
	log    *slog.Logger
	server *http.Server
	port   int
}

func New(
	log *slog.Logger,
	auth rest.Auth,
	port int,
) *App {
	router := chi.NewRouter()
	rest.InitRoutes(router, auth)
	server := &http.Server{
		Addr:    restPort(port),
		Handler: router,
	}

	return &App{
		log:    log,
		server: server,
		port:   port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "restApp.Run"
	log := a.log.With(slog.String("op", op), slog.Int("port", a.port))

	a.server.ListenAndServe()
	log.Info("rest server is running")

	return nil
}

func (a *App) Stop() {
	const op = "restApp.Stop"

	a.log.With(slog.String("op", op)).Info("stopping rest server",
		slog.Int("port", a.port),
	)

	shutdownCtx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := a.server.Shutdown(shutdownCtx)
	if err != nil {
		a.log.Error("failed to shutdown rest server", sl.Err(err))
	}
}

func restPort(port int) string {
	return fmt.Sprintf(":%d", port)
}
