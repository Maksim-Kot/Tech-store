package orders

import (
	"context"
	"errors"

	"github.com/Maksim-Kot/Tech-store-orders/internal/repository"
	"github.com/Maksim-Kot/Tech-store-orders/pkg/model"
)

var (
	ErrNotFound   = errors.New("order not found")
	ErrNotCreated = errors.New("order not created")
)

type ordersRepository interface {
	CreateOrder(ctx context.Context, userID int64, price float64, items []model.Item) (int64, error)
	OrderByID(ctx context.Context, id int64) (*model.Order, error)
	OrdersByUserID(ctx context.Context, id int64) ([]*model.Order, error)
}

type Controller struct {
	repo ordersRepository
}

func New(repo ordersRepository) *Controller {
	return &Controller{repo}
}

func (c *Controller) CreateOrder(ctx context.Context, userID int64, price float64, items []model.Item) (int64, error) {
	id, err := c.repo.CreateOrder(ctx, userID, price, items)
	if err != nil {
		if errors.Is(err, repository.ErrNotCreated) {
			return 0, ErrNotCreated
		}
		return 0, err
	}

	return id, nil
}

func (c *Controller) OrderByID(ctx context.Context, id int64) (*model.Order, error) {
	order, err := c.repo.OrderByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return order, nil
}

func (c *Controller) OrdersByUserID(ctx context.Context, id int64) ([]*model.Order, error) {
	orders, err := c.repo.OrdersByUserID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return orders, nil
}
