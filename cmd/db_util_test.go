package main

import (
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestGenerateDummyData(t *testing.T) {
	// Create a PostgreSQL database connection for testing
	dsn := getDsn()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// Auto migrate the schema
	if err := db.AutoMigrate(&Service{}, &ServiceVersion{}, &User{}, &UserProfile{}); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	// Call the function to generate dummy data
	GenerateDummyData(db)

	// Check if the service was inserted
	var service Service
	if err := db.Where("service_name = ?", "Service 1").First(&service).Error; err != nil {
		t.Fatalf("failed to find service: %v", err)
	}

	// Check if the service version was inserted
	var serviceVersion ServiceVersion
	if err := db.Where("service_version_name = ?", "Service 1 Version 1").First(&serviceVersion).Error; err != nil {
		t.Fatalf("failed to find service version: %v", err)
	}

	// Check if the user was inserted
	var user User
	if err := db.Where("username = ?", "user1").First(&user).Error; err != nil {
		t.Fatalf("failed to find user: %v", err)
	}

	// Check if the user profile was inserted
	var userProfile UserProfile
	if err := db.Where("email = ?", "abc@gmail.com").First(&userProfile).Error; err != nil {
		t.Fatalf("failed to find user profile: %v", err)
	}
}
