package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type RegisterRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	AppId    int    `json:"app_id"`
}

func handleRegister(w http.ResponseWriter, r *http.Request, a Auth) {
	reqBody := RegisterRequestBody{}
	if err := render.DecodeJSON(r.Body, &reqBody); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.PlainText(w, r, err.Error())
		return
	}

	id, err := a.RegisterNewUser(r.Context(), reqBody.Email, reqBody.Password)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.PlainText(w, r, err.Error())
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, id)
}

func handleLogin(w http.ResponseWriter, r *http.Request, a Auth) {
	reqBody := LoginRequestBody{}
	if err := render.DecodeJSON(r.Body, &reqBody); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.PlainText(w, r, err.Error())
		return
	}

	token, err := a.Login(r.Context(), reqBody.Email, reqBody.Password, reqBody.AppId)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.PlainText(w, r, err.Error())
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, token)
}

func handleIsAdmin(w http.ResponseWriter, _ *http.Request, _ Auth) {
	w.WriteHeader(http.StatusOK)
}

func handleGetUser(w http.ResponseWriter, r *http.Request, a Auth) {
	rawUserID := chi.URLParam(r, "userId")

	userID, err := uuid.Parse(rawUserID)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.PlainText(w, r, err.Error())
		return
	}

	user, err := a.UserById(r.Context(), userID)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.PlainText(w, r, err.Error())
		return
	}

	render.JSON(w, r, user)
	w.WriteHeader(http.StatusOK)
}
