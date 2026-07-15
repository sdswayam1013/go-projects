package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Global MongoDB client.
// This gets initialized from main.go.
var client *mongo.Client

// Models groups all database models together.
type Models struct {
	LogEntry LogEntry
}

// New acts like a constructor.
// Saves MongoDB connection and returns models.
func New(mongoClient *mongo.Client) Models {

	client = mongoClient

	return Models{
		LogEntry: LogEntry{},
	}
}

// LogEntry represents one MongoDB document.
type LogEntry struct {

	// MongoDB document id
	ID string `bson:"_id,omitempty" json:"id,omitempty"`

	// Service generating the log
	Name string `bson:"name" json:"name"`

	// Actual log message
	Data string `bson:"data" json:"data"`

	// Creation timestamp
	CreatedAt time.Time `bson:"created_at" json:"created_at"`

	// Last update timestamp
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

///////////////////////////////////////////////////////////
// INSERT A NEW LOG
///////////////////////////////////////////////////////////

func (l *LogEntry) Insert(entry LogEntry) error {

	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(
		context.TODO(),
		LogEntry{
			Name:      entry.Name,
			Data:      entry.Data,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	)

	if err != nil {
		log.Println("Error inserting into logs:", err)
		return err
	}

	return nil
}

///////////////////////////////////////////////////////////
// GET ALL LOGS
///////////////////////////////////////////////////////////

func (l *LogEntry) All() ([]*LogEntry, error) {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		15*time.Second,
	)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	opts := options.Find()

	// newest logs first
	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(
		context.TODO(),
		bson.D{},
		opts,
	)

	if err != nil {
		log.Println("Finding all docs error:", err)
		return nil, err
	}

	defer cursor.Close(ctx)

	var logs []*LogEntry

	for cursor.Next(ctx) {

		var item LogEntry

		err := cursor.Decode(&item)

		if err != nil {
			log.Println("Error decoding log into slice:", err)
			return nil, err
		}

		logs = append(logs, &item)
	}

	return logs, nil
}

///////////////////////////////////////////////////////////
// GET ONE LOG BY ID
///////////////////////////////////////////////////////////

func (l *LogEntry) GetOne(id string) (*LogEntry, error) {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		15*time.Second,
	)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	// Convert string id to MongoDB ObjectID
	docID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	var entry LogEntry

	err = collection.
		FindOne(ctx, bson.M{"_id": docID}).
		Decode(&entry)

	if err != nil {
		return nil, err
	}

	return &entry, nil
}

///////////////////////////////////////////////////////////
// UPDATE EXISTING LOG
///////////////////////////////////////////////////////////

func (l *LogEntry) Update() (*mongo.UpdateResult, error) {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		15*time.Second,
	)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	// Convert string id to ObjectID
	docID, err := primitive.ObjectIDFromHex(l.ID)

	if err != nil {
		return nil, err
	}

	result, err := collection.UpdateOne(
		ctx,

		// Find document by id
		bson.M{"_id": docID},

		// Fields to update
		bson.D{
			bson.E{
				Key: "$set",
				Value: bson.D{
					bson.E{Key: "name", Value: l.Name},
					bson.E{Key: "data", Value: l.Data},
					bson.E{Key: "updated_at", Value: time.Now()},
				},
			},
		},
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}

///////////////////////////////////////////////////////////
// DELETE ENTIRE COLLECTION
///////////////////////////////////////////////////////////

func (l *LogEntry) DropCollection() error {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		15*time.Second,
	)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	if err := collection.Drop(ctx); err != nil {
		return err
	}

	return nil
}