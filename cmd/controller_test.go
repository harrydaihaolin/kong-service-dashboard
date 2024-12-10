package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHandlers(t *testing.T) {
	var tests = []struct {
		name       string
		method     string
		url        string
		handler    http.HandlerFunc
		statusCode int
		body       string
	}{
		{"TestGetAllServices", "GET", "/services", GetServices, http.StatusOK, "Service 1"},
		{"TestGetAllServicesWithServiceVersion", "GET", "/services?load_version=true", GetServices, http.StatusOK, "Service 1 Version 1"},
		{"TestGetAllServicesWithPagination", "GET", "/services?page=1&limit=1", GetServices, http.StatusOK, "Service 1"},
		{"TestGetAllServicesWithPaginationSecondPage", "GET", "/services?page=2&limit=1", GetServices, http.StatusOK, "Service 2"},
		{"TestGetAllServicesWithSorting", "GET", "/services?sortBy=service_name&order=desc", GetServices, http.StatusOK, "Service 2"},
		{"TestGetAllServicesWithInvalidSorting", "GET", "/services?sort_by=invalid&order=desc", GetServices, http.StatusBadRequest, ""},
		{"TestGetAllServicesWithInvalidOrder", "GET", "/services?sort_by=service_name&order=invalid", GetServices, http.StatusBadRequest, ""},
		{"TestGetAllUsers", "GET", "/users", GetUsers, http.StatusOK, "user1"},
		{"TestGetServiceByIdNotFound", "GET", "/services?id=10000", GetServices, http.StatusNotFound, ""},
		{"TestGetUserByIdNotFound", "GET", "/users?id=10000", GetUsers, http.StatusNotFound, ""},
		{"TestSearchServicesByServiceName", "GET", "/services?search_mode=true&name=1", GetServices, http.StatusOK, "Service 1"},
		{"TestSearchServicesByServiceNameNotFound", "GET", "/services?search_mode=true&name=3", GetServices, http.StatusOK, "[]\n"},
		{"TestSearchServicesByServiceNameWithLimit", "GET", "/services?search_mode=true&name=Service&limit=1", GetServices, http.StatusOK, "Service 1"},
		{"TestGetServiceByServiceName", "GET", "/services?name=Service%201", GetServices, http.StatusOK, "Service 1"},
		{"TestGetServiceById", "GET", "/services", GetServices, http.StatusOK, "Service 1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.url, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := tt.handler
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.statusCode, rr.Code)
			if tt.body != "" {
				assert.Contains(t, rr.Body.String(), tt.body)
			}
		})
	}
}

func TestServiceMutateHandlers(t *testing.T) {
	// Seed database for delete operation
	db := GetDBInstance()
	db.Create(&Service{ServiceName: "ServiceForDeletion", ServiceDescription: "Service for deletion description"})
	db.Create(&Service{ServiceName: "ServiceForDeletion2", ServiceDescription: "Service for deletion description"})
	db.Create(&Service{ServiceName: "ServiceForUpdate", ServiceDescription: "Service for update description"})
	// fetch the service id for delete operation and update operation
	var serviceForDeletion Service
	db.Where("service_name = ?", "ServiceForDeletion").First(&serviceForDeletion)
	var serviceForUpdate Service
	db.Where("service_name = ?", "ServiceForUpdate").First(&serviceForUpdate)

	// Seed database for service version delete operation
	db.Create(&ServiceVersion{ServiceID: serviceForUpdate.ID, ServiceVersionName: "Service 1 Version 1", ServiceVersionURL: "http://service1.com", ServiceVersionDescription: "Service 1 Version 1 Description"})

	// fetch the service version id for delete operation and update operation
	var serviceVersionForDeletion ServiceVersion
	db.Where("service_version_name = ?", "Service 1 Version 1").First(&serviceVersionForDeletion)
	var serviceVersionForUpdate ServiceVersion
	db.Where("service_version_name = ?", "Service 1 Version 1").First(&serviceVersionForUpdate)

	tests := []struct {
		name       string
		method     string
		url        string
		handler    http.HandlerFunc
		statusCode int
		body       string
		payload    string
	}{
		{"TestCreateService", "POST", "/services", CreateService, http.StatusCreated, "Service 5", `{"service_name": "Service 5", "service_description": "Service 5 Description"}`},
		{"TestCreateServiceAlreadyCreated", "POST", "/services", CreateService, http.StatusConflict, "Service already exists", `{"service_name": "Service 1", "service_description": "Service 1 Description"}`},
		{"TestCreateServiceWithInvalidJsonPayload", "POST", "/services", CreateService, http.StatusBadRequest, "Invalid JSON payload", "invalid json"},
		{"TestUpdateService", "PUT", "/services", UpdateService, http.StatusOK, "Service 1 Updated", `{"id": ` + fmt.Sprint(serviceForUpdate.ID) + `, "service_name": "Service 1 Updated"}`},
		{"TestUpdateServiceNotFound", "PUT", "/services", UpdateService, http.StatusNotFound, "Service not found", `{"id": 10000, "service_name": "Non-existent Service"}`},
		{"TestUpdateServiceInvalidJsonPayload", "PUT", "/services", UpdateService, http.StatusBadRequest, "Invalid JSON payload", "invalid json"},
		{"TestUpdateServiceInternalServerError", "PUT", "/services", UpdateService, http.StatusInternalServerError, "Failed to update service", `{"id": ` + fmt.Sprint(serviceForUpdate.ID) + `, "service_name": "Service 1"}`},
		{"TestDeleteService", "DELETE", "/services?id=" + fmt.Sprint(serviceForDeletion.ID), DeleteService, http.StatusOK, "", ""},
		{"TestDeleteServiceByName", "DELETE", "/services?name=ServiceForDeletion2", DeleteService, http.StatusOK, "", ""},
		{"TestDeleteServiceNotFound", "DELETE", "/services?id=10000", DeleteService, http.StatusNotFound, "Resource not found", ""},
		{"TestCreateServiceVersion", "POST", "/service_versions", CreateServiceVersion, http.StatusCreated, "Service 1 Version 2", `{"service_id": ` + fmt.Sprint(serviceForUpdate.ID) + `, "service_version_name": "Service 1 Version 2", "service_version_url": "http://service1.com", "service_version_description": "Service 1 Version 2 Description"}`},
		{"TestCreateServiceVersionInvalidJsonPayload", "POST", "/service_versions", CreateServiceVersion, http.StatusBadRequest, "Invalid JSON payload", "invalid json"},
		{"TestCreateServiceVersionServiceNotFound", "POST", "/service_versions", CreateServiceVersion, http.StatusNotFound, "Service not found", `{"service_id": 10000, "service_version_name": "Service 1 Version 2", "service_version_url": "http://service1.com", "service_version_description": "Service 1 Version 2 Description"}`},
		{"TestUpdateServiceVersion", "PUT", "/service_versions", UpdateServiceVersion, http.StatusOK, "Service 1 Version 1", `{"id": ` + fmt.Sprint(serviceVersionForUpdate.ID) + `, "service_id": ` + fmt.Sprint(serviceForUpdate.ID) + `, "service_version_name": "Service 1 Version 1"}`},
		{"TestUpdateServiceVersionNotFound", "PUT", "/service_versions", UpdateServiceVersion, http.StatusNotFound, "Version not found", `{"id": 10000, "service_version_name": "Non-existent Service Version"}`},
		{"TestUpdateServiceVersionInvalidJsonPayload", "PUT", "/service_versions", UpdateServiceVersion, http.StatusBadRequest, "Invalid JSON payload", "invalid json"},
		{"TestDeleteServiceVersion", "DELETE", "/service_versions?id=" + fmt.Sprint(serviceVersionForDeletion.ID), DeleteServiceVersion, http.StatusOK, "", ""},
		{"TestDeleteServiceVersionNotFound", "DELETE", "/service_versions?id=10000", DeleteServiceVersion, http.StatusNotFound, "Resource not found", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			var err error

			// Create request based on method
			if tt.method == "POST" || tt.method == "PUT" {
				req, err = http.NewRequest(tt.method, tt.url, strings.NewReader(tt.payload))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req, err = http.NewRequest(tt.method, tt.url, nil)
			}

			assert.NoError(t, err)

			// Mock response recorder
			rr := httptest.NewRecorder()
			handler := tt.handler
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.statusCode, rr.Code)
			if tt.body != "" {
				assert.Contains(t, rr.Body.String(), tt.body)
			}
		})
	}
}

func TestUserMutateHandlers(t *testing.T) {
	// Seed database for delete operation
	db := GetDBInstance()
	db.Create(&User{Username: "UserForDeletion", Password: "password", Role: "user"})
	db.Create(&User{Username: "UserForDeletion2", Password: "password", Role: "user"})
	db.Create(&User{Username: "UserForUpdate", Password: "password", Role: "user"})
	// fetch the user id for delete operation and update operation
	var userForDeletion User
	db.Where("username = ?", "UserForDeletion").First(&userForDeletion)
	var userForUpdate User
	db.Where("username = ?", "UserForUpdate").First(&userForUpdate)

	tests := []struct {
		name       string
		method     string
		url        string
		handler    http.HandlerFunc
		statusCode int
		body       string
		payload    string
	}{
		{"TestCreateUser", "POST", "/users", CreateUser, http.StatusCreated, "User 5", `{"username": "User 5", "password": "password", "role": "user"}`},
		{"TestCreateUserAlreadyCreated", "POST", "/users", CreateUser, http.StatusConflict, "User already exists", `{"username": "user1", "password": "password", "role": "user"}`},
		{"TestCreateUserWithInvalidJsonPayload", "POST", "/users", CreateUser, http.StatusBadRequest, "Invalid JSON payload", "invalid json"},
		{"TestUpdateUser", "PUT", "/users", UpdateUser, http.StatusOK, "User 1 Updated", `{"id": ` + fmt.Sprint(userForUpdate.ID) + `, "username": "User 1 Updated"}`},
		{"TestUpdateUserNotFound", "PUT", "/users", UpdateUser, http.StatusNotFound, "User not found", `{"id": 10000, "username": "Non-existent User"}`},
		{"TestUpdateUserInvalidJsonPayload", "PUT", "/users", UpdateUser, http.StatusBadRequest, "Invalid JSON payload", "invalid json"},
		{"TestDeleteUser", "DELETE", "/users?id=" + fmt.Sprint(userForDeletion.ID), DeleteUser, http.StatusOK, "", ""},
		{"TestDeleteUserByName", "DELETE", "/users?username=UserForDeletion2", DeleteUser, http.StatusOK, "", ""},
		{"TestDeleteUserNotFound", "DELETE", "/users?id=10000", DeleteUser, http.StatusNotFound, "Resource not found", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			var err error

			// Create request based on method
			if tt.method == "POST" || tt.method == "PUT" {
				req, err = http.NewRequest(tt.method, tt.url, strings.NewReader(tt.payload))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req, err = http.NewRequest(tt.method, tt.url, nil)
			}

			assert.NoError(t, err)

			// Mock response recorder
			rr := httptest.NewRecorder()
			handler := tt.handler
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.statusCode, rr.Code)
			if tt.body != "" {
				assert.Contains(t, rr.Body.String(), tt.body)
			}
		})
	}
}
