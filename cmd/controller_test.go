package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetAllServices(t *testing.T) {

	services = []ServiceResponse{
		{ID: "1", Name: "Service1"},
		{ID: "2", Name: "Service2"},
	}

	req, err := http.NewRequest("GET", "/services", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllServices)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	expected := `[{"id":"1","name":"Service1"},{"id":"2","name":"Service2"}]`
	assert.JSONEq(t, expected, rr.Body.String())
}

func TestGetServiceById(t *testing.T) {
	services = []ServiceResponse{
		{ID: "1", Name: "Service1"},
		{ID: "2", Name: "Service2"},
	}

	req, err := http.NewRequest("GET", "/services/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/services/{id}", GetServiceById)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	expected := `{"id":"1","name":"Service1"}`
	assert.JSONEq(t, expected, rr.Body.String())
}

func TestGetAllUsers(t *testing.T) {
	users = []UserResponse{
		{ID: "1", Name: "User1"},
		{ID: "2", Name: "User2"},
	}

	req, err := http.NewRequest("GET", "/users", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllUsers)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	expected := `[{"id":"1","name":"User1"},{"id":"2","name":"User2"}]`
	assert.JSONEq(t, expected, rr.Body.String())
}

func TestGetUserById(t *testing.T) {
	users = []UserResponse{
		{ID: "1", Name: "User1"},
		{ID: "2", Name: "User2"},
	}

	req, err := http.NewRequest("GET", "/users/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", GetUserById)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	expected := `{"id":"1","name":"User1"}`
	assert.JSONEq(t, expected, rr.Body.String())
}

func TestGetServiceByIdNotFound(t *testing.T) {
	services = []ServiceResponse{
		{ID: "1", Name: "Service1"},
		{ID: "2", Name: "Service2"},
	}

	req, err := http.NewRequest("GET", "/services/3", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/services/{id}", GetServiceById)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetUserByIdNotFound(t *testing.T) {
	users = []UserResponse{
		{ID: "1", Name: "User1"},
		{ID: "2", Name: "User2"},
	}

	req, err := http.NewRequest("GET", "/users/3", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", GetUserById)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

// API tests
func TestGetAllServicesAPI(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(GetAllServices))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestGetServiceByIdAPI(t *testing.T) {
	services = []ServiceResponse{
		{ID: "1", Name: "Service1"},
		{ID: "2", Name: "Service2"},
	}

	router := mux.NewRouter()
	router.HandleFunc("/services/{id}", GetServiceById)
	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/services/1")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestGetAllUsersAPI(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(GetAllUsers))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestGetUserByIdAPI(t *testing.T) {
	users = []UserResponse{
		{ID: "1", Name: "User1"},
		{ID: "2", Name: "User2"},
	}

	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", GetUserById)
	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/users/1")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestGetServiceByIdNotFoundAPI(t *testing.T) {
	services = []ServiceResponse{
		{ID: "1", Name: "Service1"},
		{ID: "2", Name: "Service2"},
	}

	router := mux.NewRouter()
	router.HandleFunc("/services/{id}", GetServiceById)
	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/services/3")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestGetUserByIdNotFoundAPI(t *testing.T) {
	users = []UserResponse{
		{ID: "1", Name: "User1"},
		{ID: "2", Name: "User2"},
	}

	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", GetUserById)
	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/users/3")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}
