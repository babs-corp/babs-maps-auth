package rest

import (
	"github.com/babs-corp/babs-maps-auth/internal/domain/models"
	"github.com/google/uuid"
)

type LoginInput struct {
	Body struct {
		Email    string `json:"email" doc:"user email"`
		Password string `json:"password" doc:"user password"`
	}
}

type LoginResponse struct {
	Body struct {
		Token string `json:"token" example:"oirijt8u2j3f" doc:"jwt token"`
	}
}

type RegisterInput struct {
	Body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
}

type RegisterResponse struct {
	Body struct {
		Id uuid.UUID `json:"id"`
	}
}

type GetUserInput struct {
	Uid string `doc:"user uid" path:"userId"`
}

type GetUserResponse struct {
	Body struct {
		// TODO: hide password and private data
		User models.User `json:"user" doc:"full user info"`
	}
}

type GetUsersInput struct {
	Limit uint `doc:"limit" query:"limit"`
}

type GetUsersResponse struct {
	Body struct {
		// TODO: hide password and private data
		Users []models.User `json:"users" doc:"full users info"`
	}
}

type GetUserByTokenInput struct {
	Body struct {
		Token    string `json:"token"`
	}
}

type GetUserByTokenResponse struct {
	Body struct {
		// TODO: hide password and private data
		User models.User `json:"user" doc:"full user info"`
	}
}
