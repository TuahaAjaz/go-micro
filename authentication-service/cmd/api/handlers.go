package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (app *Config) InsertUser(c *gin.Context) {
	var requestPayload struct {
		Email     string    `json:"email"`
		FirstName string    `json:"firstName"`
		LastName  string    `json:"lastName"`
		Password  string    `json:"password"`
		Active    int       `json:"active"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	err := app.ReadJson(c, &requestPayload)

	if err != nil {
		app.ErrorJson(c, err, http.StatusBadRequest)
		return
	}

	// user, err := app.Models.User.Insert(requestPayload)

	// if err != nil {
	// 	app.ErrorJson(c, errors.New("Invalid Credentials"), http.StatusBadRequest)
	// 	return
	// }

	// app.WriteJson(c, 200, user)
}

func (app *Config) Authenticate(c *gin.Context) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.ReadJson(c, &requestPayload)

	if err != nil {
		app.ErrorJson(c, err, http.StatusBadRequest)
		return
	}

	user, err := app.Models.User.GetByEmail(requestPayload.Email)

	if err != nil {
		app.ErrorJson(c, errors.New("Invalid Credentials"), http.StatusBadRequest)
		return
	}

	valid, err := app.Models.User.PasswordMatches(requestPayload.Password)

	if err != nil || !valid {
		app.ErrorJson(c, errors.New("Invalid Credentials"), http.StatusBadRequest)
		return
	}

	app.WriteJson(c, 200, user)
}
