package web

import (
	"context"
	"errors"

	catalogmodel "github.com/Maksim-Kot/Tech-store-catalog/pkg/model"
	"github.com/Maksim-Kot/Tech-store-web/internal/gateway"
)

var ErrNotFound = errors.New("not found")

type catalogGateway interface {
	Catalog(ctx context.Context) ([]*catalogmodel.Category, error)
	ProductsByCategoryID(ctx context.Context, id int64) ([]*catalogmodel.Product, error)
	ProductByID(ctx context.Context, id int64) (*catalogmodel.Product, error)
}

type Controller struct {
	catalogGateway catalogGateway
}

func New(catalogGateway catalogGateway) *Controller {
	return &Controller{catalogGateway}
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
