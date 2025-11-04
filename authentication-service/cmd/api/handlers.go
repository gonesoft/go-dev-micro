package main

import (
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(W http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJson(W, r, &requestPayload)
	if err != nil {
		app.errorJson(W, http.ErrBodyNotAllowed)
		return
	}

	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJson(W, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJson(W, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Authenticated user %s", user.Email),
		Data:    user,
	}

	app.writeJson(W, http.StatusOK, payload)

}
