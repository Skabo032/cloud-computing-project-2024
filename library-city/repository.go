package main

import (
	"context"
	"fmt"
	"time"

	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

type Repository struct {
	client *mongo.Client
}

type LendingRequest struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	UserJMBG string `bson:"userJmbg"`
	Title string `bson:"title"`
	Author string `bson:"author"`
	ISBN string `bson:isbn`
	LendingDate string `bson:"lendingDate"`
}

type User struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	JMBG string				`bson:"jmbg"`
	Name string 			`bson:"name"`
	Address string			`bson:"address"`
	NumberOfBooks int    	`bson:"numberOfBooks"`
}

const (
	database        = "db-ns"
	collection      = "lending"
	connectTimeout  = 10 * time.Second
)

func NewRepository() (*Repository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_CONNECTION_STRING")))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB")
	return &Repository{client: client}, nil
}

func (r *Repository) Close() {
	r.client.Disconnect(context.Background())
	fmt.Println("Disconnected from MongoDB")
}

// Create
func (r *Repository) CreateLendingRequest(lendingRequest LendingRequest) error {
	collection := r.client.Database(database).Collection(collection)
	_, err := collection.InsertOne(context.Background(), lendingRequest)
	if err != nil {
		fmt.Println("Error while creating a lending request: ", err)
		return err
	}
	fmt.Println("New lending request created")
	return err
}