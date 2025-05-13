package memory

import (
	"context"
	"slices"
	"sync"
	"time"

	"github.com/Maksim-Kot/Tech-store-orders/internal/repository"
	"github.com/Maksim-Kot/Tech-store-orders/pkg/model"
)

const (
	StatusNew string = "created"
)

type Repository struct {
	sync.RWMutex
	orders map[int64]*model.Order
}

func New() (*Repository, error) {
	return &Repository{
		orders: map[int64]*model.Order{},
	}, nil
}

func (r *Repository) CreateOrder(_ context.Context, userID int64, price float64, items []model.Item) (int64, error) {
	if len(items) == 0 {
		return 0, repository.ErrNotCreated
	}

	r.Lock()
	defer r.Unlock()

	id := int64(len(r.orders) + 1)

	order := model.Order{
		ID:        id,
		UserID:    userID,
		Price:     price,
		Status:    StatusNew,
		Items:     items,
		CreatedAt: time.Now(),
	}

	r.orders[id] = &order

	return id, nil
}

func (r *Repository) OrderByID(ctx context.Context, id int64) (*model.Order, error) {
	if id < 1 {
		return nil, repository.ErrNotFound
	}

	r.RLock()
	defer r.RUnlock()

	order, exists := r.orders[id]
	if !exists {
		return nil, repository.ErrNotFound
	}

	return order, nil
}

func (r *Repository) OrdersByUserID(ctx context.Context, id int64) ([]*model.Order, error) {
	if id < 1 {
		return nil, repository.ErrNotFound
	}

	r.RLock()
	defer r.RUnlock()

	var orders []*model.Order
	for _, order := range r.orders {
		if order.UserID == id {
			orders = append(orders, order)
		}
	}

	if len(orders) == 0 {
		return nil, repository.ErrNotFound
	}

	slices.SortFunc(orders, func(a, b *model.Order) int {
		return b.CreatedAt.Compare(a.CreatedAt)
	})

	return orders, nil
}
