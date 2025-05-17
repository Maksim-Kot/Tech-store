package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Maksim-Kot/Tech-store-orders/pkg/model"
	"github.com/Maksim-Kot/Tech-store-web/internal/gateway"
)

const (
	baseURL          = "http://%s/v1"
	createOrderURL   = baseURL + "/order"
	orderByIdURL     = baseURL + "/order/%d"
	orderByUserIdURL = baseURL + "/orders/user/%d"
)

type Gateway struct {
	addr string
}

func New(addr string) *Gateway {
	return &Gateway{addr}
}

type orderResponse struct {
	Order *model.Order `json:"order"`
}

type ordersResponse struct {
	Orders []*model.Order `json:"orders"`
}

func (g *Gateway) OrderByID(ctx context.Context, id int64) (*model.Order, error) {
	url := fmt.Sprintf(orderByIdURL, g.addr, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	log.Printf("[gateway] GET %s (orders service)", url)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return nil, gateway.ErrNotFound
		default:
			return nil, fmt.Errorf("unexpected status: %s", resp.Status)
		}
	}

	var wrapper orderResponse
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, err
	}

	return wrapper.Order, nil
}

func (g *Gateway) OrdersByUserID(ctx context.Context, id int64) ([]*model.Order, error) {
	url := fmt.Sprintf(orderByUserIdURL, g.addr, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	log.Printf("[gateway] GET %s (orders service)", url)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return nil, gateway.ErrNotFound
		default:
			return nil, fmt.Errorf("unexpected status: %s", resp.Status)
		}
	}

	var wrapper ordersResponse
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, err
	}

	return wrapper.Orders, nil
}
