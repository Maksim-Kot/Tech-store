package orders

import (
	"context"
	"errors"

	ordersmodel "github.com/Maksim-Kot/Tech-store-orders/pkg/model"
	"github.com/Maksim-Kot/Tech-store-web/internal/controller"
	"github.com/Maksim-Kot/Tech-store-web/internal/gateway"
)

type ordersGateway interface {
	OrderByID(ctx context.Context, id int64) (*ordersmodel.Order, error)
	OrdersByUserID(ctx context.Context, id int64) ([]*ordersmodel.Order, error)
}

type OrdersController struct {
	ordersGateway ordersGateway
}

func New(ordersGateway ordersGateway) *OrdersController {
	return &OrdersController{ordersGateway: ordersGateway}
}

func (c *OrdersController) OrderByID(ctx context.Context, id int64) (*ordersmodel.Order, error) {
	order, err := c.ordersGateway.OrderByID(ctx, id)

	if err != nil {
		if errors.Is(err, gateway.ErrNotFound) {
			return nil, controller.ErrNotFound
		}
		return nil, err
	}

	return order, nil
}

func (c *OrdersController) OrdersByUserID(ctx context.Context, id int64) ([]*ordersmodel.Order, error) {
	orders, err := c.ordersGateway.OrdersByUserID(ctx, id)

	if err != nil {
		if errors.Is(err, gateway.ErrNotFound) {
			return nil, controller.ErrNotFound
		}
		return nil, err
	}

	return orders, nil
}
