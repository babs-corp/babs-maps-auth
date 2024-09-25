package rest

import "github.com/google/uuid"

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
