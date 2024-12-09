package main

import (
	"log"
	"os"
	"testing"

	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	db := GetDBInstance()
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
