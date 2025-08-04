package storage

import "github.com/priyanka-16/attendance-app-backEnd/internal/types"

type Storage interface {
	CreateUser(user types.User) (int64, error)
	GetUserById(id int64) (types.User, error)
	GetUsersList() ([]types.User, error)

	CreateUserOTP(otp types.UserOTP) (int64, error)
	GetUserOTPById(id int64) (types.UserOTP, error)
	GetUserOTPList() ([]types.UserOTP, error)

	CreateUserStudent(s types.UserStudent) (int64, error)
	GetUserStudentById(id int64) (types.UserStudent, error)
	GetUserStudentsList() ([]types.UserStudent, error)

	CreateUserTeacher(t types.UserTeacher) (int64, error)
	GetUserTeacherById(id int64) (types.UserTeacher, error)
	GetUserTeachersList() ([]types.UserTeacher, error)

	CreateSchool(sc types.School) (int64, error)
	GetSchoolById(id int64) (types.School, error)
	GetSchoolsList() ([]types.School, error)

	CreateSchoolGrade(g types.SchoolGrade) (int64, error)
	GetSchoolGradeById(id int64) (types.SchoolGrade, error)
	GetSchoolGradesList() ([]types.SchoolGrade, error)

	CreateSchoolGradeSection(s types.SchoolGradeSection) (int64, error)
	GetSchoolGradeSectionById(id int64) (types.SchoolGradeSection, error)
	GetSchoolGradeSectionsList() ([]types.SchoolGradeSection, error)

	CreateAttendance(a types.Attendance) (int64, error)
	GetAttendanceById(id int64) (types.Attendance, error)
	GetAttendancesList() ([]types.Attendance, error)
}
