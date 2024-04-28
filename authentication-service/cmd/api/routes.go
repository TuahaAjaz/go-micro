package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (app *Config) routes() http.Handler {
	mux := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://*", "https://*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}
	config.ExposeHeaders = []string{"Link"}
	config.AllowCredentials = true
	config.MaxAge = 300

	mux.Use(cors.Default())

	mux.POST("/authenticate", app.Authenticate)

	return mux
}
