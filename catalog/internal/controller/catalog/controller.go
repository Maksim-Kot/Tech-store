package catalog

import (
	"context"
	"errors"

	"github.com/Maksim-Kot/Tech-store-catalog/internal/repository"
	"github.com/Maksim-Kot/Tech-store-catalog/pkg/model"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrNotEnough    = errors.New("not enough quantity")
	ErrEditConflict = errors.New("edit conflict")
)

type catalogRepository interface {
	Categories(ctx context.Context) ([]*model.Category, error)
	ProductsByCategoryID(ctx context.Context, id int64) ([]*model.Product, error)
	ProductByID(ctx context.Context, id int64) (*model.Product, error)
	DecreaseProductQuantity(ctx context.Context, id int64, amount int32) error
	IncreaseProductQuantity(ctx context.Context, id int64, amount int32) error
	PutCategory(ctx context.Context, category *model.Category) error
	PutProduct(ctx context.Context, product *model.Product) error
}

type Controller struct {
	repo catalogRepository
}

func New(repo catalogRepository) *Controller {
	return &Controller{repo}
}

func (c *Controller) Categories(ctx context.Context) ([]*model.Category, error) {
	return c.repo.Categories(ctx)
}

func (c *Controller) ProductsByCategoryID(ctx context.Context, id int64) ([]*model.Product, error) {
	products, err := c.repo.ProductsByCategoryID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return products, nil
}

func (c *Controller) ProductByID(ctx context.Context, id int64) (*model.Product, error) {
	product, err := c.repo.ProductByID(ctx, id)

	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return product, nil
}

func (c *Controller) DecreaseProductQuantity(ctx context.Context, id int64, amount int32) error {
	err := c.repo.DecreaseProductQuantity(ctx, id, amount)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			return ErrNotFound
		case errors.Is(err, repository.ErrNotEnough):
			return ErrNotEnough
		case errors.Is(err, ErrEditConflict):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (c *Controller) IncreaseProductQuantity(ctx context.Context, id int64, amount int32) error {
	err := c.repo.IncreaseProductQuantity(ctx, id, amount)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

func (c *Controller) PutCategory(ctx context.Context, category *model.Category) error {
	return c.repo.PutCategory(ctx, category)
}

func (c *Controller) PutProduct(ctx context.Context, product *model.Product) error {
	return c.repo.PutProduct(ctx, product)
}
