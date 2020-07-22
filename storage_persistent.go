package main

import (
	"database/sql"
	"errors"
	"os"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type persistenseStorage struct {
	db *sql.DB
}

func newPersistenseStorageSqllite() (*persistenseStorage, error) {
	db, err := sql.Open("sqlite3", "file:db.sqlite?cache=shared")
	if err != nil {
		return nil, err
	}

	batch := []string{
		`CREATE TABLE users (id TEXT PRIMARY KEY, conversation_started BOOLEAN, level INT, age INT, candidate TEXT);`,
	}

	for _, b := range batch {
		_, _ = db.Exec(b)
	}

	return &persistenseStorage{
		db: db,
	}, nil
}

func newPersistenseStoragePq() (*persistenseStorage, error) {
	connStr := os.Getenv("DB_CONNECTION")
	if connStr == "" {
		return nil, errors.New("DB_CONNECTION is empty")
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(2)

	return &persistenseStorage{
		db: db,
	}, nil
}

func (s *persistenseStorage) load(id string) (*StorageUser, error) {
	sqlStatement := `SELECT id, level, conversation_started, age, candidate FROM users WHERE id = $1;`
	var user StorageUser
	row := s.db.QueryRow(sqlStatement, id)
	err := row.Scan(&user.Id, &user.Level, &user.ConversationStarted, &user.Age, &user.Candidate)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *persistenseStorage) clear(id string) error {
	sqlStatement := `DELETE FROM users WHERE id = $1;`
	_, err := s.db.Exec(sqlStatement, id)
	return err
}

func (s *persistenseStorage) count() (int, error) {
	sqlStatement := `SELECT COUNT(*) FROM users;`
	row := s.db.QueryRow(sqlStatement)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *persistenseStorage) save(user *StorageUser) error {
	sqlStatement := `SELECT COUNT(*) FROM users WHERE id = $1;`
	row := s.db.QueryRow(sqlStatement, user.Id)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		sqlStatement = `INSERT INTO users (id, level, conversation_started, age, candidate) VALUES ($1, $2, $3, $4, $5)`
		_, err = s.db.Exec(sqlStatement, user.Id, user.Level, user.ConversationStarted, user.Age, user.Candidate)
		return err
	}
	sqlStatement = `UPDATE users SET level=$2, conversation_started=$3, age=$4, candidate=$5 WHERE id = $1`
	_, err = s.db.Exec(sqlStatement, user.Id, user.Level, user.ConversationStarted, user.Age, user.Candidate)
	return err
}
