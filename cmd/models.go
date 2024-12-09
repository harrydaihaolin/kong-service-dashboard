package main

type Service struct {
	ServiceID          uint             `gorm:"primaryKey" json:"service_id"`
	ServiceName        string           `gorm:"type:varchar(255);unique;not null" json:"service_name"`
	ServiceDescription string           `gorm:"type:text" json:"service_description"`
	Versions           []ServiceVersion `gorm:"foreignKey:ServiceID" json:"service_versions,omitempty"`
}

type ServiceVersion struct {
	ServiceVersionID          uint   `gorm:"primaryKey" json:"service_version_id"`
	ServiceID                 uint   `gorm:"not null" json:"service_id"`
	ServiceVersionName        string `gorm:"type:varchar(255);unique" json:"service_version_name"`
	ServiceVersionURL         string `gorm:"type:text" json:"service_version_url"`
	ServiceVersionDescription string `gorm:"type:text" json:"service_version_description"`
}

type User struct {
	UserID      uint        `gorm:"primaryKey" json:"user_id"`
	Username    string      `gorm:"type:varchar(255);unique;not null" json:"username"`
	UserProfile UserProfile `gorm:"foreignKey:UserID" json:"user_profile,omitempty"`
	UserRole    UserRole    `gorm:"foreignKey:UserID" json:"user_role,omitempty"`
}

type UserProfile struct {
	UserProfileID uint   `gorm:"primaryKey" json:"user_profile_id"`
	UserID        uint   `gorm:"not null;unique" json:"user_id"`
	FirstName     string `gorm:"type:varchar(255)" json:"first_name"`
	LastName      string `gorm:"type:varchar(255)" json:"last_name"`
	Email         string `gorm:"type:varchar(255);unique;not null" json:"email"`
}

type UserRole struct {
	UserRoleID uint   `gorm:"primaryKey" json:"user_role_id"`
	UserID     uint   `gorm:"not null" json:"user_id"`
	Role       string `gorm:"type:varchar(255)" json:"role"`
}
