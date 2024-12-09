package main

import (
	"log"
	"os"
	"testing"

	"gorm.io/gorm"
)

// Test Data
var USER_TEST_DATA = `[{"CreatedAt":"2024-12-09T00:16:52.874698-08:00", "DeletedAt":null, "ID":37, "UpdatedAt":"2024-12-09T00:16:52.874698-08:00", "user_profile":{"CreatedAt":"0001-01-01T00:00:00Z", "DeletedAt":null, "ID":0, "UpdatedAt":"0001-01-01T00:00:00Z", "email":"", "first_name":"", "last_name":"", "user_id":0}, "username":"user1"}, {"CreatedAt":"2024-12-09T00:16:52.874698-08:00", "DeletedAt":null, "ID":38, "UpdatedAt":"2024-12-09T00:16:52.874698-08:00", "user_profile":{"CreatedAt":"0001-01-01T00:00:00Z", "DeletedAt":null, "ID":0, "UpdatedAt":"0001-01-01T00:00:00Z", "email":"", "first_name":"", "last_name":"", "user_id":0}, "username":"user2"}]`
var SERVICE_TEST_DATA = `[{"CreatedAt":"2024-12-09T00:16:52.868929-08:00", "DeletedAt":null, "ID":37, "UpdatedAt":"2024-12-09T00:16:52.868929-08:00", "service_description":"Service 1 Description", "service_name":"Service 1"}, {"CreatedAt":"2024-12-09T00:16:52.868929-08:00", "DeletedAt":null, "ID":38, "UpdatedAt":"2024-12-09T00:16:52.868929-08:00", "service_description":"Service 2 Description", "service_name":"Service 2"}]`
var SINGLE_SERVICE_TEST_DATA = `[{"CreatedAt":"2024-12-09T00:16:52.874698-08:00", "DeletedAt":null, "ID":37, "UpdatedAt":"2024-12-09T00:16:52.874698-08:00", "user_profile":{"CreatedAt":"0001-01-01T00:00:00Z", "DeletedAt":null, "ID":0, "UpdatedAt":"0001-01-01T00:00:00Z", "email":"", "first_name":"", "last_name":"", "user_id":0}, "username":"user1"}, {"CreatedAt":"2024-12-09T00:16:52.874698-08:00", "DeletedAt":null, "ID":38, "UpdatedAt":"2024-12-09T00:16:52.874698-08:00", "user_profile":{"CreatedAt":"0001-01-01T00:00:00Z", "DeletedAt":null, "ID":0, "UpdatedAt":"0001-01-01T00:00:00Z", "email":"", "first_name":"", "last_name":"", "user_id":0}, "username":"user2"}]`

func TestMain(m *testing.M) {
	db := setupDB(nil)
	if db == nil {
		panic("failed to connect database")
	}

	// Load test data
	log.Println("Generating dummy data")
	GenerateDummyData(db)

	// Run tests
	log.Println("Running tests")
	code := m.Run()

	// Cleanup
	log.Println("Cleaning up database")
	cleanup(db)

	os.Exit(code)
}

func cleanup(db *gorm.DB) {
	db.Exec("DELETE FROM user_profiles")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM service_versions")
	db.Exec("DELETE FROM services")
}
