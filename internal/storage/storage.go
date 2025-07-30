package storage

import "github.com/priyanka-16/attendance-app-backEnd/internal/types"

type Storage interface {
	CreateStudent(name string, rollNo int, schoolID int, grade string) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetStudentsList() ([]types.Student, error)
	CreateTeacher(name string, mobile int, schoolID int, grades []string) (int64, error)
	GetTeacherById(id int64) (types.Teacher, error)
	GetTeachersList() ([]types.Teacher, error)
	CreateAttendanceRecord(record types.AttendanceRecord) (int64, error)
	GetAttendanceRecordByID(id int64) (types.AttendanceRecord, error)
	GetAttendanceRecordsList() ([]types.AttendanceRecord, error)
	CreateSchool(school types.School) (int64, error)
	GetSchoolByID(id int64) (types.School, error)
	GetSchoolsList() ([]types.School, error)
}
