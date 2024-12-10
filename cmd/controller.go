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
	loadVersion := queryParams.Get("load_version")

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
		fetchAndRespond(w, func() error {
			query := db.First(&service, "id = ?", id)
			if loadVersion == "true" {
				query = query.Preload("Versions")
			}
			return query.Error
		}, &service)
	case searchFlag == "true" && name != "":
		// Perform a search by name
		fetchAndRespond(w, func() error {
			query := db.Where("service_name LIKE ?", "%"+name+"%").Order(sortBy + " " + order).Limit(limitInt)
			if loadVersion == "true" {
				query = query.Preload("Versions")
			}
			return query.Find(&services).Error
		}, &services)
	case name != "":
		// Get a single service by name
		fetchAndRespond(w, func() error {
			query := db.First(&service, "service_name = ?", name)
			if loadVersion == "true" {
				query = query.Preload("Versions")
			}
			return query.Error
		}, &service)
	default:
		// Fetch paginated and sorted results
		fetchAndRespond(w, func() error {
			query := db.Offset(offset).Limit(limitInt).Order(sortBy + " " + order)
			if loadVersion == "true" {
				query = query.Preload("Versions")
			}
			return query.Find(&services).Error
		}, &services)
	}
}

// UpdateService updates an existing service in the database based on the provided JSON payload.
func UpdateService(w http.ResponseWriter, r *http.Request) {
	db := GetDBInstance()

	var service Service
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Ensure ID is provided
	if service.ID == 0 {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	// Check if the service exists
	var existingService Service
	if db.First(&existingService, service.ID).Error != nil {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	if err := db.Save(&service).Error; err != nil {
		http.Error(w, "Failed to update service", http.StatusInternalServerError)
		log.Printf("Error updating service: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(service)
}

func CreateServiceVersion(w http.ResponseWriter, r *http.Request) {
	db := GetDBInstance()

	var version ServiceVersion
	if err := json.NewDecoder(r.Body).Decode(&version); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Validate that ServiceID is provided
	if version.ServiceID == 0 {
		http.Error(w, "ServiceID is required", http.StatusBadRequest)
		return
	}

	// Check if the service exists
	var service Service
	if err := db.First(&service, version.ServiceID).Error; err != nil {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	// Check if the version already exists
	var existingVersion ServiceVersion
	if err := db.Where("service_version_name = ? AND service_id = ?", version.ServiceVersionName, version.ServiceID).First(&existingVersion).Error; err == nil {
		http.Error(w, "Version already exists", http.StatusConflict)
		return
	}

	// Create the new version
	if err := db.Create(&version).Error; err != nil {
		http.Error(w, "Failed to create version", http.StatusInternalServerError)
		log.Printf("Error creating version: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(version)
}

func UpdateServiceVersion(w http.ResponseWriter, r *http.Request) {
	db := GetDBInstance()

	var version ServiceVersion
	if err := json.NewDecoder(r.Body).Decode(&version); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Ensure ID is provided
	if version.ID == 0 {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	// Check if the version exists
	var existingVersion ServiceVersion
	if err := db.First(&existingVersion, version.ID).Error; err != nil {
		http.Error(w, "Version not found", http.StatusNotFound)
		return
	}

	// Update the version
	existingVersion.ServiceVersionName = version.ServiceVersionName
	existingVersion.ServiceVersionDescription = version.ServiceVersionDescription
	existingVersion.ServiceVersionURL = version.ServiceVersionURL

	if err := db.Save(&existingVersion).Error; err != nil {
		http.Error(w, "Failed to update version", http.StatusInternalServerError)
		log.Printf("Error updating version: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingVersion)
}

func DeleteServiceVersion(w http.ResponseWriter, r *http.Request) {
	db := GetDBInstance()

	// Parse query parameters
	id := r.URL.Query().Get("id")

	// Convert ID to integer if provided
	var idInt uint
	if id != "" {
		parsedID, err := strconv.ParseUint(id, 10, 32)
		idInt = uint(parsedID)
		if err != nil {
			http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
			return
		}
	}

	// Check if the version exists by ID
	var version ServiceVersion
	if id != "" {
		if err := db.First(&version, idInt).Error; err != nil {
			http.Error(w, "Resource not found", http.StatusNotFound)
			return
		}
	} else {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	// Perform delete
	if err := db.Delete(&version).Error; err != nil {
		http.Error(w, "Failed to delete version", http.StatusInternalServerError)
		log.Printf("Error deleting version: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// CreateService creates a new service in the database based on the provided JSON payload.
func CreateService(w http.ResponseWriter, r *http.Request) {
	db := GetDBInstance()

	var service Service
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Check if the service already exists
	var existingService Service
	if db.Where("service_name = ?", service.ServiceName).First(&existingService).Error == nil {
		http.Error(w, "Service already exists", http.StatusConflict)
		return
	}

	if err := db.Create(&service).Error; err != nil {
		http.Error(w, "Failed to create service", http.StatusInternalServerError)
		log.Printf("Error creating service: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(service)
}

func DeleteService(w http.ResponseWriter, r *http.Request) {
	db := GetDBInstance()

	// Parse query parameters
	id := r.URL.Query().Get("id")
	name := r.URL.Query().Get("name")

	// Convert ID to integer if provided
	var idInt int
	var err error
	if id != "" {
		idInt, err = strconv.Atoi(id)
		if handleDBQueryError(w, err, "Invalid ID parameter", http.StatusBadRequest) {
			return
		}
	}

	// Check if the service exists by ID or name
	var service Service
	if id != "" {
		if db.First(&service, idInt).Error != nil {
			http.Error(w, "Resource not found", http.StatusNotFound)
			return
		}
	} else if name != "" {
		if db.Where("service_name = ?", name).First(&service).Error != nil {
			http.Error(w, "Resource not found", http.StatusNotFound)
			return
		}
	} else {
		http.Error(w, "ID or name parameter is required", http.StatusBadRequest)
		return
	}

	// Perform soft delete
	if err := db.Delete(&service).Error; err != nil {
		http.Error(w, "Failed to delete service", http.StatusInternalServerError)
		log.Printf("Error deleting service: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetUsers fetches user data from the database based on query parameters and responds with the results.
func GetUsers(w http.ResponseWriter, r *http.Request) {
	db := GetDBInstance()
	var users []User
	var user User

	// Get query parameters
	queryParams := r.URL.Query()
	username := queryParams.Get("username")
	id := queryParams.Get("id")

	// Fetch data based on query parameters
	if id != "" {
		// Get user by ID
		fetchAndRespond(w, func() error { return db.First(&user, "id = ?", id).Error }, &user)
		return
	}
	if username != "" {
		// Get user by username
		fetchAndRespond(w, func() error { return db.First(&user, "user_name = ?", username).Error }, &user)
	} else {
		// Get all users
		fetchAndRespond(w, func() error { return db.Find(&users).Error }, &users)
	}
}

// UpdateUser updates an existing user in the database based on the provided JSON payload.
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	db := GetDBInstance()

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Ensure ID is provided
	if user.ID == 0 {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	// Check if the user exists
	var existingUser User
	if db.First(&existingUser, user.ID).Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if err := db.Save(&user).Error; err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		log.Printf("Error updating user: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// CreateUser creates a new user in the database based on the provided JSON payload.
func CreateUser(w http.ResponseWriter, r *http.Request) {
	db := GetDBInstance()

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Check if the user already exists
	var existingUser User
	if db.Where("username = ?", user.Username).First(&existingUser).Error == nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	if err := db.Create(&user).Error; err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		log.Printf("Error creating user: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	db := GetDBInstance()

	// Parse query parameters
	id := r.URL.Query().Get("id")
	username := r.URL.Query().Get("username")

	// Convert ID to integer if provided
	var idInt int
	var err error
	if id != "" {
		idInt, err = strconv.Atoi(id)
		if handleDBQueryError(w, err, "Invalid ID parameter", http.StatusBadRequest) {
			return
		}
	}

	// Check if the user exists by ID or username
	var user User
	if id != "" {
		if db.First(&user, idInt).Error != nil {
			http.Error(w, "Resource not found", http.StatusNotFound)
			return
		}
	} else if username != "" {
		if db.Where("username = ?", username).First(&user).Error != nil {
			http.Error(w, "Resource not found", http.StatusNotFound)
			return
		}
	} else {
		http.Error(w, "ID or username parameter is required", http.StatusBadRequest)
		return
	}

	// Perform soft delete
	if err := db.Delete(&user).Error; err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		log.Printf("Error deleting user: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
