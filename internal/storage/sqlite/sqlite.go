package sqlite

import (
	"fmt"

	"github.com/priyanka-16/attendance-app-backEnd/internal/config"
	"github.com/priyanka-16/attendance-app-backEnd/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Sqlite struct {
	Db *gorm.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := gorm.Open(sqlite.Open(cfg.StoragePath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}
	// AutoMigrate will create/update all tables & relationships
	err = db.AutoMigrate(
		&models.User{},
		&models.UserOTP{},
		&models.UserStudent{},
		&models.UserTeacher{},
		&models.School{},
		&models.SchoolGrade{},
		&models.SchoolGradeSection{},
		&models.Attendance{},
	)
	if err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return &Sqlite{Db: db}, nil
}

// ---------------- USER ----------------
func (s *Sqlite) CreateUser(u *models.User) (uint, error) {
	if err := s.Db.Create(u).Error; err != nil {
		return 0, err
	}
	return u.ID, nil
}

func (s *Sqlite) GetUserByMobile(mobile string) (*models.User, error) {
	var user models.User
	if err := s.Db.Where("mobile = ?", mobile).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Sqlite) GetUserById(id uint) (*models.User, error) {
	var user models.User
	if err := s.Db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// func (s *Sqlite) GetUsersList() ([]models.User, error) {
// 	var users []models.User
// 	if err := s.Db.Find(&users).Error; err != nil {
// 		return nil, err
// 	}
// 	return users, nil
// }

// ---------------- USER OTP ----------------
func (s *Sqlite) CreateUserOTP(otp *models.UserOTP) (uint, error) {
	if err := s.Db.Create(otp).Error; err != nil {
		return 0, err
	}
	return otp.ID, nil
}

func (s *Sqlite) GetUserOTPByMobile(mobile string) (*models.UserOTP, error) {
	var otp models.UserOTP
	if err := s.Db.Where("mobile = ?", mobile).Last(&otp).Error; err != nil {
		return nil, err
	}
	return &otp, nil
}

// func (s *Sqlite) GetUserOTPList() ([]models.UserOTP, error) {
// 	var otps []models.UserOTP
// 	if err := s.Db.Find(&otps).Error; err != nil {
// 		return nil, err
// 	}
// 	return otps, nil
// }

// // ---------------- USER STUDENT ----------------
func (s *Sqlite) CreateUserStudent(userStudent *models.UserStudent) (uint, error) {
	if err := s.Db.Create(userStudent).Error; err != nil {
		return 0, err
	}
	return userStudent.ID, nil
}

func (s *Sqlite) GetUserStudentById(id uint) (*models.UserStudent, error) {
	var userStudent models.UserStudent
	if err := s.Db.Preload("User").First(&userStudent, id).Error; err != nil {
		return nil, err
	}
	return &userStudent, nil
}
func (s *Sqlite) UpdateUserStudent(stu models.UserStudent) error {
	return s.Db.Model(&models.UserStudent{}).Where("user_id = ?", stu.UserID).Updates(stu).Error
}

// func (s *Sqlite) GetUserStudentsList() ([]models.UserStudent, error) {
// 	var userStudents []models.UserStudent
// 	if err := s.Db.Preload("User").Find(&userStudents).Error; err != nil {
// 		return nil, err
// 	}
// 	return userStudents, nil
// }

// // ---------------- USER TEACHER ----------------
func (s *Sqlite) CreateUserTeacher(userTeacher *models.UserTeacher) (uint, error) {
	if err := s.Db.Create(userTeacher).Error; err != nil {
		return 0, err
	}
	return userTeacher.ID, nil
}

func (s *Sqlite) GetUserTeacherById(id uint) (*models.UserTeacher, error) {
	var userTeacher models.UserTeacher
	if err := s.Db.Preload("User").Preload("School").First(&userTeacher, id).Error; err != nil {
		return nil, err
	}
	return &userTeacher, nil
}

func (s *Sqlite) UpdateUserTeacher(tea models.UserTeacher) error {
	return s.Db.Model(&models.UserTeacher{}).Where("user_id = ?", tea.UserID).Updates(tea).Error
}

// func (s *Sqlite) GetUserTeachersList() ([]models.UserTeacher, error) {
// 	var userTeachers []models.UserTeacher
// 	if err := s.Db.Preload("User").Preload("School").Find(&userTeachers).Error; err != nil {
// 		return nil, err
// 	}
// 	return userTeachers, nil
// }

// // ---------------- SCHOOL ----------------
// func (s *Sqlite) CreateSchool(school *models.School) (uint, error) {
// 	if err := s.Db.Create(school).Error; err != nil {
// 		return 0, err
// 	}
// 	return school.ID, nil
// }

// func (s *Sqlite) GetSchoolById(id uint) (*models.School, error) {
// 	var school models.School
// 	if err := s.Db.First(&school, id).Error; err != nil {
// 		return nil, err
// 	}
// 	return &school, nil
// }

// func (s *Sqlite) GetSchoolsList() ([]models.School, error) {
// 	var schools []models.School
// 	if err := s.Db.Find(&schools).Error; err != nil {
// 		return nil, err
// 	}
// 	return schools, nil
// }

// // ---------------- SCHOOL GRADE ----------------
// func (s *Sqlite) CreateSchoolGrade(schoolGrade *models.SchoolGrade) (uint, error) {
// 	if err := s.Db.Create(schoolGrade).Error; err != nil {
// 		return 0, err
// 	}
// 	return schoolGrade.ID, nil
// }

// func (s *Sqlite) GetSchoolGradeById(id uint) (*models.SchoolGrade, error) {
// 	var schoolGrade models.SchoolGrade
// 	if err := s.Db.Preload("School").First(&schoolGrade, id).Error; err != nil {
// 		return nil, err
// 	}
// 	return &schoolGrade, nil
// }

// func (s *Sqlite) GetSchoolGradesList() ([]models.SchoolGrade, error) {
// 	var schoolGrades []models.SchoolGrade
// 	if err := s.Db.Preload("School").Find(&schoolGrades).Error; err != nil {
// 		return nil, err
// 	}
// 	return schoolGrades, nil
// }

// // ---------------- SCHOOL GRADE SECTION ----------------
// func (s *Sqlite) CreateSchoolGradeSection(schoolGradeSection *models.SchoolGradeSection) (uint, error) {
// 	if err := s.Db.Create(schoolGradeSection).Error; err != nil {
// 		return 0, err
// 	}
// 	return schoolGradeSection.ID, nil
// }

// func (s *Sqlite) GetSchoolGradeSectionById(id uint) (*models.SchoolGradeSection, error) {
// 	var schoolGradeSection models.SchoolGradeSection
// 	if err := s.Db.Preload("Grade.School").
// 		Preload("ClassTeacher.User").
// 		Preload("ClassTeacher.School").First(&schoolGradeSection, id).Error; err != nil {
// 		return nil, err
// 	}
// 	return &schoolGradeSection, nil
// }

// func (s *Sqlite) GetSchoolGradeSectionsList() ([]models.SchoolGradeSection, error) {
// 	var schoolGradeSections []models.SchoolGradeSection
// 	if err := s.Db.Preload("Grade.School").
// 		Preload("ClassTeacher.User").
// 		Preload("ClassTeacher.School").Find(&schoolGradeSections).Error; err != nil {
// 		return nil, err
// 	}
// 	return schoolGradeSections, nil
// }

// // ---------------- ATTENDANCE ----------------
// func (s *Sqlite) CreateAttendance(attendance *models.Attendance) (uint, error) {
// 	if err := s.Db.Create(attendance).Error; err != nil {
// 		return 0, err
// 	}
// 	return attendance.ID, nil
// }

// func (s *Sqlite) GetAttendanceById(id uint) (*models.Attendance, error) {
// 	var attendance models.Attendance
// 	if err := s.Db.Preload("Student.User").
// 		Preload("Teacher.User").
// 		Preload("Teacher.School").First(&attendance, id).Error; err != nil {
// 		return nil, err
// 	}
// 	return &attendance, nil
// }

// func (s *Sqlite) GetAttendancesList() ([]models.Attendance, error) {
// 	var attendances []models.Attendance
// 	if err := s.Db.Preload("Student.User").
// 		Preload("Teacher.User").
// 		Preload("Teacher.School").Find(&attendances).Error; err != nil {
// 		return nil, err
// 	}
// 	return attendances, nil
// }
