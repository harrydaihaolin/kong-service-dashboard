package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/gorilla/mux"
)

func main() {
	InitDB()

	router := mux.NewRouter()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Println("ok")
	})
	router.HandleFunc("/v1/services", GetServices).Methods("GET")
	router.HandleFunc("/v1/services", CreateService).Methods("POST")
	router.HandleFunc("/v1/services", UpdateService).Methods("PUT")
	router.HandleFunc("/v1/services", DeleteService).Methods("DELETE")
	router.HandleFunc("/v1/service_versions", CreateServiceVersion).Methods("POST")
	router.HandleFunc("/v1/service_versions", UpdateServiceVersion).Methods("PUT")
	router.HandleFunc("/v1/service_versions", DeleteServiceVersion).Methods("DELETE")
	router.HandleFunc("/v1/users", GetUsers).Methods("GET")
	router.HandleFunc("/v1/users", CreateUser).Methods("POST")
	router.HandleFunc("/v1/users", UpdateUser).Methods("PUT")
	router.HandleFunc("/v1/users", DeleteUser).Methods("DELETE")
	router.HandleFunc("/v1/auth", UserAuthentication).Methods("POST")

	// Add logger middleware to the router
	loggedMux := LoggerMiddleware(router)
	// Add Role Based middleware to the router
	roleBasedMux := RoleBasedMiddleware(loggedMux)

	if err := http.ListenAndServe(":8080", roleBasedMux); err != nil {
		fmt.Println("Error starting server:", err)
	}
	log.Println("Starting server on :8080")
}
