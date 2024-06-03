package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID             int
	Email          string
	HashedPassword []byte
	Name           string
	Created        time.Time
}

// define a new UserModel type which wraps a database connection pool
type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) Insert(name, email, password string) error {
	return nil
}

func (u *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}
func (u *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
