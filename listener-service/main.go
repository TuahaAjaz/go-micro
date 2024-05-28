package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Panic("error connecting to rabbitmq")
		os.Exit(1)
	}
	defer rabbitConn.Close()

	log.Println("Connecte to rabbitmq successfully")
	// Listen for events

	// Create a consumer

	// Create events on consumed listeners
}

func connect() (*amqp.Connection, error) {
	var count int64 = 1
	var connection *amqp.Connection
	var backoff = 1 * time.Second

	for {
		conn, err := amqp.Dial("amqp://guest:guest@localhost")
		if err != nil {
			count++
			fmt.Println("amqp connection error", err)
		} else {
			connection = conn
			break
		}

		if count > 5 {
			return nil, err
		}

		backoff = time.Duration(math.Pow(float64(count), 2)) * time.Second

		time.Sleep(backoff)
		continue
	}

	return connection, nil
}
