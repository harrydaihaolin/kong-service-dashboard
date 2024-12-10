#!/bin/bash

# Variables
USERNAME="user1"
PASSWORD="password"
AUTH_URL="http://localhost:8080/v1/auth"
SERVICES_URL="http://localhost:8080/v1/services"

# Fetch JWT token
TOKEN=$(curl -s -X POST $AUTH_URL \
    -H "Content-Type: application/json" \
    -d '{"username": "'$USERNAME'", "password": "'$PASSWORD'"}' | jq -r '.token')

# Check if token is not empty
if [ -z "$TOKEN" ]; then
    echo "Failed to fetch JWT token"
    exit 1
fi

# Use JWT token to get all services
curl -s -X GET $SERVICES_URL \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json"