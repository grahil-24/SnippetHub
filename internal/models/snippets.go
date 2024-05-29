package models

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Snippet struct {
	ID      int
	Content string
	Title   string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(Content string, Title string, Expires int) (int, error) {
	//? denotes the placeholder where we insert our data. This helps to prevent SQL injection attacks
	stmt := "INSERT INTO snippets (title, content, created, expires) VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))"

	//better to use transactions when dealing with database to preserve data consistency
	tx, err := m.DB.Begin()
	if err != nil {
		return 0, err
	}

	//if tx operations return any error the transaction is rolled back
	defer tx.Rollback()

	//exec is used for operations which do not return a row like delete and create
	result, err := tx.Exec(stmt, Content, Title, Expires)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT  id, title, content, created, expires FROM snippets WHERE id = ? AND expires > UTC_TIMESTAMP()`
	s := &Snippet{}
	row := m.DB.QueryRow(stmt, id)
	//use Scan to get the returned values
	//char, varchar string maps to string, int to int, BIGINT to int64 and Boolean to bool
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		//return ErrNoRows error on not finding any rows
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	//return s if all went ok
	return s, nil
}

// get the latest 10 snippets
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets ORDER BY created DESC LIMIT 10`
	snippets := []*Snippet{}
	rows, err := m.DB.Query(stmt)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		s := &Snippet{}
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
