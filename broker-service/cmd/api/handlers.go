package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data:omitempty"`
}

func (app *Config) Broker(c *gin.Context) {
	payload := jsonResponse{
		Error:   false,
		Message: "Broker is up",
	}

	out, err := json.MarshalIndent(payload, "", "\t")

	if err != nil {
		panic(err.Error())
	}

	fmt.Println(payload)

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(out)
}
