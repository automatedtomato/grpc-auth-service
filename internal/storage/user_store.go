package storage

import (
	"errors"
	"sync"

	"github.com/automatedtomato/grpc-auth-service/internal/model"
)

type UserStore interface {
	Create(user *model.User) error
	GetByUsername(username string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	GetByID(id string) (*model.User, error)
	GetByResetToken(token string) (*model.User, error)
	Update(user *model.User) error
}

type InMemoryUserStore struct {
	users   map[string]*model.User
	byName  map[string]string
	byEmail map[string]string
	byToken map[string]string
	mu      sync.RWMutex
}

func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		users:   make(map[string]*model.User),
		byName:  make(map[string]string),
		byEmail: make(map[string]string),
		byToken: make(map[string]string),
	}
}

func (s *InMemoryUserStore) Create(user *model.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// validate uniqueness of username and email
	if _, exists := s.byName[user.Username]; exists {
		return errors.New("username already exists")
	}

	if _, exists := s.byEmail[user.Email]; exists {
		return errors.New("email already exists")
	}

	// save user
	s.users[user.ID] = user
	s.byName[user.Username] = user.ID
	s.byEmail[user.Email] = user.ID
	return nil
}

func (s *InMemoryUserStore) GetByUsername(username string) (*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	id, exists := s.byName[username]
	if !exists {
		return nil, errors.New("user not found")
	}
	return s.users[id], nil
}

func (s *InMemoryUserStore) GetByEmail(email string) (*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	id, exists := s.byEmail[email]
	if !exists {
		return nil, errors.New("user not found")
	}
	return s.users[id], nil
}

func (s *InMemoryUserStore) GetByID(id string) (*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *InMemoryUserStore) GetByResetToken(token string) (*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	id, exists := s.byToken[token]
	if !exists {
		return nil, errors.New("invalid reset token")
	}
	return s.users[id], nil
}

func (s *InMemoryUserStore) Update(user *model.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[user.ID]; !exists {
		return errors.New("user not found")
	}

	for token, id := range s.byToken {
		if id == user.ID {
			delete(s.byToken, token)
		}
	}
	if user.ResetToken != "" {
		s.byToken[user.ResetToken] = user.ID
	}

	s.users[user.ID] = user
	return nil
}
