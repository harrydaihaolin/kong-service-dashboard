package main

import (
	"gorm.io/gorm"
)

type Service struct {
	gorm.Model

	ServiceName        string           `gorm:"unique;not null" json:"service_name"`
	ServiceDescription string           `gorm:"type:text" json:"service_description"`
	Versions           []ServiceVersion `gorm:"foreignKey:ServiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"service_versions,omitempty"`
}

type ServiceVersion struct {
	gorm.Model

	ServiceID                 uint   `gorm:"not null" json:"service_id"`
	ServiceVersionName        string `gorm:"unique" json:"service_version_name"`
	ServiceVersionURL         string `gorm:"type:text" json:"service_version_url"`
	ServiceVersionDescription string `gorm:"type:text" json:"service_version_description"`
}

type User struct {
	gorm.Model

	Username    string      `gorm:"unique;not null" json:"username"`
	UserProfile UserProfile `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user_profile,omitempty"`
}

type UserProfile struct {
	gorm.Model

	UserID    uint   `gorm:"not null;unique" json:"user_id"`
	FirstName string `gorm:"type:varchar(255)" json:"first_name"`
	LastName  string `gorm:"type:varchar(255)" json:"last_name"`
	Email     string `gorm:"unique;not null" json:"email"`
}
