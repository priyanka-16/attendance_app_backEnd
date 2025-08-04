package types

import "time"

type User struct {
	ID        int64     `json:"id"`
	Mobile    string    `json:"mobile" validate:"required"`
	LoginHash string    `json:"loginHash" validate:"required"`
	Password  string    `json:"password" validate:"required"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserOTP struct {
	ID        int64     `json:"id"`
	Mobile    string    `json:"mobile" validate:"required"`
	Code      string    `json:"code" validate:"required"`
	IsUsed    bool      `json:"isUsed"`
	ExpiresAt time.Time `json:"expiresAt" validate:"required"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserStudent struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"userID" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserTeacher struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"userID" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	SchoolID  int64     `json:"schoolID" validate:"required"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type School struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name" validate:"required"`
	Address   string    `json:"address"`
	District  string    `json:"district" validate:"required"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type SchoolGrade struct {
	ID        int64     `json:"id"`
	SchoolID  int64     `json:"schoolID" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	Slug      string    `json:"slug" validate:"required"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type SchoolGradeSection struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name" validate:"required"`
	FullName       string    `json:"fullName" validate:"required"`
	Slug           string    `json:"slug" validate:"required"`
	GradeID        int64     `json:"gradeID" validate:"required"`
	ClassTeacherID int64     `json:"classTeacherID"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type Attendance struct {
	ID        int64     `json:"id"`
	StudentID int64     `json:"studentID" validate:"required"`
	Date      time.Time `json:"date" validate:"required"`
	Status    string    `json:"status" validate:"required"`
	TakenBy   int64     `json:"takenBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
