package models

import (
	"gorm.io/gorm"
	"time"
)

// GormModel  definition
type GormModel struct {
	ID        uint           `gorm:"primaryKey" json:"ID"`
	CreatedAt time.Time      `gorm:"index" json:"created"`
	UpdatedAt time.Time      `json:"updated"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted"`
}

type User struct {
	GormModel
	Name     string  `gorm:"not null" json:"name"`
	Email    string  `gorm:"uniqueIndex;not null;size:255" json:"eMailAddress"`
	Password string  `json:"password"`
	Address  string  `json:"address"`
	UserType string  `json:"userType"`
	Profile  Profile `gorm:"foreignKey:ApplicantID"`
	Tokens   []Token
}

// Profile represents the profile model in the database.
type Profile struct {
	gorm.Model
	ApplicantID       uint `gorm:"not null"`
	ResumeFileAddress string
	Skills            string
	Education         string
	Experience        string
	Name              string `gorm:"not null"`
	Email             string `gorm:"not null"`
	Phone             string
}
type Job struct {
	gorm.Model
	Title             string    `gorm:"not null"`
	Description       string    `gorm:"not null"`
	PostedOn          time.Time `gorm:"not null"`
	TotalApplications int
	CompanyName       string `gorm:"not null"`
	PostedByID        uint   `gorm:"not null"`
	PostedBy          User   `gorm:"foreignKey:PostedByID"`
}

// ValidateUserType validates the UserType field in User model.
func (u *User) ValidateUserType() error {
	if u.UserType != "applicant" && u.UserType != "admin" {
		return gorm.ErrInvalidField
	}
	return nil
}

// BeforeSave hook to validate the UserType field.
func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	if err := u.ValidateUserType(); err != nil {
		return err
	}
	return nil
}
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	// Set the CreatedAt and UpdatedAt timestamps
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	// Update the UpdatedAt timestamp
	u.UpdatedAt = time.Now()
	return nil
}
