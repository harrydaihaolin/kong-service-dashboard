package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

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

// LoggerMiddleware logs details about each HTTP request
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Capture the start time
		startTime := time.Now()

		// Combine the query params into a single URL
		urlWithParams := r.URL.Path + "?" + r.URL.RawQuery
		log.Printf("[Requests] %s %s", r.Method, urlWithParams)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)

		// Log the duration
		duration := time.Since(startTime)
		log.Printf("Completed in %v", duration)
	})
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

func GetServices(w http.ResponseWriter, r *http.Request) {
	db := GetDBInstance()
	var services []Service
	var service Service

	// Get pagination and sorting parameters from query string
	queryParams := r.URL.Query()
	page := queryParams.Get("page")
	limit := queryParams.Get("limit")
	sortBy := queryParams.Get("sort_by")
	order := queryParams.Get("order")
	searchFlag := queryParams.Get("search_mode")
	name := queryParams.Get("name")
	id := queryParams.Get("id")

	// Set default values if parameters are not provided
	if page == "" {
		page = "1"
	}
	if limit == "" {
		limit = "10"
	}
	if sortBy == "" {
		sortBy = "id"
	}
	if order == "" {
		order = "asc"
	}

	// Convert parameters to integers
	pageInt, err := strconv.Atoi(page)
	if handleDBQueryError(w, err, "Invalid page parameter", http.StatusBadRequest) {
		return
	}
	limitInt, err := strconv.Atoi(limit)
	if handleDBQueryError(w, err, "Invalid limit parameter", http.StatusBadRequest) {
		return
	}

	// Validate sorting parameters
	validSortBy := map[string]bool{"id": true, "service_name": true, "created_at": true}
	validOrder := map[string]bool{"asc": true, "desc": true}

	if sortBy != "" && !validSortBy[sortBy] {
		http.Error(w, "Invalid sort_by parameter", http.StatusBadRequest)
		return
	}
	if order != "" && !validOrder[order] {
		http.Error(w, "Invalid order parameter", http.StatusBadRequest)
		return
	}
	if (sortBy == "" && order != "") || (sortBy != "" && order == "") {
		http.Error(w, "Both sort_by and order parameters must be provided together", http.StatusBadRequest)
		return
	}

	// Calculate offset
	offset := (pageInt - 1) * limitInt

	// Fetch data based on search criteria
	switch {
	case id != "":
		// Get service by ID
		fetchAndRespond(w, func() error { return db.First(&service, "id = ?", id).Error }, &service)
	case searchFlag == "true" && name != "":
		// Perform a search by name
		fetchAndRespond(w, func() error {
			return db.Where("service_name LIKE ?", "%"+name+"%").Order(sortBy + " " + order).Limit(limitInt).Find(&services).Error
		}, &services)
	case name != "":
		// Get a single service by name
		fetchAndRespond(w, func() error { return db.First(&service, "service_name = ?", name).Error }, &service)
	default:
		// Fetch paginated and sorted results
		fetchAndRespond(w, func() error {
			return db.Offset(offset).Limit(limitInt).Order(sortBy + " " + order).Find(&services).Error
		}, &services)
	}
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	db := GetDBInstance()
	var users []User
	var user User

	// Get query parameters
	queryParams := r.URL.Query()
	username := queryParams.Get("username")

	// Fetch data based on query parameters
	if username != "" {
		// Get user by username
		fetchAndRespond(w, func() error { return db.First(&user, "user_name = ?", username).Error }, &user)
	} else {
		// Get all users
		fetchAndRespond(w, func() error { return db.Find(&users).Error }, &users)
	}
}
