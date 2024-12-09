package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func getGormDB() (*gorm.DB, error) {
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

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := gorm.Open(gormPostgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}

func initDB() (*gorm.DB, error) {
	db, err := getGormDB()
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Auto migrate the schema
	db.AutoMigrate(&Service{}, &ServiceVersion{}, &User{}, &UserProfile{})

	// Run migrations from the migrations folder
	migrationsPath := "file://./migrations" // Path to your migrations folder
	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres", driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange && err != migrate.ErrNilVersion {
		log.Fatal(err)
	}

	// Inject DummyData only if the table is empty
	if err := db.First(&Service{}).Error; err != nil {
		GenerateDummyData(db)
	} else {
		log.Println("Service table already has data, skipping dummy data injection")
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
