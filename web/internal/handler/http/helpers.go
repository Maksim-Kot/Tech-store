package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Maksim-Kot/Tech-store-catalog/pkg/model"
	"github.com/Maksim-Kot/Tech-store-web/internal/contexkeys"

	"github.com/go-playground/form/v4"
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

func (h *Handler) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := h.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		h.ServerError(w, err)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}

func (h *Handler) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           h.SessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: h.IsAuthenticated(r),
	}
}

func (h *Handler) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = h.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		return err
	}

	return nil
}

func (h *Handler) IsAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(contexkeys.IsAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
