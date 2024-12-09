package main

import (
	"context"
	"errors"
	"log"

	"gorm.io/gorm"
)

func GenerateDummyData(db *gorm.DB) {
	ctx := context.Background()

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
	}

	var existingService Service
	if err := db.WithContext(ctx).Where("id = ?", service.ID).First(&existingService).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Record does not exist, so create it
			if err := db.WithContext(ctx).Create(&service).Error; err != nil {
				log.Fatal(err)
			}
		} else {
			// An actual error occurred
			log.Fatal(err)
		}
	}

	var existingUser User
	if err := db.WithContext(ctx).Where("id = ?", user.ID).First(&existingUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Record does not exist, so create it
			if err := db.WithContext(ctx).Create(&user).Error; err != nil {
				log.Fatal(err)
			}
		} else {
			// An actual error occurred
			log.Fatal(err)
		}
	}

}
