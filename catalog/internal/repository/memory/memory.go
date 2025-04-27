package memory

import (
	"cmp"
	"context"
	"slices"
	"sync"

	"github.com/Maksim-Kot/Tech-store-catalog/internal/repository"
	"github.com/Maksim-Kot/Tech-store-catalog/pkg/model"
)

type Repository struct {
	sync.RWMutex
	categories map[int64]*model.Category
	products   map[int64]*model.Product
}

func New() (*Repository, error) {
	return &Repository{
		categories: map[int64]*model.Category{},
		products:   map[int64]*model.Product{},
	}, nil
}

func (r *Repository) Categories(_ context.Context) ([]*model.Category, error) {
	r.RLock()
	defer r.RUnlock()

	var categories []*model.Category
	for _, c := range r.categories {
		categories = append(categories, c)
	}
	if len(categories) == 0 {
		return nil, repository.ErrNotFound
	}

	slices.SortFunc(categories, func(a, b *model.Category) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return categories, nil
}

func (r *Repository) ProductsByCategoryID(_ context.Context, id int64) ([]*model.Product, error) {
	r.RLock()
	defer r.RUnlock()

	var products []*model.Product
	for _, p := range r.products {
		if p.CategoryID == id {
			products = append(products, p)
		}
	}

	if len(products) == 0 {
		return nil, repository.ErrNotFound
	}

	slices.SortFunc(products, func(a, b *model.Product) int {
		return cmp.Compare(a.ID, b.ID)
	})

	return products, nil
}

func (r *Repository) ProductByID(_ context.Context, id int64) (*model.Product, error) {
	r.RLock()
	defer r.RUnlock()

	product, ok := r.products[id]
	if !ok {
		return nil, repository.ErrNotFound
	}

	return product, nil
}

func (r *Repository) DecreaseProductQuantity(_ context.Context, id int64, amount int32) error {
	r.Lock()
	defer r.Unlock()

	product, exists := r.products[id]
	if !exists {
		return repository.ErrNotFound
	}

	if product.Quantity < amount {
		return repository.ErrNotEnough
	}

	product.Quantity -= amount
	return nil
}

func (r *Repository) IncreaseProductQuantity(_ context.Context, id int64, amount int32) error {
	r.Lock()
	defer r.Unlock()

	product, exists := r.products[id]
	if !exists {
		return repository.ErrNotFound
	}

	product.Quantity += amount
	return nil
}

func (r *Repository) PutCategory(ctx context.Context, category *model.Category) error {
	r.Lock()
	defer r.Unlock()

	id := int64(len(r.categories) + 1)
	category.ID = id

	r.categories[id] = category

	return nil
}

func (r *Repository) PutProduct(ctx context.Context, product *model.Product) error {
	r.Lock()
	defer r.Unlock()

	id := int64(len(r.products) + 1)
	product.ID = id

	r.products[id] = product

	return nil
}
