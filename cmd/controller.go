package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func setJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func handleDBQueryError(w http.ResponseWriter, err error, message string, statusCode int) bool {
	if err != nil {
		http.Error(w, message, statusCode)
		log.Printf("Query error: %v", err)
		return true
	}
	return false
}

func handleEncodingError(w http.ResponseWriter, err error) bool {
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Encoding error: %v", err)
		return true
	}
	return false
}

func handleDataNotFound(w http.ResponseWriter) {
	http.Error(w, "Data not found", http.StatusNotFound)
}

// fetchAndRespond is a helper function that handles fetching data and responding to an HTTP request.
//
// It calls the provided fetch function, checks for errors, and writes the appropriate HTTP response.
//
// Args:
//
//	w (http.ResponseWriter): The HTTP response writer.
//	fetchFunc (func() error): A function that fetches data and returns an error if something goes wrong.
//	data (interface{}): The data to be encoded and sent in the HTTP response if the fetch is successful.
//
// The function performs the following steps:
// 1. Calls the fetch function.
// 2. Checks if the error returned by the fetch function is a "record not found" error and responds with a 404 status code if true.
// 3. Handles any other errors by responding with a 500 status code and logging the error.
// 4. If no errors occur, it sets the response content type to "application/json" and encodes the provided data into the response body.
func fetchAndRespond(w http.ResponseWriter, fetchFunc func() error, data interface{}) {
	// Call the fetch function
	err := fetchFunc()

	// Check for "record not found" error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, "Resource not found", http.StatusNotFound)
		return
	}

	// Handle other potential errors
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Database error: %v", err)
		return
	}

	// Respond with the fetched data encoded as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Encoding error: %v", err)
	}
}

func GetAllServices(w http.ResponseWriter, r *http.Request) {
	db := GetDBInstance()
	var services []Service
	fetchAndRespond(w, func() error { return db.Find(&services).Error }, &services)
}

func GetServiceById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	db := GetDBInstance()
	var service Service
	fetchAndRespond(w, func() error { return db.First(&service, "id = ?", params["id"]).Error }, &service)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	db := GetDBInstance()
	var users []User
	fetchAndRespond(w, func() error { return db.Find(&users).Error }, &users)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	db := GetDBInstance()
	var user User
	fetchAndRespond(w, func() error { return db.First(&user, "id = ?", params["id"]).Error }, &user)
}
