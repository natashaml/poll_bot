package main

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
)

type Storage struct {
	users sync.Map
	db    *sql.DB
}

func newStorage() *Storage {
	return &Storage{
		users: sync.Map{},
	}
}

type StorageUser struct {
	Id    string
	Level int
	Age   int

	Candidate string
}

func (s *Storage) Obtain(id string) (*StorageUser, error) {
	user, err := s.fromCache(id)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}
	persistedUser, err := s.fromPersisted(id)
	if err != nil {
		return nil, err
	}
	if persistedUser != nil {
		return persistedUser, nil
	}

	newUser := &StorageUser{Id: id}
	s.users.Store(id, newUser)

	return newUser, nil
}

func (s *Storage) fromCache(id string) (*StorageUser, error) {
	user, ok := s.users.Load(id)
	if ok && user != nil {
		return user.(*StorageUser), nil
	}
	return nil, nil
}

func (s *Storage) fromPersisted(id string) (*StorageUser, error) {
	if s.db == nil {
		return nil, errors.New("persistence not enabled")
	}

	sqlStatement := `SELECT id, level, age, candidate FROM users WHERE id = $1;`
	var user StorageUser
	row := s.db.QueryRow(sqlStatement, id)
	err := row.Scan(&user.Id, &user.Level, &user.Age, &user.Candidate)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
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

	user, err := s.fromCache(id)
	if user == nil {
		return fmt.Errorf("%v missed in cache", id)
	}

	sqlStatement := `SELECT COUNT(*) FROM users WHERE id = $1;`
	row := s.db.QueryRow(sqlStatement, user.Id)
	var count int
	err = row.Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		sqlStatement = `INSERT INTO users (id, level, age, candidate) VALUES ($1, $2, $3, $4)`
		_, err = s.db.Exec(sqlStatement, user.Id, user.Level, user.Age, user.Candidate)
		return err
	}
	sqlStatement = `UPDATE users SET level=$2, age=$3, candidate=$4 WHERE id = $1`
	_, err = s.db.Exec(sqlStatement, user.Id, user.Level, user.Age, user.Candidate)
	return err
}
