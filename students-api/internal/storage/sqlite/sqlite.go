package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/kumar-shreyash/students-api/internal/config"
	"github.com/kumar-shreyash/students-api/internal/types"
	_ "github.com/mattn/go-sqlite3" //this package is not getting used directly so we add a _ before this
)

type Sqlite struct { //for storing db connection
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) { //in go there is nothing as constructor so we use a convention and create a function named New

	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	//table creation
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT,
		age INTEGER

	)`)

	if err != nil {
		return nil, err
	}
	return &Sqlite{
		Db: db,
	}, nil
}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")

	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, email, age) //pass these fields in order or the data will get saved in different fields

	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return lastId, nil

	// return 0, nil
}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {

	stmt, err := s.Db.Prepare("SELECT * FROM students WHERE id = ? LIMIT 1")

	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()

	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)

	if err != nil {

		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id %v", fmt.Sprint(id))
		}
		return types.Student{}, fmt.Errorf("Query error: %w", err)
	}
	return student, nil
}

func (s *Sqlite) GetStudents() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM STUDENTS")

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

		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)

		if err != nil {
			return nil, err
		}
		students = append(students, student)

	}

	return students, nil
}
