package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

var JwtSecretKey = []byte("secret") //Note - This is a sample key. For production, use a secure key.
var AllowedRoles = []string{"admin", "user"}
var Permissions = map[string]map[string]bool{
	"admin": {
		"GET":    true,
		"POST":   true,
		"PUT":    true,
		"DELETE": true,
	},
	"user": {
		"GET": true,
	},
}
var whitelistedPaths = []string{"/v1/auth"}

// UserAuthentication is a handler function that authenticates a user based on the provided username and password.
func UserAuthentication(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	username := creds.Username
	password := creds.Password

	// Here you should validate the username and password with your user database
	// For simplicity, let's assume any username and password combination is valid
	if username == "" || password == "" {
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		return
	}
	db := GetDBInstance()
	var user User
	if err := db.Where("username = ? AND password = ?", username, password).First(&user).Error; err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	claims := CustomClaims{
		Role: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtSecretKey)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"token":"` + tokenString + `"}`))
}

// RoleBasedMiddleware is a middleware that checks if the request has a valid JWT token with the required role.
func RoleBasedMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests to whitelisted paths without token
		for _, path := range whitelistedPaths {
			if r.URL.Path == path {
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization token not provided", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(tokenString, "Bearer ") {
			http.Error(w, "Invalid authorization header format, requires Bearer prefix", http.StatusUnauthorized)
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return JwtSecretKey, nil
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("Token parsing failed: %v", err), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*CustomClaims)
		if !ok || !token.Valid {
			var errMsg string
			if !ok {
				errMsg = "Invalid token claims"
			} else if !token.Valid {
				errMsg = "Invalid token"
			}
			http.Error(w, errMsg, http.StatusUnauthorized)
			return
		}

		// Check if the role has the required permission
		method := r.Method
		if !checkPermission(claims.Role, method) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// CheckPermission is a helper function to check if a role has permission to perform an action.
func checkPermission(role, action string) bool {
	actions, exists := Permissions[role]
	if !exists {
		return false
	}
	if allowed, exists := actions[action]; exists && allowed {
		return true
	}
	return false
}
