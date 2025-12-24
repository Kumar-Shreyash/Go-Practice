package storage

import "github.com/kumar-shreyash/students-api/internal/types"


type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)//this is a method with its signature
	GetStudentById(id int64) (types.Student, error)
	GetStudents() ([]types.Student, error)
}
