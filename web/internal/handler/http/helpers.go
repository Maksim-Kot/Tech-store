package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Maksim-Kot/Tech-store-catalog/pkg/model"
)

func (h *Handler) getID(r *http.Request) (int64, error) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

// processedProduct represents a product with parsed and normalized attributes
// for safe and readable rendering in the UI.
type processedProduct struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Price       float64        `json:"price"`
	Quantity    int32          `json:"quantity"`
	ImageURL    string         `json:"image_url,omitempty"`
	Attributes  map[string]any `json:"attributes"`
	CategoryID  int64          `json:"category_id"`
}

// transformProductForView converts the original Product structure, received
// from the catalog service, into a format suitable for UI rendering.
func transformProductAttributes(product *model.Product) (*processedProduct, error) {
	var attributes map[string]any
	if err := json.Unmarshal(product.Attributes, &attributes); err != nil {
		return nil, fmt.Errorf("invalid attributes format: %w", err)
	}

	processedAttributes := make(map[string]any)
	for key, value := range attributes {
		newKey := strings.ReplaceAll(key, "_", " ")
		processedAttributes[newKey] = value
	}

	return &processedProduct{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Quantity:    product.Quantity,
		ImageURL:    product.ImageURL,
		Attributes:  processedAttributes,
		CategoryID:  product.CategoryID,
	}, nil
}
