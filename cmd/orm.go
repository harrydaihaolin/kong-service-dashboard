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
	ServiceVersionName        string `gorm:"not null" json:"service_version_name"`
	ServiceVersionURL         string `gorm:"type:text" json:"service_version_url"`
	ServiceVersionDescription string `gorm:"type:text" json:"service_version_description"`
}

type User struct {
	gorm.Model

	Username    string      `gorm:"unique;not null" json:"username"`
	Password    string      `gorm:"not null" json:"password"`
	UserProfile UserProfile `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user_profile,omitempty"`
	Role        string      `gorm:"not null" json:"role"`
}

type UserProfile struct {
	gorm.Model

	UserID    uint   `gorm:"unique;not null" json:"user_id"`
	FirstName string `gorm:"type:varchar(255)" json:"first_name"`
	LastName  string `gorm:"type:varchar(255)" json:"last_name"`
	Email     string `gorm:"not null" json:"email"`
}
