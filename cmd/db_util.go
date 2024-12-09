package main

import (
	"context"
	"log"

	"gorm.io/gorm"
)

func GenerateDummyData(db *gorm.DB) {
	ctx := context.Background()

	var serviceCount int64
	var userCount int64

	// Check if the Service table is empty
	if err := db.WithContext(ctx).Model(&Service{}).Count(&serviceCount).Error; err != nil {
		log.Fatal(err)
	}

	// Check if the User table is empty
	if err := db.WithContext(ctx).Model(&User{}).Count(&userCount).Error; err != nil {
		log.Fatal(err)
	}

	// If both tables are empty, inject dummy data
	if serviceCount == 0 && userCount == 0 {
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

		if err := db.WithContext(ctx).Create(&service).Error; err != nil {
			log.Fatal(err)
		}

		if err := db.WithContext(ctx).Create(&user).Error; err != nil {
			log.Fatal(err)
		}
	}
}
