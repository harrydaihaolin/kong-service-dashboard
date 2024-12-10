package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	postgresGorm "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Required for file-based migrations
)

var (
	dbInstance *gorm.DB
	once       sync.Once
)

func GetDBInstance() *gorm.DB {
	once.Do(func() {
		var err error
		dsn := getDsn()
		dbInstance, err = gorm.Open(postgresGorm.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("failed to connect database")
		}
	})
	return dbInstance
}

func getDsn() string {
	dbHost := getEnv("SERVICE_DASHBOARD_DB_HOST", "host.docker.internal")
	dbPort := getEnv("SERVICE_DASHBOARD_DB_PORT", "5432")
	dbUser := getEnv("SERVICE_DASHBOARD_DB_USER", "postgres")
	dbPassword := getEnv("SERVICE_DASHBOARD_DB_PASSWORD", "example")
	dbName := getEnv("SERVICE_DASHBOARD_DB_NAME", "postgres")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		if key == "SERVICE_DASHBOARD_DB_HOST" && os.Getenv("UNIT_TEST") == "True" {
			return "localhost"
		}
		log.Printf("Warning: %s is not set. Using default value '%s'.", key, fallback)
		return fallback
	}
	return value
}

func InitDB() {
	db := GetDBInstance()

	db.AutoMigrate(&Service{}, &ServiceVersion{}, &User{}, &UserProfile{})

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})

	m, err := migrate.NewWithDatabaseInstance("file://./migrations", "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange && err != migrate.ErrNilVersion {
		log.Fatal(err)
	}

	GenerateDummyData(db)
}

func GenerateDummyData(db *gorm.DB) {
	ctx := context.Background()

	var serviceCount, userCount int64
	if err := db.WithContext(ctx).Model(&Service{}).Count(&serviceCount).Error; err != nil {
		log.Fatal(err)
	}
	if err := db.WithContext(ctx).Model(&User{}).Count(&userCount).Error; err != nil {
		log.Fatal(err)
	}

	if serviceCount == 0 && userCount == 0 {
		services := []Service{
			{
				ServiceName:        "Service 1",
				ServiceDescription: "Service 1 Description",
				Versions: []ServiceVersion{
					{
						ServiceVersionName:        "Service 1 Version 1",
						ServiceVersionURL:         "http://service1.com",
						ServiceVersionDescription: "Service 1 Version 1 Description",
					},
				},
			},
			{
				ServiceName:        "Service 2",
				ServiceDescription: "Service 2 Description",
				Versions: []ServiceVersion{
					{
						ServiceVersionName:        "Service 2 Version 1",
						ServiceVersionURL:         "http://service2.com",
						ServiceVersionDescription: "Service 2 Version 1 Description",
					},
				},
			},
		}

		users := []User{
			{
				Username: "user1",
				Password: "password",
				Role:     "user",
				UserProfile: UserProfile{
					FirstName: "User",
					LastName:  "One",
					Email:     "abc@gmail.com",
				},
			},
			{
				Username: "user2",
				Password: "password",
				Role:     "admin",
				UserProfile: UserProfile{
					FirstName: "User",
					LastName:  "Two",
					Email:     "def@gmail.com",
				},
			},
		}

		if err := db.WithContext(ctx).Create(&services).Error; err != nil {
			log.Fatal(err)
		}
		log.Printf("Inserted services %+v", services)
		if err := db.WithContext(ctx).Create(&users).Error; err != nil {
			log.Fatal(err)
		}
		log.Printf("Inserted users: %+v", users)
	}

	return
}
