package main

import (
        //"go.mongodb.org/mongo-driver/mongo"
		"fmt"
		"log"
		"os"
		"net/http"
		"encoding/json"
		"io/ioutil"
		"bytes"
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
	http.ListenAndServe(":8081", nil) //replace later with a env var
}

func userHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodPost {
		fmt.Println("Usao u post user")
		// Read the request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		// // Unmarshal JSON into a LendingRequest struct
		// var bodyJson LendingRequest
		// err = json.Unmarshal(body, &bodyJson)
		// if err != nil {
		// 	http.Error(w, "Error unmarshaling JSON", http.StatusBadRequest)
		// 	return
		// }

		// newLendingRequest := LendingRequest{
		// 	JMBG:           bodyJson.UserJMBG,
		// 	Title:           bodyJson.Title,
		// 	Author:        bodyJson.Author,
		// 	LendingDate:  bodyJson.LendingDate,
		// }

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
		
		targetURL := "http://localhost:8080/user"

		// Marshal payload to JSON
		payloadBytes, err := json.Marshal(newUser)
		if err != nil {
			fmt.Println("Error marshaling payload:", err)
			return
		}

		// Make a POST request
		resp, err := http.Post(targetURL, "application/json", bytes.NewBuffer(payloadBytes))
		if err != nil {
			fmt.Println("Error making POST request:", err)
			return
		}
		defer resp.Body.Close()

		// Check the response status
		if resp.StatusCode == 501 { // user already exists
			fmt.Println("User with that JMBG already exists: ", newUser.JMBG)
			w.WriteHeader(501);
			return
		} else if resp.StatusCode != http.StatusOK { // server error
			fmt.Println("Unexpected response status: ", resp.Status)
			w.WriteHeader(500);
			return
		} else { // happy ending
			fmt.Println("New user sucessfully created ")
			w.WriteHeader(200);
			return
		}
	}
}

func bookLendingHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodPost {
		fmt.Println("Usao u post lending request")
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

		newLendingRequest := LendingRequest{
			UserJMBG:           bodyJson.UserJMBG,
			Title:           bodyJson.Title,
			Author:        bodyJson.Author,
			ISBN: bodyJson.ISBN,
			LendingDate:  bodyJson.LendingDate,
		}

		targetURL := "http://localhost:8080/user/bookLending"

		// Marshal payload to JSON
		payloadBytes, err := json.Marshal(newLendingRequest)
		if err != nil {
			fmt.Println("Error marshaling payload:", err)
			return
		}

		// Make a POST request
		resp, err := http.Post(targetURL, "application/json", bytes.NewBuffer(payloadBytes))
		if err != nil {
			fmt.Println("Error making POST request:", err)
			return
		}
		defer resp.Body.Close()

		// Check the response status from central server
		if resp.StatusCode == 501 { // user has 3 or more books		
			fmt.Println("User already has 3 or more books, can't take more");
			w.WriteHeader(501);
			return
		} else if resp.StatusCode != http.StatusOK { // server error
			fmt.Println("Unexpected response status: ", resp.Status)
			w.WriteHeader(500);
			return
		} else { // happy ending
			err = repo.CreateLendingRequest(newLendingRequest)
			if err != nil {
				fmt.Println("Error while creating the lending request!")
				w.WriteHeader(500)
				return
			}
			fmt.Println("New lending request sucessfully created ")
			w.WriteHeader(200);
			return
		}
	}
}

