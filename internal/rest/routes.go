package rest

import (
	"context"
	"net/http"

	"github.com/babs-corp/babs-maps-auth/internal/domain/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type Auth interface {
	Login(ctx context.Context,
		email string,
		password string,
		appId int,
	) (token string, err error)
	RegisterNewUser(ctx context.Context,
		email string,
		password string,
	) (userId uuid.UUID, err error)
	IsAdmin(ctx context.Context, userId uuid.UUID) (bool, error)
	UserById(ctx context.Context, userId uuid.UUID) (models.User, error)
}

const (
	PostRegisterURL = "/register"
	PostLoginURL    = "/login"
	GetIsAdminURL   = "/isAdmin"
	GetUserURL      = "/user/{userId}"
)

func InitRoutes(r chi.Router, auth Auth) {
	r.Use(middleware.Logger)

	r.Post(PostRegisterURL, func(w http.ResponseWriter, r *http.Request) {
		handleRegister(w, r, auth)
	})

	r.Post(PostLoginURL, func(w http.ResponseWriter, r *http.Request) {
		handleLogin(w, r, auth)
	})

	r.Get(GetUserURL, func(w http.ResponseWriter, r *http.Request) {
		handleGetUser(w, r, auth)
	})
}
