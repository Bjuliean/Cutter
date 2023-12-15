package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func (s *Storage)AAA() {
	fmt.Println(123)
}

func New(storagePath string) (*Storage, error) {
	const ferr = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ferr, err)
	}

	st, err := db.Prepare("CREATE TABLE IF NOT EXISTS url(id INTEGER PRIMARY KEY, alias TEXT NOT NULL UNIQUE, url TEXT NOT NULL); CREATE INDEX IF NOT EXIST idx_alias ON url(alias);")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ferr, err)
	}
	_, err = st.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ferr, err)
	}

	return &Storage{db: db}, nil
}