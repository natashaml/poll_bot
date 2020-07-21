package main

type Storage struct {
	users map[string]*StorageUser
}

func newStorage() *Storage {
	return &Storage{
		users: map[string]*StorageUser{},
	}
}

type StorageUser struct {
	Id    string
	Level int
}

func (s *Storage) Obtain(id string) (*StorageUser, error) {
	user := s.users[id]
	if user != nil {
		return user, nil
	}
	user = &StorageUser{Id: id}
	s.users[id] = user

	return user, nil
}
