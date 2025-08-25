package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Mobile    string `gorm:"uniqueIndex;not null"`
	LoginHash string `gorm:"not null"`
	Password  string `gorm:"not null"`
	IsActive  bool   `gorm:"default:true"`
}

type UserOTP struct {
	gorm.Model
	Mobile    string `gorm:"not null;index"`
	Code      string `gorm:"not null"`
	IsUsed    bool   `gorm:"default:false"`
	ExpiresAt int64  `gorm:"not null"`
}

type UserStudent struct {
	gorm.Model
	UserID   uint   `gorm:"not null;index"`
	Name     string `gorm:"not null"`
	IsActive bool   `gorm:"default:true"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

type UserTeacher struct {
	gorm.Model
	UserID   uint   `gorm:"not null;index"`
	Name     string `gorm:"not null"`
	SchoolID uint   `gorm:"not null;index"`
	IsActive bool   `gorm:"default:true"`

	User   User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	School School `gorm:"foreignKey:SchoolID;constraint:OnDelete:CASCADE;"`
}

type School struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Address  string
	District string `gorm:"not null"`
	Phone    string
	Email    string
}

type SchoolGrade struct {
	gorm.Model
	SchoolID uint   `gorm:"not null;index"`
	Name     string `gorm:"not null"`
	Slug     string `gorm:"not null"`

	School School `gorm:"foreignKey:SchoolID;constraint:OnDelete:CASCADE;"`
}

type SchoolGradeSection struct {
	gorm.Model
	Name           string `gorm:"not null"`
	FullName       string `gorm:"not null"`
	Slug           string `gorm:"not null"`
	GradeID        uint   `gorm:"not null;index"`
	ClassTeacherID uint

	Grade        SchoolGrade `gorm:"foreignKey:GradeID;constraint:OnDelete:CASCADE;"`
	ClassTeacher UserTeacher `gorm:"foreignKey:ClassTeacherID;constraint:OnDelete:SET NULL;"`
}

type Attendance struct {
	gorm.Model
	StudentID uint   `gorm:"not null;index"`
	Date      string `gorm:"not null"`
	Status    string `gorm:"not null"`
	TakenBy   uint

	Student UserStudent `gorm:"foreignKey:StudentID;constraint:OnDelete:CASCADE;"`
	Teacher UserTeacher `gorm:"foreignKey:TakenBy;constraint:OnDelete:SET NULL;"`
}
