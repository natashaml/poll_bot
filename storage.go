package main

import (
	"sync"
)

type Storage struct {
	users sync.Map
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
	user, ok := s.users.Load(id)
	if ok && user != nil {
		return user.(*StorageUser), nil
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

func (s *Storage) fromPersisted(id string) (*StorageUser, error) {
	return nil, nil
}

func (s *Storage) Persist(id string) error {
	return nil
}
