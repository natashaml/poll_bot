package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

type Storage struct {
	users sync.Map
	db    *sql.DB
}

func newStorage() (*Storage, error) {
	connStr := os.Getenv("DB_CONNECTION")
	if connStr == "" {
		return nil, errors.New("DB_CONNECTION is empty")
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(2)

	return &Storage{
		users: sync.Map{},
		db:    db,
	}, nil
}

type StorageUser struct {
	Id                  string
	Level               int
	ConversationStarted bool
	Age                 int

	Candidate string
}

func (u *StorageUser) validate() error {
	if u.Id == "" {
		return errors.New("Empty user id")
	}
	return nil
}

func (s *Storage) Obtain(id string) (*StorageUser, error) {
	if id == "" {
		return nil, errors.New("Unable to obtain empty id")
	}

	user := s.fromCache(id)
	if user != nil {
		return user, user.validate()
	}
	persistedUser, err := s.fromPersisted(id)
	if err != nil {
		return nil, err
	}
	if persistedUser != nil {
		return persistedUser, persistedUser.validate()
	}

	newUser := &StorageUser{Id: id}
	s.users.Store(id, newUser)

	return newUser, newUser.validate()
}

func (s *Storage) fromCache(id string) *StorageUser {
	user, ok := s.users.Load(id)
	if ok && user != nil {
		return user.(*StorageUser)
	}
	return nil
}

func (s *Storage) fromPersisted(id string) (*StorageUser, error) {
	if s.db == nil {
		return nil, errors.New("persistence not enabled")
	}

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

func (s *Storage) Clear(id string) error {
	s.users.Delete(id)

	sqlStatement := `DELETE FROM users WHERE id = $1;`
	_, err := s.db.Exec(sqlStatement, id)
	return err
}

func (s *Storage) PersistCount() (int, error) {
	if s.db == nil {
		return 0, errors.New("persistence not enabled")
	}

	sqlStatement := `SELECT COUNT(*) FROM users;`
	row := s.db.QueryRow(sqlStatement)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Storage) Persist(id string) error {
	if s.db == nil {
		return errors.New("persistence not enabled")
	}

	user := s.fromCache(id)
	if user == nil {
		return fmt.Errorf("%v missed in cache", id)
	}

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
