package storage

import "github.com/priyanka-16/attendance-app-backEnd/internal/models"

type Storage interface {
	// ---------------- USER ----------------
	CreateUser(user *models.User) (uint, error)
	// GetUserById(id uint) (*models.User, error)
	// GetUsersList() ([]models.User, error)

	// // ---------------- USER OTP ----------------
	CreateUserOTP(otp *models.UserOTP) (uint, error)
	GetUserOTPByMobile(mobile string) (*models.UserOTP, error)
	GetUserByMobile(mobile string) (*models.User, error)
	// GetUserOTPById(id uint) (*models.UserOTP, error)
	// GetUserOTPList() ([]models.UserOTP, error)

	// // ---------------- USER STUDENT ----------------
	// CreateUserStudent(s *models.UserStudent) (uint, error)
	// GetUserStudentById(id uint) (*models.UserStudent, error)
	// GetUserStudentsList() ([]models.UserStudent, error)

	// // ---------------- USER TEACHER ----------------
	// CreateUserTeacher(t *models.UserTeacher) (uint, error)
	// GetUserTeacherById(id uint) (*models.UserTeacher, error)
	// GetUserTeachersList() ([]models.UserTeacher, error)

	// // ---------------- SCHOOL ----------------
	// CreateSchool(sc *models.School) (uint, error)
	// GetSchoolById(id uint) (*models.School, error)
	// GetSchoolsList() ([]models.School, error)

	// // ---------------- SCHOOL GRADE ----------------
	// CreateSchoolGrade(g *models.SchoolGrade) (uint, error)
	// GetSchoolGradeById(id uint) (*models.SchoolGrade, error)
	// GetSchoolGradesList() ([]models.SchoolGrade, error)

	// // ---------------- SCHOOL GRADE SECTION ----------------
	// CreateSchoolGradeSection(s *models.SchoolGradeSection) (uint, error)
	// GetSchoolGradeSectionById(id uint) (*models.SchoolGradeSection, error)
	// GetSchoolGradeSectionsList() ([]models.SchoolGradeSection, error)

	// // ---------------- ATTENDANCE ----------------
	// CreateAttendance(a *models.Attendance) (uint, error)
	// GetAttendanceById(id uint) (*models.Attendance, error)
	// GetAttendancesList() ([]models.Attendance, error)
}
