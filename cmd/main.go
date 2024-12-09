package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func initDB() (*gorm.DB, error) {
	// Get all DB connection variables from the environment
	migrationsPath := "./migrations" // Path to your migrations folder
	dbHost := os.Getenv("SERVICE_DASHBOARD_DB_HOST")
	if dbHost == "" {
		dbHost = "host.docker.internal"
	}
	dbPort := os.Getenv("SERVICE_DASHBOARD_DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbUser := os.Getenv("SERVICE_DASHBOARD_DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}
	dbPassword := os.Getenv("SERVICE_DASHBOARD_DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "example"
	}
	dbName := os.Getenv("SERVICE_DASHBOARD_DB_NAME")
	if dbName == "" {
		dbName = "postgres"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		return nil, fmt.Errorf("error creating extension: %w", err)
	}

	if err := db.AutoMigrate(&User{}, &Service{}, &ServiceVersion{}, &UserRole{}, &UserProfile{}); err != nil {
		return nil, fmt.Errorf("error during migration: %w", err)
	}

	// Inject DummyData
	dummyData := &DummyData{}
	dummyData.Generate(db)

	// Run migrations
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
	m, err := migrate.New("file://"+migrationsPath, dbURL)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v\n", err)
	}

	// Run all available migrations up
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migrations: %v\n", err)
	} else {
		log.Println("Migrations ran successfully")
	}

	return db, nil
}

func main() {
	_, err := initDB()
	if err != nil {
		fmt.Println(err)
		panic("failed to initialize database")
	}

	http.HandleFunc("/", handler)
	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
