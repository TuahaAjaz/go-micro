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

func (app *Config) ReadJson(c *gin.Context, data any) error {
	// Declare a variable to store the JSON data as a raw message
	var rawData json.RawMessage

	// Decode the JSON data from the request body
	err := json.NewDecoder(c.Request.Body).Decode(&rawData)
	if err != nil {
		// Handle error if decoding JSON data fails
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	return nil
}

func (app *Config) WriteJson(c *gin.Context, status int, data any, headers ...map[string]string) error {
	out, err := json.Marshal(data)

	fmt.Println(data)

	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			c.Writer.Header().Set(key, value)
		}
	}

	c.Writer.Header().Set("Content-Type", "application/json")

	c.Writer.WriteHeader(status)

	_, err = c.Writer.Write(out)

	if err != nil {
		return err
	}

	return nil
}

func (app *Config) ErrorJson(c *gin.Context, err error, status ...int) {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var res jsonResponse

	res.Error = true
	res.Message = err.Error()

	app.WriteJson(c, statusCode, res)
}
