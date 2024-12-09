package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// Load Test Data

func TestGetAllServices(t *testing.T) {
	req, err := http.NewRequest("GET", "/services", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllServices)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	expected_service_names := []string{"Service 1", "Service 2"}
	// the response body should contain the expected service names
	for _, name := range expected_service_names {
		assert.Contains(t, rr.Body.String(), name)
	}
}

func TestGetAllServicesWithPagination(t *testing.T) {
	req, err := http.NewRequest("GET", "/services?page=1&limit=1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllServices)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Service 1")
	assert.NotContains(t, rr.Body.String(), "Service 2")
}

func TestGetAllServicesWithPaginationSecondPage(t *testing.T) {
	req, err := http.NewRequest("GET", "/services?page=2&limit=1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllServices)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NotContains(t, rr.Body.String(), "Service 1")
	assert.Contains(t, rr.Body.String(), "Service 2")
}

func TestGetAllServicesWithSorting(t *testing.T) {
	req, err := http.NewRequest("GET", "/services?sortBy=service_name&order=desc", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllServices)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Service 2")
	assert.Contains(t, rr.Body.String(), "Service 1")
}

func TestGetAllServicesWithInvalidSorting(t *testing.T) {
	req, err := http.NewRequest("GET", "/services?sort_by=invalid&order=desc", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllServices)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetAllServicesWithInvalidOrder(t *testing.T) {
	req, err := http.NewRequest("GET", "/services?sort_by=service_name&order=invalid", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllServices)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetAllUsers(t *testing.T) {
	req, err := http.NewRequest("GET", "/users", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllUsers)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	expected_usernames := []string{"user1", "user2"}
	// the response body should contain the expected usernames
	for _, name := range expected_usernames {
		assert.Contains(t, rr.Body.String(), name)
	}
}

func TestGetServiceByIdNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/services/10000", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/services/{id}", GetServiceById)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetUserByIdNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/users/10000", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", GetUserById)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}
