package main

import (
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

// Generate dummy data for testing
type DummyData struct{}

func (d *DummyData) Generate(db *gorm.DB) {
	service := Service{
		ServiceName:        "Service 1",
		ServiceDescription: "Service 1 Description",
		Versions: []ServiceVersion{
			{
				ServiceVersionName:        "Service 1 Version 1",
				ServiceVersionURL:         "http://service1.com",
				ServiceVersionDescription: "Service 1 Version 1 Description",
			},
		},
	}

	user := User{
		Username: "user1",
		UserProfile: UserProfile{
			FirstName: "User",
			LastName:  "One",
			Email:     "abc@gmail.com",
		},
		UserRole: UserRole{
			Role: "admin",
		},
	}

	db.Create(&service)
	db.Create(&user)
}
