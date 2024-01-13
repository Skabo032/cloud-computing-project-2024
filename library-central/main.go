package main

import (
        //"context"
        //"go.mongodb.org/mongo-driver/mongo"
        //"go.mongodb.org/mongo-driver/mongo/options"
        //"go.mongodb.org/mongo-driver/mongo/readpref"
		//"go.mongodb.org/mongo-driver/bson"
		//"go.mongodb.org/mongo-driver/bson/primitive"
		"fmt"
		"log"
		"os"
		"net/http"
		//"time"
		//"library-central/repository"
)

func main() {
	http.HandleFunc("/user", userHandler)
	http.HandleFunc("/user/bookLending", bookLendingHandler)
	http.ListenAndServe(8080, nil) //replace later with a env var

	repo, err := NewRepository()
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
		os.Exit(1)
	}
	defer repo.Close()

	// Create
	user := User{
		JMBG:           "1234567890123",
		Name:           "John Doe",
		Address:        "123 Main St",
		NumberOfBooks:  5,
	}
	err = repo.CreateUser(user)
	if err != nil {
		log.Fatal("Error creating user:", err)
	}

	// Read
	users, err := repo.ReadUsers()
	if err != nil {
		log.Fatal("Error reading users:", err)
	}
	fmt.Println("Users:", users)

	// Read by JMBG
	readUser, err := repo.ReadUserByJmbg("1234567890123")
	if err != nil {
		log.Fatal("Error reading user by JMBG:", err)
	}
	fmt.Println("User by JMBG:", readUser)

	// Update
	if readUser != nil {
		err = repo.UpdateUser(readUser.ID, "Updated Name")
		if err != nil {
			log.Fatal("Error updating user:", err)
		}
	}

	// Read after update
	users, err = repo.ReadUsers()
	if err != nil {
		log.Fatal("Error reading users:", err)
	}
	fmt.Println("Users after update:", users)

	// Delete
	if readUser != nil {
		err = repo.DeleteUser(readUser.ID)
		if err != nil {
			log.Fatal("Error deleting user:", err)
		}
	}

	// Read after delete
	users, err = repo.ReadUsers()
	if err != nil {
		log.Fatal("Error reading users:", err)
	}
	fmt.Println("Users after delete:", users)
}

func userHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodGet {
		jmbg = r.FormValue("jmbg")

		user, err := repo.ReadUserByJmbg(jmbg)
		if err != nil {
			log.Fatal("Error reading user by JMBG:", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		fmt.Fprintf(w, "%s,%s,%s,%v\n", user.JMBG, user.Name, user.Address, user.NumberOfBooks)
	}
	else if r.Method == http.MethodPost {
		jmbg := r.FormValue("jmbg")
		name := r.FormValue("name")
		address := r.FormValue("address")

		newUser := User{
			JMBG:           jmbg,
			Name:           name,
			Address:        address,
			NumberOfBooks:  0,
		}

		user, err := repo.ReadUserByJmbg(jmbg)
		if err == mongo.ErrNoDocuments {
			err2 = repo.CreateUser(newUser)
			if err2 != nil {
				log.Fatal("Error creating user:", err)
			}
		}
	}
}