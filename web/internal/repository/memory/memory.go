package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/Maksim-Kot/Tech-store-web/internal/model"
	"github.com/Maksim-Kot/Tech-store-web/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type Repository struct {
	sync.RWMutex
	// Contains email -> user instance
	users map[string]*model.User
}

func New() (*Repository, error) {
	return &Repository{
		users: map[string]*model.User{},
	}, nil
}

func (r *Repository) Insert(_ context.Context, name, email, password string) error {
	r.Lock()
	defer r.Unlock()

	if _, exists := r.users[email]; exists {
		return repository.ErrDuplicateEmail
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	id := int64(len(r.users) + 1)
	user := &model.User{
		ID:             id,
		Name:           name,
		Email:          email,
		HashedPassword: hashedPassword,
		Created:        time.Now(),
	}

	r.users[email] = user

	return nil
}

func (r *Repository) Authenticate(_ context.Context, email, password string) (int64, error) {
	r.RLock()
	defer r.RUnlock()

	if _, exists := r.users[email]; !exists {
		return 0, repository.ErrInvalidCredentials
	}

	user := r.users[email]

	id := user.ID
	hashedPassword := user.HashedPassword

	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, repository.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (r *Repository) Exists(_ context.Context, id int64) (bool, error) {
	r.RLock()
	defer r.RUnlock()

	for _, u := range r.users {
		if u.ID == id {
			return true, nil
		}
	}

	return false, nil
}
