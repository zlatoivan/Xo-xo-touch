package user

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrNoUser  = errors.New("no user found")
	ErrBadPass = errors.New("invalid password")
)

type UserMemoryRepository struct {
	data  map[string]*User
	mutex *sync.Mutex
}

func NewMemoryRepo() *UserMemoryRepository {
	return &UserMemoryRepository{
		data:  map[string]*User{},
		mutex: &sync.Mutex{},
	}
}

func (repo *UserMemoryRepository) Register(username, password string) (*User, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	newID := uuid.New()
	user := &User{
		ID:       newID.String(),
		Username: username,
		Password: password,
	}
	repo.data[user.ID] = user

	return user, nil
}

func (repo *UserMemoryRepository) Authorize(username, password string) (*User, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	user, ok := repo.data[username]
	if !ok {
		return nil, ErrNoUser
	}

	if user.Password != password {
		return nil, ErrBadPass
	}

	return user, nil
}
