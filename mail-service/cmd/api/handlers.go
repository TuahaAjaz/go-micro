package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *Config) SendMail(c *gin.Context) {
	fmt.Println("Sending Mail...")
	var requestPayload struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	err := app.ReadJson(c, &requestPayload)
	if err != nil {
		fmt.Println("error reading json: ", err)
		app.ErrorJson(c, err)
		return
	}

	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		fmt.Println("error sending mail: ", err)
		app.ErrorJson(c, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Mail sent to %s", requestPayload.To),
	}

	app.WriteJson(c, http.StatusAccepted, payload)
}
