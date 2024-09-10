package rest

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/chi/v5"
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
	) (userId int64, err error)
	IsAdmin(ctx context.Context, userId int) (bool, error)
}

const (
	PostRegisterURL = "/register"
	PostLoginURL    = "/login"
	GetIsAdminURL   = "/isAdmin"
)

func InitRoutes(r chi.Router, auth Auth) {
	r.Use(middleware.Logger)

	r.Post(PostRegisterURL, func(w http.ResponseWriter, r *http.Request) {
		handleRegister(w, r, auth)
	})
	r.Post(PostLoginURL, func(w http.ResponseWriter, r *http.Request) {
		handleLogin(w, r, auth)
	})
}
