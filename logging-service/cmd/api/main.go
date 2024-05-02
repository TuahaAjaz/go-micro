package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/username/log-service/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURI = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	//Connect to Mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic("Error connecting to Mongo", err)
	}
	client = mongoClient

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	go app.service()
}

func (app *Config) service() {
	srv := &http.Server{
		Addr:    fmt.Sprintf("Server starting at port %s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic("Error connecting to server")
	}
}

func connectToMongo() (*mongo.Client, error) {
	//Create connection options
	clientOptions := options.Client().ApplyURI(mongoURI)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	//Connect to mongo
	connection, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("Error coneecting to Mongo: ", err)
		return nil, err
	}

	return connection, err
}
