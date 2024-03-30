package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *Config) Broker(c *gin.Context) {
	payload := jsonResponse{
		Error:   false,
		Message: "Broker is up",
	}

	_ = app.WriteJson(c, http.StatusOK, payload)
}
