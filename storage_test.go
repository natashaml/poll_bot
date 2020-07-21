package main

import (
	"database/sql"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	_ "github.com/proullon/ramsql/driver"
)

func randomString() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZÅÄÖ" +
		"abcdefghijklmnopqrstuvwxyzåäö" +
		"0123456789")
	length := 8
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String() // E.g. "ExcbsVQs"
	return str
}

func newTestStorage() (*Storage, error) {

	db, err := sql.Open("ramsql", randomString())
	if err != nil {
		return nil, err
	}

	return &Storage{
		users: sync.Map{},
		db:    db,
	}, nil
}

func (s *Storage) init() error {
	batch := []string{
		`CREATE TABLE users (id TEXT PRIMARY KEY, level INT, age INT, candidate TEXT);`,
	}

	for _, b := range batch {
		_, err := s.db.Exec(b)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestMapStorage(t *testing.T) {
	s, err := newTestStorage()
	require.NoError(t, err)
	err = s.init()
	require.NoError(t, err)

	user, err := s.Obtain("12")
	require.NoError(t, err)
	require.Equal(t, user.Id, "12")
	require.Equal(t, user.Age, 0)
	user.Age = 16

	user, err = s.Obtain("12")
	require.NoError(t, err)
	require.Equal(t, user.Id, "12")
	require.Equal(t, user.Age, 16)

	err = s.Persist("12")
	require.NoError(t, err)

	count, err := s.PersistCount()
	require.NoError(t, err)
	require.Equal(t, count, 1)

	err = s.Persist("12")
	require.NoError(t, err)

	count, err = s.PersistCount()
	require.NoError(t, err)
	require.Equal(t, count, 1)
}
