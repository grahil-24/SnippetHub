package models

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"strings"
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

	//12 is the cost of hashing. higher the cost more secure the hash value is but also takes longer
	//to hash the password. So choosing the right balance is important
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users(name, email, hashed_password, created) VALUES (?, ?, ?, UTC_TIMESTAMP())`

	//use the Exec() method to insert the data into database
	_, err = u.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		// If this returns an error, we use the errors.As() function to check
		// whether the error has the type *mysql.MySQLError. If it does, the
		// error will be assigned to the mySQLError variable. We can then check
		// whether or not the error relates to our users_uc_email key by
		// checking if the error code equals 1062 and the contents of the error
		// message string. If it does, we return an ErrDuplicateEmail error.
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (u *UserModel) Authenticate(email, password string) (int, error) {

	var id int
	var hashedPassword []byte

	//fetch the id and password of user from database if the email provided was correct
	stmt := `SELECT id, hashed_password FROM users WHERE email = ?`
	err := u.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)

	//if email was not present in db, then we return an invalid creds error
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	//if email was correct, we compare the hased password with the plain passwd provided by user
	//if they do not match we return a invalid cred error
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))

	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	//if we reached here, means passwd and email were correct. So we return the id and nil as error
	return id, nil
}
func (u *UserModel) Exists(id int) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"
	err := u.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}
