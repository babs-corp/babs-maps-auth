package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/babs-corp/babs-maps-auth/internal/domain/models"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Auth interface {
	Login(ctx context.Context,
		email string,
		password string,
	) (token string, err error)
	RegisterNewUser(ctx context.Context,
		email string,
		password string,
	) (userId uuid.UUID, err error)
	IsAdmin(ctx context.Context, userId uuid.UUID) (bool, error)
	UserById(ctx context.Context, userId uuid.UUID) (models.User, error)
	Users(ctx context.Context, limit uint) ([]models.User, error)
}

const (
	PostRegisterURL = "/register"
	PostLoginURL    = "/login"
	GetIsAdminURL   = "/isAdmin"
	GetUserURL      = "/user/{userId}"
	GetUsersURL     = "/users"
)

func InitRoutes(router *chi.Mux, auth Auth) {

	router.Use(middleware.Logger)

	api := humachi.New(router, huma.DefaultConfig("My API", "1.0.0"))

	huma.Register(api, huma.Operation{
		OperationID:   "register-user",
		Method:        http.MethodPost,
		Path:          PostRegisterURL,
		Summary:       "Register new user",
		Tags:          []string{"auth"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, input *RegisterInput) (*RegisterResponse, error) {
		id, err := auth.RegisterNewUser(context.Background(), input.Body.Email, input.Body.Password)
		if err != nil {
			return nil, fmt.Errorf("cannot create user: %w", err)
		}
		resp := RegisterResponse{}
		resp.Body.Id = id
		return &resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "login-user",
		Method:        http.MethodPost,
		Path:          PostLoginURL,
		Summary:       "Login user",
		Tags:          []string{"auth"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, input *LoginInput) (*LoginResponse, error) {
		token, err := auth.Login(context.Background(), input.Body.Email, input.Body.Password)
		if err != nil {
			return nil, fmt.Errorf("cannot login user: %w", err)
		}
		resp := LoginResponse{}
		resp.Body.Token = token
		return &resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "get-user",
		Method:        http.MethodGet,
		Path:          GetUserURL,
		Summary:       "Get user by uuid",
		Tags:          []string{"users"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, input *GetUserInput) (*GetUserResponse, error) {
		rawUserID := input.Uid

		userID, err := uuid.Parse(rawUserID)
		if err != nil {
			return nil, fmt.Errorf("cannot parse user id: %w", err)
		}

		user, err := auth.UserById(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("cannot get user: %w", err)
		}

		resp := GetUserResponse{}
		resp.Body.User = user
		return &resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "get-users",
		Method:        http.MethodGet,
		Path:          GetUsersURL,
		Summary:       "Get users",
		Tags:          []string{"users"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, input *GetUsersInput) (*GetUsersResponse, error) {
		limit := input.Limit
		if limit == 0 {
			limit = 10
		}
		users, err := auth.Users(ctx, limit)
		if err != nil {
			return nil, fmt.Errorf("cannot get users: %w", err)
		}
		resp := GetUsersResponse{}
		resp.Body.Users = users
		return &resp, nil
	})
}
