package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Maksim-Kot/Commons/discovery"
	"github.com/Maksim-Kot/Commons/httputil"
	"github.com/Maksim-Kot/Tech-store-orders/pkg/model"
	"github.com/Maksim-Kot/Tech-store-web/internal/gateway"
)

const (
	serviceName = "orders"

	baseURL          = "http://%s/v1"
	createOrderURL   = baseURL + "/order"
	orderByIdURL     = baseURL + "/order/%d"
	orderByUserIdURL = baseURL + "/orders/user/%d"
)

type Gateway struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

type orderResponse struct {
	Order *model.Order `json:"order"`
}

type ordersResponse struct {
	Orders []*model.Order `json:"orders"`
}

type createOrderResponse struct {
	Order struct {
		ID int64 `json:"id"`
	} `json:"order"`
}

func (g *Gateway) OrderByID(ctx context.Context, id int64) (*model.Order, error) {
	addr, err := httputil.ServiceAddr(ctx, serviceName, g.registry)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf(orderByIdURL, addr, id)

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
	addr, err := httputil.ServiceAddr(ctx, serviceName, g.registry)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf(orderByUserIdURL, addr, id)

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

func (g *Gateway) CreateOrder(ctx context.Context, userID int64, price float64, items []*model.Item) (int64, error) {
	addr, err := httputil.ServiceAddr(ctx, serviceName, g.registry)
	if err != nil {
		return 0, err
	}
	url := fmt.Sprintf(createOrderURL, addr)

	orderReq := struct {
		UserID int64         `json:"user_id"`
		Price  float64       `json:"price"`
		Items  []*model.Item `json:"items"`
	}{
		UserID: userID,
		Price:  price,
		Items:  items,
	}

	body, err := json.Marshal(orderReq)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	log.Printf("[gateway] POST %s (orders service)", url)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return 0, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var wrapper createOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return 0, err
	}

	return wrapper.Order.ID, nil
}
