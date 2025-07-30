package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3" //this package is required for initialization but is not called (used behind scenes), hence prefixed with _
	"github.com/priyanka-16/attendance-app-backEnd/internal/config"
	"github.com/priyanka-16/attendance-app-backEnd/internal/types"
)

type Sqlite struct {
	Db *sql.DB
}

//now to create an instance of struct we do like below

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)

	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	rollNo INTEGER NOT NULL,
	schoolID INTEGER NOT NULL,
	grade TEXT NOT NULL,
	FOREIGN KEY (schoolID) REFERENCES schools(id)
	)`)

	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS teachers(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	mobile INTEGER NOT NULL,
	schoolID INTEGER NOT NULL,
	grades TEXT NOT NULL,
	FOREIGN KEY (schoolID) REFERENCES schools(id)
	)`)

	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS attendanceRecords(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	studentID INTEGER NOT NULL,
	date TEXT NOT NULL,
	status TEXT CHECK (status IN ('present', 'absent')) NOT NULL,
	markedBy int NOT NULL,
	markedAt TEXT NOT NULL,
	FOREIGN KEY (studentID) REFERENCES students(id),
    FOREIGN KEY (markedBy) REFERENCES teachers(id)
	)`)

	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS schools(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	adress TEXT NOT NULL
	)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil

}

func (s *Sqlite) CreateStudent(name string, rollNo int, schoolID int, grade string) (int64, error) {
	//way to declare Sqlite struct implements Storage interface by defining common method

	stmt, err := s.Db.Prepare("INSERT INTO students (name, rollNo, schoolID, grade) VALUES (?,?,?,?)") //we put placeholders first and values later to save from SQL injection
	if err != nil {
		return 0, nil
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, rollNo, schoolID, grade)
	if err != nil {
		return 0, nil
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	return lastId, nil
}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {

	stmt, err := s.Db.Prepare("SELECT * FROM students where id = ? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()

	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.ID, &student.Name, &student.RollNo, &student.SchoolID, &student.Grade)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("No student found with id %s", id)
		}

		return types.Student{}, fmt.Errorf("query Error:%w", err)
	}

	return student, nil
}

func (s *Sqlite) GetStudentsList() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT * FROM students")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []types.Student
	for rows.Next() {
		var student types.Student
		err := rows.Scan(&student.ID, &student.Name, &student.RollNo, &student.SchoolID, &student.Grade)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}
	return students, nil
}

func (s *Sqlite) CreateTeacher(name string, mobile int, schoolID int, grades []string) (int64, error) {
	//way to declare Sqlite struct implements Storage interface by defining common method
	stmt, err := s.Db.Prepare("INSERT INTO teachers (name, mobile, schoolID, grades) VALUES (?,?,?,?)") //we put placeholders first and values later to save from SQL injection
	if err != nil {
		return 0, nil
	}
	defer stmt.Close()

	gradesJSON, err := json.Marshal(grades)
	if err != nil {
		return 0, nil
	}

	result, err := stmt.Exec(name, mobile, schoolID, string(gradesJSON))
	if err != nil {
		return 0, nil
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	return lastId, nil
}

func (s *Sqlite) GetTeacherById(id int64) (types.Teacher, error) {

	stmt, err := s.Db.Prepare("SELECT * FROM teachers where id = ? LIMIT 1")
	if err != nil {
		return types.Teacher{}, err
	}
	defer stmt.Close()

	var teacher types.Teacher
	var gradesJSON string

	err = stmt.QueryRow(id).Scan(&teacher.ID, &teacher.Name, &teacher.Mobile, &teacher.SchoolID, &gradesJSON)
	json.Unmarshal([]byte(gradesJSON), &teacher.Grades)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Teacher{}, fmt.Errorf("No teacher found with id %s", id)
		}

		return types.Teacher{}, fmt.Errorf("query Error:%w", err)
	}

	return teacher, nil
}

func (s *Sqlite) GetTeachersList() ([]types.Teacher, error) {
	stmt, err := s.Db.Prepare("SELECT * FROM teachers")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teachers []types.Teacher
	for rows.Next() {
		var teacher types.Teacher
		var gradesJSON string
		err := rows.Scan(&teacher.ID, &teacher.Name, &teacher.Mobile, &teacher.SchoolID, &gradesJSON)
		json.Unmarshal([]byte(gradesJSON), &teacher.Grades)
		if err != nil {
			return nil, err
		}
		teachers = append(teachers, teacher)
	}
	return teachers, nil
}

func (s *Sqlite) CreateAttendanceRecord(record types.AttendanceRecord) (int64, error) {
	stmt, err := s.Db.Prepare(`
    INSERT INTO attendanceRecords 
    (studentID, date, status, markedBy, markedAt) 
    VALUES (?, ?, ?, ?, ?)
  `)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(
		record.StudentID,
		record.Date.Format(time.RFC3339),
		record.Status,
		record.MarkedBy,
		record.MarkedAt.Format(time.RFC3339),
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (s *Sqlite) GetAttendanceRecordByID(id int64) (types.AttendanceRecord, error) {
	stmt, err := s.Db.Prepare(`
    SELECT id, studentID, date, status, markedBy, markedAt 
    FROM attendanceRecords 
    WHERE id = ? LIMIT 1
  `)
	if err != nil {
		return types.AttendanceRecord{}, err
	}
	defer stmt.Close()

	var record types.AttendanceRecord
	var markedAtStr string
	var dateStr string

	err = stmt.QueryRow(id).Scan(
		&record.ID,
		&record.StudentID,
		&dateStr,
		&record.Status,
		&record.MarkedBy,
		&markedAtStr,
	)
	if err != nil {
		return types.AttendanceRecord{}, err
	}

	record.Date, err = time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return types.AttendanceRecord{}, fmt.Errorf("invalid date format: %w", err)
	}

	record.MarkedAt, err = time.Parse(time.RFC3339, markedAtStr)
	if err != nil {
		return types.AttendanceRecord{}, fmt.Errorf("invalid markedAt format: %w", err)
	}
	return record, nil
}

func (s *Sqlite) GetAttendanceRecordsList() ([]types.AttendanceRecord, error) {
	stmt, err := s.Db.Prepare(`
    SELECT id, studentID, date, status, markedBy, markedAt 
    FROM attendanceRecords
  `)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []types.AttendanceRecord

	for rows.Next() {
		var record types.AttendanceRecord
		var markedAtStr string
		var dateStr string

		err := rows.Scan(
			&record.ID,
			&record.StudentID,
			&dateStr,
			&record.Status,
			&record.MarkedBy,
			&markedAtStr,
		)
		if err != nil {
			return nil, err
		}

		record.Date, err = time.Parse(time.RFC3339, dateStr)
		if err != nil {
			return []types.AttendanceRecord{}, fmt.Errorf("invalid date format: %w", err)
		}

		record.MarkedAt, err = time.Parse(time.RFC3339, markedAtStr)
		if err != nil {
			return []types.AttendanceRecord{}, fmt.Errorf("invalid markedAt format: %w", err)
		}
		records = append(records, record)
	}

	return records, nil
}

func (s *Sqlite) CreateSchool(school types.School) (int64, error) {
	stmt, err := s.Db.Prepare(`
    INSERT INTO schools 
    (name, adress) 
    VALUES (?, ?)
  `)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(
		school.Name,
		school.Adress,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (s *Sqlite) GetSchoolByID(id int64) (types.School, error) {
	stmt, err := s.Db.Prepare(`
    SELECT id, name, adress 
    FROM schools 
    WHERE id = ? LIMIT 1
  `)
	if err != nil {
		return types.School{}, err
	}
	defer stmt.Close()

	var school types.School

	err = stmt.QueryRow(id).Scan(
		&school.ID,
		&school.Name,
		&school.Adress,
	)
	if err != nil {
		return types.School{}, err
	}

	return school, nil
}

func (s *Sqlite) GetSchoolsList() ([]types.School, error) {
	stmt, err := s.Db.Prepare(`
    SELECT id, name, adress 
    FROM schools
  `)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schools []types.School

	for rows.Next() {
		var school types.School

		err := rows.Scan(
			&school.ID,
			&school.Name,
			&school.Adress,
		)
		if err != nil {
			return nil, err
		}
		schools = append(schools, school)
	}

	return schools, nil
}
