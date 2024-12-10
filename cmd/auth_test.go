package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserAuthentication(t *testing.T) {
	tests := []struct {
		name           string
		username       string
		password       string
		expectedStatus int
	}{
		{"AdminUser", "user1", "password", http.StatusOK},
		{"RegularUser", "user2", "password", http.StatusOK},
		{"WrongUser", "user3", "password", http.StatusUnauthorized},
		{"WrongPassword", "user1", "wrongpassword", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonStr := `{"username":"` + tt.username + `","password":"` + tt.password + `"}`
			req, err := http.NewRequest("POST", "/v1/auth", strings.NewReader(jsonStr))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(UserAuthentication)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestRoleBasedMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		role           string
		method         string
		expectedStatus int
	}{
		{"AdminGET", "admin", "GET", http.StatusOK},
		{"AdminPOST", "admin", "POST", http.StatusOK},
		{"AdminPUT", "admin", "PUT", http.StatusOK},
		{"AdminDELETE", "admin", "DELETE", http.StatusOK},
		{"UserGET", "user", "GET", http.StatusOK},
		{"UserPOST", "user", "POST", http.StatusForbidden},
		{"UserPUT", "user", "PUT", http.StatusForbidden},
		{"UserDELETE", "user", "DELETE", http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use /v1/auth to get the JWT token
			var jsonStr string
			if tt.role == "admin" {
				jsonStr = `{"username":"user2","password":"password"}`
			} else {
				jsonStr = `{"username":"user1","password":"password"}`
			}
			authReq, err := http.NewRequest("POST", "/v1/auth", strings.NewReader(jsonStr))
			assert.NoError(t, err)
			authReq.Header.Set("Content-Type", "application/json")

			// Get the token from the response
			authRr := httptest.NewRecorder()
			authHandler := http.HandlerFunc(UserAuthentication)
			authHandler.ServeHTTP(authRr, authReq)

			assert.Equal(t, http.StatusOK, authRr.Code)
			var response map[string]string
			err = json.Unmarshal(authRr.Body.Bytes(), &response)
			assert.NoError(t, err)
			token, ok := response["token"]
			assert.True(t, ok)

			// Use the token to access the protected resource
			req, err := http.NewRequest(tt.method, "/v1/services", nil)
			assert.NoError(t, err)

			req.Header.Set("Authorization", "Bearer "+token)

			rr := httptest.NewRecorder()
			handler := RoleBasedMiddleware(http.HandlerFunc(GetServices))
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestAuthFlow(t *testing.T) {
	tests := []struct {
		name           string
		username       string
		password       string
		expectedStatus int
	}{
		{"AdminUser", "user1", "password", http.StatusOK},
		{"RegularUser", "user2", "password", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Login
			jsonStr := `{"username":"` + tt.username + `","password":"` + tt.password + `"}`
			req, err := http.NewRequest("POST", "/v1/auth", strings.NewReader(jsonStr))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(UserAuthentication)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)

			// Get protected resource
			req, err = http.NewRequest("GET", "/v1/services", nil)
			assert.NoError(t, err)

			token := rr.Header().Get("Authorization")
			req.Header.Set("Authorization", token)

			rr = httptest.NewRecorder()
			handler = http.HandlerFunc(GetServices)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)
			expected_service_names := []string{"Service 1", "Service 2"}
			// the response body should contain the expected service names
			for _, name := range expected_service_names {
				assert.Contains(t, rr.Body.String(), name)
			}
		})
	}
}
