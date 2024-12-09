package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/gorilla/mux"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	InitDB()

	router := mux.NewRouter()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Println("ok")
	})
	router.HandleFunc("/v1/services", GetAllServices).Methods("GET")
	router.HandleFunc("/v1/services/{id}", GetServiceById).Methods("GET")
	router.HandleFunc("/v1/users", GetAllUsers).Methods("GET")
	router.HandleFunc("/v1/users/{id}", GetUserById).Methods("GET")

	// Wrap the mux with LoggerMiddleware
	loggedMux := LoggerMiddleware(router)

	if err := http.ListenAndServe(":8080", loggedMux); err != nil {
		fmt.Println("Error starting server:", err)
	}
	log.Println("Starting server on :8080")
}
