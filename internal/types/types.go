package types

import "time"

type Student struct {
	ID       int64  `json:"id"`
	Name     string `json:"name" validate:"required"`
	RollNo   int    `json:"rollNo" validate:"required"`
	SchoolID int    `json:"schoolID" validate:"required"`
	Grade    string `json:"grade" validate:"required"`
}
type Teacher struct {
	ID       int      `json:"id"`
	Name     string   `json:"name" validate:"required"`
	Mobile   int      `json:"mobile" validate:"required"`
	SchoolID int      `json:"schoolID" validate:"required"`
	Grades   []string `json:"grades" validate:"required"`
}
type AttendanceRecord struct {
	ID        int       `json:"id"`
	StudentID int       `json:"studentID" validate:"required"`
	Date      time.Time `json:"date" validate:"required"`
	Status    string    `json:"status" validate:"required"`
	MarkedBy  int       `json:"markedBy" validate:"required"`
	MarkedAt  time.Time `json:"markedAt" validate:"required"`
}
type School struct {
	ID     int    `json:"id"`
	Name   string `json:"name" validate:"required"`
	Adress string `json:"adress" validate:"required"`
}
