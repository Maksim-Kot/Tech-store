package web

import (
	"context"
	"errors"

	catalogmodel "github.com/Maksim-Kot/Tech-store-catalog/pkg/model"
	"github.com/Maksim-Kot/Tech-store-web/internal/gateway"
	"github.com/Maksim-Kot/Tech-store-web/internal/repository"
)

var (
	ErrNotFound           = errors.New("not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrDuplicateEmail     = errors.New("duplicate email")
)

type catalogGateway interface {
	Catalog(ctx context.Context) ([]*catalogmodel.Category, error)
	ProductsByCategoryID(ctx context.Context, id int64) ([]*catalogmodel.Product, error)
	ProductByID(ctx context.Context, id int64) (*catalogmodel.Product, error)
}

type userRepo interface {
	Insert(ctx context.Context, name, email, password string) error
	Authenticate(ctx context.Context, email, password string) (int64, error)
	Exists(ctx context.Context, id int64) (bool, error)
}

type Controller struct {
	catalogGateway catalogGateway
	userRepo       userRepo
}

func New(catalogGateway catalogGateway, userRepo userRepo) *Controller {
	return &Controller{catalogGateway, userRepo}
}

func (c *Controller) Catalog(ctx context.Context) ([]*catalogmodel.Category, error) {
	categories, err := c.catalogGateway.Catalog(ctx)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (c *Controller) ProductsByCategoryID(ctx context.Context, id int64) ([]*catalogmodel.Product, error) {
	products, err := c.catalogGateway.ProductsByCategoryID(ctx, id)

	if err != nil {
		if errors.Is(err, gateway.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return products, nil
}

func (c *Controller) ProductByID(ctx context.Context, id int64) (*catalogmodel.Product, error) {
	product, err := c.catalogGateway.ProductByID(ctx, id)

	if err != nil {
		if errors.Is(err, gateway.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return product, nil
}

func (c *Controller) InsertUser(ctx context.Context, name, email, password string) error {
	err := c.userRepo.Insert(ctx, name, email, password)

	if err != nil {
		if errors.Is(err, repository.ErrDuplicateEmail) {
			return ErrDuplicateEmail
		}
		return err
	}

	return nil
}

func (c *Controller) AuthenticateUser(ctx context.Context, email, password string) (int64, error) {
	id, err := c.userRepo.Authenticate(ctx, email, password)

	if err != nil {
		if errors.Is(err, repository.ErrInvalidCredentials) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	return id, nil
}

func (c *Controller) UserExists(ctx context.Context, id int64) (bool, error) {
	return c.userRepo.Exists(ctx, id)
}
