package main

import (
	"log"
	"net/http"

	"fitpass.com/database"
	"fitpass.com/handlers"
	"github.com/gorilla/mux"
)

func main() {

	database.InitDB()

	router := mux.NewRouter()

	// User Routes
	router.HandleFunc("/users", handlers.CreateUser).Methods("POST")

	// Profile Routes
	router.HandleFunc("/profiles", handlers.CreateProfile).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}
