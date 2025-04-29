package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Maksim-Kot/Tech-store-catalog/config"
	"github.com/Maksim-Kot/Tech-store-catalog/internal/controller/catalog"
	"github.com/Maksim-Kot/Tech-store-catalog/pkg/model"
)

type Handler struct {
	ctrl *catalog.Controller
	cfg  config.APIConfig
}

func New(ctrl *catalog.Controller, cfg config.APIConfig) *Handler {
	return &Handler{ctrl, cfg}
}

func (h *Handler) HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": h.cfg.Env,
			"version":     h.cfg.Version,
			"name":        h.cfg.Name,
		},
	}

	err := h.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

func (h *Handler) CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	categories, err := h.ctrl.Categories(ctx)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	err = h.writeJSON(w, http.StatusOK, envelope{"categories": categories}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

func (h *Handler) ProductsByCategoryIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := h.getID(r)
	if err != nil || id < 1 {
		h.notFoundResponse(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	products, err := h.ctrl.ProductsByCategoryID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, catalog.ErrNotFound):
			h.notFoundResponse(w, r)
		default:
			h.serverErrorResponse(w, r, err)
		}
		return
	}

	err = h.writeJSON(w, http.StatusOK, envelope{"products": products}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

func (h *Handler) ProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := h.getID(r)
	if err != nil || id < 1 {
		h.notFoundResponse(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	product, err := h.ctrl.ProductByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, catalog.ErrNotFound):
			h.notFoundResponse(w, r)
		default:
			h.serverErrorResponse(w, r, err)
		}
		return
	}

	err = h.writeJSON(w, http.StatusOK, envelope{"product": product}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

func (h *Handler) DecreaseProductQuantityHandler(w http.ResponseWriter, r *http.Request) {
	id, err := h.getID(r)
	if err != nil || id < 1 {
		h.notFoundResponse(w, r)
		return
	}

	amount, err := h.getAmount(r)
	if err != nil || amount < 1 {
		h.badRequestResponse(w, r, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = h.ctrl.DecreaseProductQuantity(ctx, id, amount)
	if err != nil {
		switch {
		case errors.Is(err, catalog.ErrNotFound):
			h.notFoundResponse(w, r)
		case errors.Is(err, catalog.ErrNotEnough):
			h.badRequestResponse(w, r, err)
		case errors.Is(err, catalog.ErrEditConflict):
			h.editConflictResponse(w, r)
		default:
			h.serverErrorResponse(w, r, err)
		}
		return
	}
}

func (h *Handler) IncreaseProductQuantityHandler(w http.ResponseWriter, r *http.Request) {
	id, err := h.getID(r)
	if err != nil || id < 1 {
		h.notFoundResponse(w, r)
		return
	}

	amount, err := h.getAmount(r)
	if err != nil || amount < 1 {
		h.badRequestResponse(w, r, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = h.ctrl.IncreaseProductQuantity(ctx, id, amount)
	if err != nil {
		switch {
		case errors.Is(err, catalog.ErrNotFound):
			h.notFoundResponse(w, r)
		default:
			h.serverErrorResponse(w, r, err)
		}
		return
	}
}

func (h *Handler) PutCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
	}

	err := h.readJSON(w, r, &input)
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	category := &model.Category{
		Name: input.Name,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = h.ctrl.PutCategory(ctx, category)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	err = h.writeJSON(w, http.StatusCreated, envelope{"category": category}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

func (h *Handler) PutProductHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string          `json:"name"`
		Description string          `json:"description,omitempty"`
		Price       float64         `json:"price"`
		Quantity    int32           `json:"quantity"`
		ImageURL    string          `json:"image_url,omitempty"`
		Attributes  json.RawMessage `json:"attributes"`
		CategoryID  int64           `json:"category_id"`
	}

	err := h.readJSON(w, r, &input)
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	product := &model.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Quantity:    input.Quantity,
		ImageURL:    input.ImageURL,
		Attributes:  input.Attributes,
		CategoryID:  input.CategoryID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = h.ctrl.PutProduct(ctx, product)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	err = h.writeJSON(w, http.StatusCreated, envelope{"product": product}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}
