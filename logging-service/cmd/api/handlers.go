package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/username/log-service/data"
)

type requestPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(c *gin.Context) {
	var requestPayload requestPayload
	err := app.ReadJson(c, &requestPayload)
	if err != nil {
		fmt.Print("error reading json: ", err)
		return
	}

	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err = app.Models.LogEntry.Insert(event)
	if err != nil {
		fmt.Print("error inserting data: ", err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "Successfully inserted log",
		Data:    event,
	}

	app.WriteJson(c, http.StatusOK, resp)
}
