package http

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Maksim-Kot/Tech-store-orders/config"
	"github.com/Maksim-Kot/Tech-store-orders/internal/controller/orders"
	"github.com/Maksim-Kot/Tech-store-orders/pkg/model"
)

type Handler struct {
	ctrl *orders.Controller
	cfg  config.APIConfig
}

func New(ctrl *orders.Controller, cfg config.APIConfig) *Handler {
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
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID int64        `json:"user_id"`
		Price  float64      `json:"price"`
		Items  []model.Item `json:"items"`
	}

	err := h.readJSON(w, r, &input)
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	id, err := h.ctrl.CreateOrder(ctx, input.UserID, input.Price, input.Items)
	if err != nil {
		switch {
		case errors.Is(err, orders.ErrNotCreated):
			h.badRequestResponse(w, r, err)
		default:
			h.ServerErrorResponse(w, r, err)
		}
		return
	}

	message := map[string]int64{
		"id": id,
	}

	err = h.writeJSON(w, http.StatusCreated, envelope{"order": message}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) OrderByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := h.getID(r)
	if err != nil || id < 1 {
		h.notFoundResponse(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	order, err := h.ctrl.OrderByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, orders.ErrNotFound):
			h.notFoundResponse(w, r)
		default:
			h.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = h.writeJSON(w, http.StatusOK, envelope{"order": order}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) OrdersByUserIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := h.getID(r)
	if err != nil || id < 1 {
		h.notFoundResponse(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	purchases, err := h.ctrl.OrdersByUserID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, orders.ErrNotFound):
			h.notFoundResponse(w, r)
		default:
			h.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = h.writeJSON(w, http.StatusOK, envelope{"orders": purchases}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}
