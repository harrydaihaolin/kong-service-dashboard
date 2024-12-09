package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type ServiceResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Add other fields as necessary
}

type UserResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Add other fields as necessary
}

var services []ServiceResponse // Assuming services is a slice of ServiceResponse
var users []UserResponse       // Assuming users is a slice of UserResponse

func GetAllServices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

func GetServiceById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range services {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	http.NotFound(w, r)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range users {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	http.NotFound(w, r)
}
