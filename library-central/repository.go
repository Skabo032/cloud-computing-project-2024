package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	client *mongo.Client
}

type User struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	JMBG string				`bson:"jmbg"`
	Name string 			`bson:"name"`
	Address string			`bson:"address"`
	NumberOfBooks int    	`bson:"numberOfBooks"`
}

type LendingRequest struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	UserJMBG string `bson:"userJmbg"`
	Title string `bson:"title"`
	Author string `bson:"author"`
	ISBN string `bson:isbn`
	LendingDate string `bson:"lendingDate"`
}

const (
	mongoURI        = "mongodb://localhost:27017"
	database        = "testing"
	collection      = "users"
	connectTimeout  = 10 * time.Second
)

func NewRepository() (*Repository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
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
func (r *Repository) CreateUser(user User) error {
	collection := r.client.Database(database).Collection(collection)
	_, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		fmt.Println("Error while creating a user: ", err)
		return err
	}
	fmt.Println("New user created")
	return err
}

// Read
func (r *Repository) ReadUsers() ([]User, error) {
	collection := r.client.Database(database).Collection(collection)
	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []User
	err = cursor.All(context.Background(), &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// ReadUserByJmbg
func (r *Repository) ReadUserByJmbg(jmbg string) (*User, error) {
	collection := r.client.Database(database).Collection(collection)
	filter := bson.M{"jmbg": jmbg}
	var user User
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update
func (r *Repository) UpdateUser(id primitive.ObjectID, newName string) error {
	collection := r.client.Database(database).Collection(collection)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"name": newName}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	return err
}

// Increment num of books lent
func (r *Repository) IncrementNumOfBooksLent(id primitive.ObjectID) error {
	collection := r.client.Database(database).Collection(collection)
	filter := bson.M{"_id": id}
	update := bson.D{{"$inc", bson.D{{"numberOfBooks", 1}}}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	fmt.Println("Number of books incremented for user: ", id)
	return err
}

// Delete
func (r *Repository) DeleteUser(id primitive.ObjectID) error {
	collection := r.client.Database(database).Collection(collection)
	filter := bson.M{"_id": id}
	_, err := collection.DeleteOne(context.Background(), filter)
	return err
}
