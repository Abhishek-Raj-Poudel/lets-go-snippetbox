package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID         int
	Title      string
	Content    string
	Created_At time.Time
	Expires    time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {

	stmt := `INSERT INTO snippets (title,content,created_at,expires)
VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
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
	//This is the long form of the code below
	stmt := `SELECT id, title, content, created_at, expires FROM snippets
	 WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created_At, &s.Expires)

	// s := &Snippet{}
	//
	// err := m.DB.QueryRow("SELECT ...", id).Scan(&s.ID, &s.Title, &s.Content, &s.Created_At, &s.Expires)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT id, title, content, created_at,expires FROM snippets
  WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	row, err := m.DB.Query(stmt)

	if err != nil {
		return nil, err
	}

  defer row.Close()

	snippets := []*Snippet{}

  for row.Next() {
    s := &Snippet{}
    err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created_At, &s.Expires)
    if err != nil {
      return nil,err
    }
    snippets = append(snippets,s)
  }

  if err := row.Err(); err != nil {
    return nil,err
  }

	return snippets, nil
}
