package main

import (
	"encoding/json"
	"fmt"
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
	handler := http.HandlerFunc(GetServices)
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
	handler := http.HandlerFunc(GetServices)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Service 1")
	assert.NotContains(t, rr.Body.String(), "Service 2")
}

func TestGetAllServicesWithPaginationSecondPage(t *testing.T) {
	req, err := http.NewRequest("GET", "/services?page=2&limit=1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetServices)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NotContains(t, rr.Body.String(), "Service 1")
	assert.Contains(t, rr.Body.String(), "Service 2")
}

func TestGetAllServicesWithSorting(t *testing.T) {
	req, err := http.NewRequest("GET", "/services?sortBy=service_name&order=desc", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetServices)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Service 2")
	assert.Contains(t, rr.Body.String(), "Service 1")
}

func TestGetAllServicesWithInvalidSorting(t *testing.T) {
	req, err := http.NewRequest("GET", "/services?sort_by=invalid&order=desc", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetServices)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetAllServicesWithInvalidOrder(t *testing.T) {
	req, err := http.NewRequest("GET", "/services?sort_by=service_name&order=invalid", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetServices)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetAllUsers(t *testing.T) {
	req, err := http.NewRequest("GET", "/users", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetUsers)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	expected_usernames := []string{"user1", "user2"}
	// the response body should contain the expected usernames
	for _, name := range expected_usernames {
		assert.Contains(t, rr.Body.String(), name)
	}
}

func TestGetServiceByIdNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/services?id=10000", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/services", GetServices)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetUserByIdNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/users?id=10000", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", GetUsers)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestSearchServicesByServiceName(t *testing.T) {
	req, err := http.NewRequest("GET", "/services?search_mode=true&name=1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetServices)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Service 1")    // Contains 1
	assert.NotContains(t, rr.Body.String(), "Service 2") // Does not contain 1
}

func TestSearchServicesByServiceNameNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/services?search_mode=true&name=3", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetServices)
	handler.ServeHTTP(rr, req)

	// return empty array
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "[]\n", rr.Body.String())
}

func TestSearchServicesByServiceNameWithLimit(t *testing.T) {
	req, err := http.NewRequest("GET", "/services?search_mode=true&name=Service&limit=1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetServices)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Service 1")    // Contains 1
	assert.NotContains(t, rr.Body.String(), "Service 2") // Does not contain 1
}

func TestGetServiceByServiceName(t *testing.T) {
	req, err := http.NewRequest("GET", "/services?name=Service%201", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetServices)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Service 1")
	assert.NotContains(t, rr.Body.String(), "Service 2")
}

func TestGetServiceById(t *testing.T) {
	// find latest service id
	req, err := http.NewRequest("GET", "/services", nil)
	assert.NoError(t, err)

	// extract the id from the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetServices)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Service 1")

	var services []map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &services)
	assert.NoError(t, err)
	serviceID := services[0]["ID"].(float64)

	req, err = http.NewRequest("GET", "/services?id="+fmt.Sprintf("%v", serviceID), nil)
	assert.NoError(t, err)

	// validate the response
	rr = httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/services", GetServices)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Service 1")
}
