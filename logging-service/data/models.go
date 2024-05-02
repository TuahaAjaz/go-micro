package data

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		LogEntry: LogEntry{},
	}
}

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID        string    `bson:"_id:omitempty" json:"id:omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_At" json:"created_At"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (l *LogEntry) Insert(entry LogEntry) error {
	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		fmt.Println("Error inserting into logs: ", err)
		return err
	}

	return nil
}

func (l *LogEntry) All() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", "-1"}})

	records, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		fmt.Println("Error fetching all records: ", err)
		return nil, err
	}

	defer records.Close(ctx)

	var loggedRecords []*LogEntry

	for records.Next(ctx) {
		var item LogEntry

		err := records.Decode(&item)
		if err != nil {
			fmt.Println("Error deccoding Records: ", err)
			return nil, err
		} else {
			loggedRecords = append(loggedRecords, &item)
		}
	}

	return loggedRecords, nil
}

func (l *LogEntry) GetARecord(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Print("Error converting ID: ", docID)
		return nil, err
	}

	var log LogEntry

	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&log)
	if err != nil {
		fmt.Print("Error Fetching a record: ", err)
		return nil, err
	}

	return &log, nil
}

func (l *LogEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	err := collection.Drop(ctx)
	if err != nil {
		fmt.Print("Error dropping Collection: ", err)
		return err
	}
	return nil
}

func (l *LogEntry) UpdateRecord() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	docID, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		fmt.Print("Error converting ID: ", docID)
		return nil, err
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": docID}, bson.D{
		{"$set", bson.D{
			{"name", l.Name},
			{"data", l.Data},
			{"updated_at", time.Now()},
		}},
	},
	)

	if err != nil {
		fmt.Print("Error updating record: ", err)
		return nil, err
	}

	return result, nil
}
