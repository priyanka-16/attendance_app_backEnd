package storage

import "github.com/priyanka-16/attendance-app-backEnd/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
}
