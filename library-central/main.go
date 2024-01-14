package main

import (
        //"context"
        "go.mongodb.org/mongo-driver/mongo"
        //"go.mongodb.org/mongo-driver/mongo/options"
        //"go.mongodb.org/mongo-driver/mongo/readpref"
		//"go.mongodb.org/mongo-driver/bson"
		//"go.mongodb.org/mongo-driver/bson/primitive"
		"fmt"
		"log"
		"os"
		"net/http"
		"encoding/json"
		"io/ioutil"
)

var repo *Repository

func main() {
	// Starting the mongoDB
	var err error
	repo, err = NewRepository()
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
		os.Exit(1)
	}
	defer repo.Close()

	// HTTP server
	http.HandleFunc("/user", userHandler)
	http.HandleFunc("/user/bookLending", bookLendingHandler)
	http.ListenAndServe(":8080", nil) //replace later with a env var
}

func userHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodGet {
		jmbg := r.FormValue("jmbg")

		user, err := repo.ReadUserByJmbg(jmbg)
		if err != nil {
			log.Println("Error reading user by JMBG:", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		fmt.Fprintf(w, "%s,%s,%s,%v\n", user.JMBG, user.Name, user.Address, user.NumberOfBooks)

		// happy ending
		w.WriteHeader(200);
		return
	} else if r.Method == http.MethodPost {
		// Read the request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		// Unmarshal JSON into a User struct
		var bodyJson User
		err = json.Unmarshal(body, &bodyJson)
		if err != nil {
			http.Error(w, "Error unmarshaling JSON", http.StatusBadRequest)
			return
		}

		newUser := User{
			JMBG:           bodyJson.JMBG,
			Name:           bodyJson.Name,
			Address:        bodyJson.Address,
			NumberOfBooks:  0,
		}

		fmt.Println(bodyJson.JMBG)
		_, err = repo.ReadUserByJmbg(bodyJson.JMBG)
		if err == mongo.ErrNoDocuments {
			err2 := repo.CreateUser(newUser)
			if err2 != nil {
				log.Println("Error creating user:", err)
				http.Error(w, http.StatusText(500), 500)
				return
			}
		} else if err != nil {
			log.Println("Error while reading a user: ", err)
			http.Error(w, http.StatusText(500), 500)
			return
		} else {
			log.Println("User with that JMBG already exists!")
			http.Error(w, http.StatusText(501), 501)
			return
		}

		// happy ending
		w.WriteHeader(200);
		return
	}
}


func bookLendingHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodPost {
		// Read the request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		// Unmarshal JSON into a LendingRequest struct
		var bodyJson LendingRequest
		err = json.Unmarshal(body, &bodyJson)
		if err != nil {
			http.Error(w, "Error unmarshaling JSON", http.StatusBadRequest)
			return
		}

		user, err := repo.ReadUserByJmbg(bodyJson.UserJMBG)
		if err != nil {
			log.Println("Error reading user by JMBG:", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}

		if user.NumberOfBooks < 3 {
			err = repo.IncrementNumOfBooksLent(user.ID)
			if err != nil {
				log.Fatal("Error updating user:", err)
			}
			// happy ending
			w.WriteHeader(200);
			return
		} else {
			log.Println("User already has 3 or more books, can't take more");
			http.Error(w, http.StatusText(501), 501)
			return
		}
	}
}