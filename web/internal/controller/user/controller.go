package user

import (
	"context"
	"errors"

	"github.com/Maksim-Kot/Tech-store-web/internal/controller"
	"github.com/Maksim-Kot/Tech-store-web/internal/repository"
)

type userRepo interface {
	Insert(ctx context.Context, name, email, password string) error
	Authenticate(ctx context.Context, email, password string) (int64, error)
	Exists(ctx context.Context, id int64) (bool, error)
}

type UserController struct {
	userRepo userRepo
}

func New(userRepo userRepo) *UserController {
	return &UserController{userRepo: userRepo}
}

func (c *UserController) InsertUser(ctx context.Context, name, email, password string) error {
	err := c.userRepo.Insert(ctx, name, email, password)

	if err != nil {
		if errors.Is(err, repository.ErrDuplicateEmail) {
			return controller.ErrDuplicateEmail
		}
		return err
	}

	return nil
}

func (c *UserController) AuthenticateUser(ctx context.Context, email, password string) (int64, error) {
	id, err := c.userRepo.Authenticate(ctx, email, password)

	if err != nil {
		if errors.Is(err, repository.ErrInvalidCredentials) {
			return 0, controller.ErrInvalidCredentials
		}
		return 0, err
	}

	return id, nil
}

func (c *UserController) UserExists(ctx context.Context, id int64) (bool, error) {
	return c.userRepo.Exists(ctx, id)
}
