package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(c *gin.Context) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.ReadJson(c, &requestPayload)

	fmt.Println(requestPayload)

	if err != nil {
		app.ErrorJson(c, errors.New("error reading json"))
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Broker is up",
	}

	_ = app.WriteJson(c, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(c *gin.Context) {
	var requestPayload RequestPayload

	err := app.ReadJson(c, &requestPayload)
	if err != nil {
		app.ErrorJson(c, err)
	}

	fmt.Println("RequestPayload => ", requestPayload)
	fmt.Println("AuthPayload in handleSubmission =>", requestPayload.Auth)

	switch requestPayload.Action {
	case "auth":
		app.authenticate(c, requestPayload.Auth)
	default:
		app.ErrorJson(c, errors.New("unknown action"))
	}
}

func (app *Config) authenticate(c *gin.Context, authPayload AuthPayload) {
	//Create some json  we'll send to the auth service
	fmt.Println("AuthPayload in authenticate =>", authPayload)
	jsonRequest, _ := json.MarshalIndent(authPayload, "", "\t")
	fmt.Println("JSON Request=> ", jsonRequest)

	//call the service
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonRequest))
	if err != nil {
		app.ErrorJson(c, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	fmt.Println("Response => ", response)
	if err != nil {
		app.ErrorJson(c, err)
		return
	}
	defer response.Body.Close()

	//make sure we get back the response
	if response.StatusCode == http.StatusUnauthorized {
		app.ErrorJson(c, errors.New("invalid creds"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.ErrorJson(c, errors.New("some error with auth occured"))
		return
	}

	var authResponse jsonResponse

	err = json.NewDecoder(response.Body).Decode(&authResponse)
	fmt.Println("AuthResponse => ", authResponse)

	if authResponse.Error {
		app.ErrorJson(c, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = authResponse.Data

	app.WriteJson(c, http.StatusAccepted, payload)
}
