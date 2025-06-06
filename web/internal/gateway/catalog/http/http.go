package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Maksim-Kot/Commons/discovery"
	"github.com/Maksim-Kot/Commons/httputil"
	"github.com/Maksim-Kot/Tech-store-catalog/pkg/model"
	"github.com/Maksim-Kot/Tech-store-web/internal/gateway"
)

const (
	serviceName = "catalog"

	baseURL               = "http://%s/v1"
	catalogURL            = baseURL + "/catalog"
	productsByCategoryURL = baseURL + "/category/%d"
	productURL            = baseURL + "/product/%d"
	decreaseProductURL    = baseURL + "/product/%d/decrease/%d"
	increaseProductURL    = baseURL + "/product/%d/increase/%d"
)

type Gateway struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

type categoriesResponse struct {
	Categories []*model.Category `json:"categories"`
}

type productsResponse struct {
	Products []*model.Product `json:"products"`
}

type productResponse struct {
	Product *model.Product `json:"product"`
}

func (g *Gateway) Catalog(ctx context.Context) ([]*model.Category, error) {
	addr, err := httputil.ServiceAddr(ctx, serviceName, g.registry)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf(catalogURL, addr)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	log.Printf("[gateway] GET %s (catalog service)", url)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var wrapper categoriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, err
	}

	return wrapper.Categories, nil
}

func (g *Gateway) ProductsByCategoryID(ctx context.Context, id int64) ([]*model.Product, error) {
	addr, err := httputil.ServiceAddr(ctx, serviceName, g.registry)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf(productsByCategoryURL, addr, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	log.Printf("[gateway] GET %s (catalog service)", url)

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

	var wrapper productsResponse
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, err
	}

	return wrapper.Products, nil
}

func (g *Gateway) ProductByID(ctx context.Context, id int64) (*model.Product, error) {
	addr, err := httputil.ServiceAddr(ctx, serviceName, g.registry)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf(productURL, addr, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	log.Printf("[gateway] GET %s (catalog service)", url)

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

	var wrapper productResponse
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, err
	}

	return wrapper.Product, nil
}

func (g *Gateway) DecreaseProductQuantity(ctx context.Context, id int64, amount int32) error {
	addr, err := httputil.ServiceAddr(ctx, serviceName, g.registry)
	if err != nil {
		return err
	}
	url := fmt.Sprintf(decreaseProductURL, addr, id, amount)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	log.Printf("[gateway] POST %s (catalog service)", url)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return gateway.ErrNotFound
		case http.StatusBadRequest:
			return gateway.ErrNotEnough
		case http.StatusConflict:
			return gateway.ErrEditConflict
		default:
			return fmt.Errorf("unexpected status: %s", resp.Status)
		}
	}

	return nil
}

func (g *Gateway) IncreaseProductQuantity(ctx context.Context, id int64, amount int32) error {
	addr, err := httputil.ServiceAddr(ctx, serviceName, g.registry)
	if err != nil {
		return err
	}
	url := fmt.Sprintf(increaseProductURL, addr, id, amount)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	log.Printf("[gateway] POST %s (catalog service)", url)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return gateway.ErrNotFound
		default:
			return fmt.Errorf("unexpected status: %s", resp.Status)
		}
	}

	return nil
}
