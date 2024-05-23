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
	Log    LogPayload  `json:"log"`
	Mail   MailPayload `json:"mail"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
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

	fmt.Println("Log Payload in handleSubmission => ", requestPayload.Log)
	fmt.Println("AuthPayload in handleSubmission =>", requestPayload.Auth)
	fmt.Println("MailPayload in handleSubmission => ", requestPayload.Mail)

	switch requestPayload.Action {
	case "auth":
		app.authenticate(c, requestPayload.Auth)
	case "log":
		app.log(c, requestPayload.Log)
	case "mail":
		app.mail(c, requestPayload.Mail)
	default:
		app.ErrorJson(c, errors.New("unknown action"))
	}
}

func (app *Config) mail(c *gin.Context, mailPayload MailPayload) {
	fmt.Println("Sending mail!!!")

	jsonPayload, err := json.MarshalIndent(mailPayload, "", "\t")
	if err != nil {
		fmt.Println("Error while marshaliing data, ", err)
		app.ErrorJson(c, err)
		return
	}

	request, err := http.NewRequest("POST", "http://mail-service/send", bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error while creating request, ", err)
		app.ErrorJson(c, err)
		return
	}

	client := &http.Client{}

	response, err := client.Do(request)
	fmt.Println("Response => ", response)
	if err != nil {
		fmt.Println("error while creating request, ", err)
		app.ErrorJson(c, err)
		return
	}

	if response.StatusCode != http.StatusAccepted {
		app.ErrorJson(c, errors.New("error in mail service"))
		return
	}

	defer response.Body.Close()

	//unmarshall response and send to client
	var writeResponse jsonResponse
	err = json.NewDecoder(response.Body).Decode(&writeResponse)
	if err != nil {
		fmt.Println(err)
		app.ErrorJson(c, err)
		return
	}

	if writeResponse.Error {
		fmt.Println(writeResponse.Message)
		app.ErrorJson(c, errors.New(writeResponse.Message))
		return
	}

	app.WriteJson(c, http.StatusAccepted, writeResponse)
}

func (app *Config) log(c *gin.Context, logPayload LogPayload) {
	//Marshal payload
	fmt.Println("Unmarshalled Log Payload => ", logPayload)
	fmt.Println(logPayload.Name)
	fmt.Println(logPayload.Data)
	marshalledLogPayload, err := json.MarshalIndent(logPayload, "", "\t")
	if err != nil {
		fmt.Println(err)
		app.ErrorJson(c, err)
		return
	}
	fmt.Println("JSON Request=> ", marshalledLogPayload)

	//Create Request
	request, err := http.NewRequest("POST", "http://logger-service/logs", bytes.NewBuffer(marshalledLogPayload))
	if err != nil {
		fmt.Println(err)
		app.ErrorJson(c, err)
		return
	}

	//initialize Client to make request and send request to log service
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		app.ErrorJson(c, err)
		return
	}
	fmt.Println("Response => ", response)
	defer response.Body.Close()

	//Check for bad request
	if response.StatusCode != http.StatusAccepted {
		fmt.Println(err)
		app.ErrorJson(c, errors.New("error in handling log request"), http.StatusBadRequest)
		return
	}

	//unmarshall response and send to client
	var writeResponse jsonResponse
	err = json.NewDecoder(response.Body).Decode(&writeResponse)
	if err != nil {
		fmt.Println(err)
		app.ErrorJson(c, err)
		return
	}

	if writeResponse.Error {
		fmt.Println(writeResponse.Message)
		app.ErrorJson(c, errors.New(writeResponse.Message))
		return
	}

	app.WriteJson(c, http.StatusAccepted, writeResponse)
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
