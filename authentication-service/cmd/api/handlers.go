package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

	fmt.Println("RequestPayload in Auth service => ", requestPayload)

	if err != nil {
		app.ErrorJson(c, err, http.StatusBadRequest)
		return
	}

	user, err := app.Models.User.GetByEmail(requestPayload.Email)

	if err != nil {
		app.ErrorJson(c, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)

	if err != nil || !valid {
		app.ErrorJson(c, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	//log the request
	err = app.logRequest("Log", fmt.Sprintf("%s logged in successfully", requestPayload.Email))
	if err != nil {
		app.ErrorJson(c, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	app.WriteJson(c, http.StatusAccepted, payload)
}

func (app *Config) logRequest(name string, data string) error {
	var logPayload struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	logPayload.Name = name
	logPayload.Data = data

	marshalleLogPayload, _ := json.MarshalIndent(logPayload, "", "\t")
	fmt.Println(marshalleLogPayload)

	request, err := http.NewRequest("POST", "http://logger-service/logs", bytes.NewBuffer(marshalleLogPayload))
	if err != nil {
		return err
	}

	client := &http.Client{}

	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}
