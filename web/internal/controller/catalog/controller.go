package catalog

import (
	"context"
	"errors"

	"github.com/Maksim-Kot/Tech-store-catalog/pkg/model"
	"github.com/Maksim-Kot/Tech-store-web/internal/controller"
	"github.com/Maksim-Kot/Tech-store-web/internal/gateway"
)

type catalogGateway interface {
	Catalog(ctx context.Context) ([]*model.Category, error)
	ProductsByCategoryID(ctx context.Context, id int64) ([]*model.Product, error)
	ProductByID(ctx context.Context, id int64) (*model.Product, error)
}

type CatalogController struct {
	catalogGateway catalogGateway
}

func New(catalogGateway catalogGateway) *CatalogController {
	return &CatalogController{catalogGateway: catalogGateway}
}

func (c *CatalogController) Catalog(ctx context.Context) ([]*model.Category, error) {
	categories, err := c.catalogGateway.Catalog(ctx)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (c *CatalogController) ProductsByCategoryID(ctx context.Context, id int64) ([]*model.Product, error) {
	products, err := c.catalogGateway.ProductsByCategoryID(ctx, id)

	if err != nil {
		if errors.Is(err, gateway.ErrNotFound) {
			return nil, controller.ErrNotFound
		}
		return nil, err
	}

	return products, nil
}

func (c *CatalogController) ProductByID(ctx context.Context, id int64) (*model.Product, error) {
	product, err := c.catalogGateway.ProductByID(ctx, id)

	if err != nil {
		if errors.Is(err, gateway.ErrNotFound) {
			return nil, controller.ErrNotFound
		}
		return nil, err
	}

	return product, nil
}
