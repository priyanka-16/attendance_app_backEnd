package sqlite

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/priyanka-16/attendance-app-backEnd/internal/config"
	"github.com/priyanka-16/attendance-app-backEnd/internal/types"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	// USERS
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS user (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	mobile TEXT NOT NULL UNIQUE,
	loginHash TEXT NOT NULL,
	password TEXT NOT NULL,
	isActive BOOLEAN DEFAULT 1,
	createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
	updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		return nil, err
	}

	// USER_OTPS
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS user_otp (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	mobile TEXT NOT NULL,
	code TEXT NOT NULL,
	isUsed BOOLEAN DEFAULT 0,
	expiresAt DATETIME NOT NULL,
	createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (mobile) REFERENCES user (mobile) ON DELETE CASCADE
	)`)
	if err != nil {
		return nil, err
	}

	// USER_STUDENTS
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS user_student (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	userID INTEGER NOT NULL,
	name TEXT NOT NULL,
	isActive BOOLEAN DEFAULT 1,
	createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
	updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (userID) REFERENCES user (id) ON DELETE CASCADE
	)`)
	if err != nil {
		return nil, err
	}

	// USER_TEACHERS
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS user_teacher (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	userID INTEGER NOT NULL,
	name TEXT NOT NULL,
	schoolID INTEGER NOT NULL,
	isActive BOOLEAN DEFAULT 1,
	createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
	updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (userID) REFERENCES user (id) ON DELETE CASCADE,
	FOREIGN KEY (schoolID) REFERENCES school (id) ON DELETE CASCADE
	)`)
	if err != nil {
		return nil, err
	}

	// SCHOOLS
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS school (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	address TEXT,
	district TEXT NOT NULL,
	phone TEXT,
	email TEXT,
	createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
	updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return nil, err
	}

	// SCHOOL_GRADES
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS school_grade (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	schoolID INTEGER NOT NULL,
	name TEXT NOT NULL,
	slug TEXT NOT NULL,
	createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
	updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (schoolID) REFERENCES school (id) ON DELETE CASCADE
	)`)
	if err != nil {
		return nil, err
	}

	// SCHOOL_GRADE_SECTIONS
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS school_grade_section (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	fullName TEXT NOT NULL,
	slug TEXT NOT NULL,
	gradeID INTEGER NOT NULL,
	classTeacherID INTEGER,
	createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
	updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (gradeID) REFERENCES school_grade (id) ON DELETE CASCADE
	)`)
	if err != nil {
		return nil, err
	}

	// ATTENDANCE
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS attendance (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		studentID INTEGER NOT NULL,
		date TEXT NOT NULL,
		status TEXT NOT NULL,
		takenBy INTEGER,
		createdAt TEXT NOT NULL,
		updatedAt TEXT NOT NULL
	)`)
	if err != nil {
		return nil, err
	}

	return &Sqlite{Db: db}, nil
}

// ---------------- USER ----------------
func (s *Sqlite) CreateUser(u types.User) (int64, error) {
	stmt, err := s.Db.Prepare(`INSERT INTO user (mobile, loginHash, password, isActive, createdAt, updatedAt) VALUES (?,?,?,?,?,?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(u.Mobile, u.LoginHash, u.Password, u.IsActive, u.CreatedAt.Format(time.RFC3339), u.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (s *Sqlite) GetUserById(id int64) (types.User, error) {
	var u types.User
	stmt := `SELECT id, mobile, loginHash, password, isActive, createdAt, updatedAt FROM user WHERE id = ?`
	row := s.Db.QueryRow(stmt, id)
	var createdAt, updatedAt string
	err := row.Scan(&u.ID, &u.Mobile, &u.LoginHash, &u.Password, &u.IsActive, &createdAt, &updatedAt)
	if err != nil {
		return u, err
	}
	u.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	u.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	return u, nil
}

func (s *Sqlite) GetUsersList() ([]types.User, error) {
	rows, err := s.Db.Query("SELECT id, mobile, loginHash, password, isActive, createdAt, updatedAt FROM user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []types.User
	for rows.Next() {
		var u types.User
		var createdAt, updatedAt string

		err := rows.Scan(&u.ID, &u.Mobile, &u.LoginHash, &u.Password, &u.IsActive, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}

		u.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		u.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

		users = append(users, u)
	}
	return users, nil
}

// ---------------- USER OTP ----------------
func (s *Sqlite) CreateUserOTP(o types.UserOTP) (int64, error) {
	stmt, err := s.Db.Prepare(`INSERT INTO user_otp (mobile, code, isUsed, expiresAt, createdAt) VALUES (?,?,?,?,?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(o.Mobile, o.Code, o.IsUsed, o.ExpiresAt.Format(time.RFC3339), o.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (s *Sqlite) GetUserOTPById(id int64) (types.UserOTP, error) {
	var otp types.UserOTP
	stmt := `SELECT id, mobile, code, isUsed, expiresAt, createdAt FROM user_otp WHERE id = ?`
	row := s.Db.QueryRow(stmt, id)

	var expiresAt, createdAt string
	err := row.Scan(&otp.ID, &otp.Mobile, &otp.Code, &otp.IsUsed, &expiresAt, &createdAt)
	if err != nil {
		return otp, err
	}

	otp.ExpiresAt, _ = time.Parse(time.RFC3339, expiresAt)
	otp.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)

	return otp, nil
}

func (s *Sqlite) GetUserOTPList() ([]types.UserOTP, error) {
	rows, err := s.Db.Query("SELECT id, mobile, code, isUsed, expiresAt, createdAt FROM user_otp")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var otps []types.UserOTP
	for rows.Next() {
		var otp types.UserOTP
		var expiresAt, createdAt string

		err := rows.Scan(&otp.ID, &otp.Mobile, &otp.Code, &otp.IsUsed, &expiresAt, &createdAt)
		if err != nil {
			return nil, err
		}

		otp.ExpiresAt, _ = time.Parse(time.RFC3339, expiresAt)
		otp.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)

		otps = append(otps, otp)
	}
	return otps, nil
}

// ---------------- USER STUDENT ----------------
func (s *Sqlite) CreateUserStudent(us types.UserStudent) (int64, error) {
	stmt, err := s.Db.Prepare(`INSERT INTO user_student (userID, name, isActive, createdAt, updatedAt) VALUES (?,?,?,?,?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(us.UserID, us.Name, us.IsActive, us.CreatedAt.Format(time.RFC3339), us.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (s *Sqlite) GetUserStudentById(id int64) (types.UserStudent, error) {
	var us types.UserStudent
	stmt := `SELECT id, userID, name, isActive, createdAt, updatedAt FROM user_student WHERE id = ?`
	row := s.Db.QueryRow(stmt, id)

	var createdAt, updatedAt string
	err := row.Scan(&us.ID, &us.UserID, &us.Name, &us.IsActive, &createdAt, &updatedAt)
	if err != nil {
		return us, err
	}

	us.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	us.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return us, nil
}

func (s *Sqlite) GetUserStudentsList() ([]types.UserStudent, error) {
	rows, err := s.Db.Query("SELECT id, userID, name, isActive, createdAt, updatedAt FROM user_student")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []types.UserStudent
	for rows.Next() {
		var us types.UserStudent
		var createdAt, updatedAt string

		err := rows.Scan(&us.ID, &us.UserID, &us.Name, &us.IsActive, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}

		us.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		us.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

		students = append(students, us)
	}
	return students, nil
}

// ---------------- USER TEACHER ----------------
func (s *Sqlite) CreateUserTeacher(ut types.UserTeacher) (int64, error) {
	stmt, err := s.Db.Prepare(`INSERT INTO user_teacher (userID, name, schoolID, isActive, createdAt, updatedAt) VALUES (?,?,?,?,?,?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(ut.UserID, ut.Name, ut.SchoolID, ut.IsActive, ut.CreatedAt.Format(time.RFC3339), ut.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (s *Sqlite) GetUserTeacherById(id int64) (types.UserTeacher, error) {
	var ut types.UserTeacher
	stmt := `SELECT id, userID, name, schoolID, isActive, createdAt, updatedAt FROM user_teacher WHERE id = ?`
	row := s.Db.QueryRow(stmt, id)

	var createdAt, updatedAt string
	err := row.Scan(&ut.ID, &ut.UserID, &ut.Name, &ut.SchoolID, &ut.IsActive, &createdAt, &updatedAt)
	if err != nil {
		return ut, err
	}

	ut.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	ut.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return ut, nil
}

func (s *Sqlite) GetUserTeachersList() ([]types.UserTeacher, error) {
	rows, err := s.Db.Query("SELECT id, userID, name, schoolID, isActive, createdAt, updatedAt FROM user_teacher")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teachers []types.UserTeacher
	for rows.Next() {
		var ut types.UserTeacher
		var createdAt, updatedAt string

		err := rows.Scan(&ut.ID, &ut.UserID, &ut.Name, &ut.SchoolID, &ut.IsActive, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}

		ut.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		ut.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

		teachers = append(teachers, ut)
	}
	return teachers, nil
}

// ---------------- SCHOOL ----------------
func (s *Sqlite) CreateSchool(sc types.School) (int64, error) {
	stmt, err := s.Db.Prepare(`INSERT INTO school (name, address, district, phone, email, createdAt, updatedAt) VALUES (?,?,?,?,?,?,?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(sc.Name, sc.Address, sc.District, sc.Phone, sc.Email, sc.CreatedAt.Format(time.RFC3339), sc.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (s *Sqlite) GetSchoolById(id int64) (types.School, error) {
	var sc types.School
	stmt := `SELECT id, name, address, district, phone, email, createdAt, updatedAt FROM school WHERE id = ?`
	row := s.Db.QueryRow(stmt, id)

	var createdAt, updatedAt string
	err := row.Scan(&sc.ID, &sc.Name, &sc.Address, &sc.District, &sc.Phone, &sc.Email, &createdAt, &updatedAt)
	if err != nil {
		return sc, err
	}

	sc.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	sc.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return sc, nil
}

func (s *Sqlite) GetSchoolsList() ([]types.School, error) {
	rows, err := s.Db.Query("SELECT id, name, address, district, phone, email, createdAt, updatedAt FROM school")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schools []types.School
	for rows.Next() {
		var sc types.School
		var createdAt, updatedAt string

		err := rows.Scan(&sc.ID, &sc.Name, &sc.Address, &sc.District, &sc.Phone, &sc.Email, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}

		sc.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		sc.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

		schools = append(schools, sc)
	}
	return schools, nil
}

// ---------------- SCHOOL GRADE ----------------
func (s *Sqlite) CreateSchoolGrade(g types.SchoolGrade) (int64, error) {
	stmt, err := s.Db.Prepare(`INSERT INTO school_grade (schoolID, name, slug, createdAt, updatedAt) VALUES (?,?,?,?,?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(g.SchoolID, g.Name, g.Slug, g.CreatedAt.Format(time.RFC3339), g.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
func (s *Sqlite) GetSchoolGradeById(id int64) (types.SchoolGrade, error) {
	var g types.SchoolGrade
	stmt := `SELECT id, schoolID, name, slug, createdAt, updatedAt FROM school_grade WHERE id = ?`
	row := s.Db.QueryRow(stmt, id)

	var createdAt, updatedAt string
	err := row.Scan(&g.ID, &g.SchoolID, &g.Name, &g.Slug, &createdAt, &updatedAt)
	if err != nil {
		return g, err
	}

	g.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	g.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return g, nil
}

func (s *Sqlite) GetSchoolGradesList() ([]types.SchoolGrade, error) {
	rows, err := s.Db.Query("SELECT id, schoolID, name, slug, createdAt, updatedAt FROM school_grade")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var grades []types.SchoolGrade
	for rows.Next() {
		var g types.SchoolGrade
		var createdAt, updatedAt string

		err := rows.Scan(&g.ID, &g.SchoolID, &g.Name, &g.Slug, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}

		g.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		g.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

		grades = append(grades, g)
	}
	return grades, nil
}

// ---------------- SCHOOL GRADE SECTION ----------------
func (s *Sqlite) CreateSchoolGradeSection(sec types.SchoolGradeSection) (int64, error) {
	stmt, err := s.Db.Prepare(`INSERT INTO school_grade_section (name, fullName, slug, gradeID, classTeacherID, createdAt, updatedAt) VALUES (?,?,?,?,?,?,?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(sec.Name, sec.FullName, sec.Slug, sec.GradeID, sec.ClassTeacherID, sec.CreatedAt.Format(time.RFC3339), sec.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
func (s *Sqlite) GetSchoolGradeSectionById(id int64) (types.SchoolGradeSection, error) {
	var sec types.SchoolGradeSection
	stmt := `SELECT id, name, fullName, slug, gradeID, classTeacherID, createdAt, updatedAt FROM school_grade_section WHERE id = ?`
	row := s.Db.QueryRow(stmt, id)

	var createdAt, updatedAt string
	err := row.Scan(&sec.ID, &sec.Name, &sec.FullName, &sec.Slug, &sec.GradeID, &sec.ClassTeacherID, &createdAt, &updatedAt)
	if err != nil {
		return sec, err
	}

	sec.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	sec.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return sec, nil
}

func (s *Sqlite) GetSchoolGradeSectionsList() ([]types.SchoolGradeSection, error) {
	rows, err := s.Db.Query("SELECT id, name, fullName, slug, gradeID, classTeacherID, createdAt, updatedAt FROM school_grade_section")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sections []types.SchoolGradeSection
	for rows.Next() {
		var sec types.SchoolGradeSection
		var createdAt, updatedAt string

		err := rows.Scan(&sec.ID, &sec.Name, &sec.FullName, &sec.Slug, &sec.GradeID, &sec.ClassTeacherID, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}

		sec.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		sec.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

		sections = append(sections, sec)
	}
	return sections, nil
}

// ---------------- ATTENDANCE ----------------
func (s *Sqlite) CreateAttendance(a types.Attendance) (int64, error) {
	stmt, err := s.Db.Prepare(`INSERT INTO attendance (studentID, date, status, takenBy, createdAt, updatedAt) VALUES (?,?,?,?,?,?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(a.StudentID, a.Date.Format(time.RFC3339), a.Status, a.TakenBy, a.CreatedAt.Format(time.RFC3339), a.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
func (s *Sqlite) GetAttendanceById(id int64) (types.Attendance, error) {
	var a types.Attendance
	stmt := `SELECT id, studentID, date, status, takenBy, createdAt, updatedAt FROM attendance WHERE id = ?`
	row := s.Db.QueryRow(stmt, id)

	var date, createdAt, updatedAt string
	err := row.Scan(&a.ID, &a.StudentID, &date, &a.Status, &a.TakenBy, &createdAt, &updatedAt)
	if err != nil {
		return a, err
	}

	a.Date, _ = time.Parse(time.RFC3339, date)
	a.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	a.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return a, nil
}

func (s *Sqlite) GetAttendancesList() ([]types.Attendance, error) {
	rows, err := s.Db.Query("SELECT id, studentID, date, status, takenBy, createdAt, updatedAt FROM attendance")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attendances []types.Attendance
	for rows.Next() {
		var a types.Attendance
		var date, createdAt, updatedAt string

		err := rows.Scan(&a.ID, &a.StudentID, &date, &a.Status, &a.TakenBy, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}

		a.Date, _ = time.Parse(time.RFC3339, date)
		a.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		a.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

		attendances = append(attendances, a)
	}
	return attendances, nil
}
