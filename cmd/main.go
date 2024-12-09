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

	"github.com/gorilla/mux"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func getDsn() string {
	dbHost := os.Getenv("SERVICE_DASHBOARD_DB_HOST")
	if dbHost == "" {
		if os.Getenv("UNIT_TEST") == "True" {
			dbHost = "localhost"
		} else {
			log.Println("Warning: SERVICE_DASHBOARD_DB_HOST is not set. Using default value 'host.docker.internal'.")
			dbHost = "host.docker.internal"
		}
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

	return dsn
}

func getGormDB() (*gorm.DB, error) {
	db, err := gorm.Open(gormPostgres.Open(getDsn()), &gorm.Config{
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
	GenerateDummyData(db)

	return db, nil
}

func main() {
	_, err := initDB()
	if err != nil {
		fmt.Println(err)
		panic("failed to initialize database")
	}

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	})
	router.HandleFunc("/v1/services", GetAllServices).Methods("GET")
	router.HandleFunc("/v1/services/{id}", GetServiceById).Methods("GET")
	router.HandleFunc("/v1/users", GetAllUsers).Methods("GET")
	router.HandleFunc("/v1/users/{id}", GetUserById).Methods("GET")

	// mount the router on the server
	http.Handle("/", router)

	fmt.Println("Starting server on :8080")
}
