package main

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupDB(t *testing.T) *gorm.DB {
	os.Setenv("UNIT_TEST", "True")
	dsn := getDsn()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)
	return db
}

func TestService_Create(t *testing.T) {
	db := setupDB(t)
	db.AutoMigrate(&Service{}, &ServiceVersion{})

	service := Service{
		ServiceName:        "Test Service " + uuid.New().String(),
		ServiceDescription: "This is a test service",
	}

	result := db.Create(&service)
	assert.NoError(t, result.Error)
	assert.NotZero(t, service.ID)
}

func TestServiceVersion_Create(t *testing.T) {
	db := setupDB(t)
	db.AutoMigrate(&Service{}, &ServiceVersion{})

	service := Service{
		ServiceName:        "Test Service " + uuid.New().String(),
		ServiceDescription: "This is a test service",
	}
	db.Create(&service)

	serviceVersion := ServiceVersion{
		ServiceID:                 service.ID,
		ServiceVersionName:        "v1.0",
		ServiceVersionURL:         "http://example.com/v1.0",
		ServiceVersionDescription: "Initial version",
	}

	result := db.Create(&serviceVersion)
	assert.NoError(t, result.Error)
	assert.NotZero(t, serviceVersion.ID)
}

func TestUser_Create(t *testing.T) {
	db := setupDB(t)
	db.AutoMigrate(&User{}, &UserProfile{})

	user := User{
		Username: "testuser_" + uuid.New().String(),
		UserProfile: UserProfile{
			FirstName: "Test",
			LastName:  "User",
			Email:     "testuser@example.com",
		},
	}

	result := db.Create(&user)
	assert.NoError(t, result.Error)
	assert.NotZero(t, user.ID)
	assert.NotZero(t, user.UserProfile.ID)
}

func TestUserProfile_Create(t *testing.T) {
	db := setupDB(t)
	db.AutoMigrate(&User{}, &UserProfile{})

	user := User{
		Username: "testuser_" + uuid.New().String(),
	}
	db.Create(&user)

	userProfile := UserProfile{
		UserID:    user.ID,
		FirstName: "Test",
		LastName:  "User",
		Email:     "testuser@example.com",
	}

	result := db.Create(&userProfile)
	assert.NoError(t, result.Error)
	assert.NotZero(t, userProfile.ID)
}

func TestService_Update(t *testing.T) {
	db := setupDB(t)
	db.AutoMigrate(&Service{}, &ServiceVersion{})

	service := Service{
		ServiceName:        "Test Service " + uuid.New().String(),
		ServiceDescription: "This is a test service",
	}
	db.Create(&service)

	service.ServiceDescription = "Updated description"
	result := db.Save(&service)
	assert.NoError(t, result.Error)
	assert.Equal(t, "Updated description", service.ServiceDescription)
}

func TestUser_Update(t *testing.T) {
	db := setupDB(t)
	db.AutoMigrate(&User{}, &UserProfile{})

	user := User{
		Username: "testuser_" + uuid.New().String(),
		UserProfile: UserProfile{
			FirstName: "Test",
			LastName:  "User",
			Email:     "testuser@example.com",
		},
	}
	db.Create(&user)

	user.UserProfile.FirstName = "Updated"
	result := db.Save(&user.UserProfile)
	assert.NoError(t, result.Error)
	assert.Equal(t, "Updated", user.UserProfile.FirstName)
}
