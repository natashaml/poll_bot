package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

type Storage struct {
	users      sync.Map
	persistent *persistenseStorage
}

func newStorage() (*Storage, error) {
	persistent, err := newPersistenseStoragePq()
	if err != nil {
		return nil, err
	}

	return &Storage{
		users:      sync.Map{},
		persistent: persistent,
	}, nil
}

type StorageUser struct {
	Id                  string
	Level               int
	ConversationStarted bool
	Age                 string
	WillTakePart        string
	WhenVote            string
	Gender              string
	Education           string

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
	if s.persistent == nil {
		return nil, errors.New("persistence not enabled")
	}
	user, err := s.persistent.load(id)
	if err != nil {
		log.Printf("Unable to load persistent user %v", err)
		return nil, nil
	}
	return user, nil
}

func (s *Storage) Clear(id string) error {
	s.users.Delete(id)

	err := s.persistent.clear(id)
	if err != nil {
		log.Printf("Unable to clear persistent user %v", err)
		return nil
	}

	return nil
}

func (s *Storage) PersistCount() (int, error) {
	if s.persistent == nil {
		return 0, errors.New("persistence not enabled")
	}

	count, err := s.persistent.count()
	if err != nil {
		log.Printf("Unable to count persistent user %v", err)
		return 0, nil
	}

	return count, nil
}

func (s *Storage) Persist(id string) error {
	if s.persistent == nil {
		return errors.New("persistence not enabled")
	}

	user := s.fromCache(id)
	if user == nil {
		return fmt.Errorf("%v missed in cache", id)
	}

	err := s.persistent.save(user)
	if err != nil {
		log.Printf("Unable to save persistent user %v", err)
		return nil
	}

	return nil
}
