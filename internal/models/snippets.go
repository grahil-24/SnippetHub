package models

import (
	"database/sql"
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
	stmt := "Insert into snippets(content, title,  created, expires) values(?, ?, UTC_TIMESTAMP(), DATE ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))"

	//exec is used for operations which do not return a row like delete and create
	result, err := m.DB.Exec(stmt, Content, Title, Expires)
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
	return nil, nil
}

// get the latest 10 snippets
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
