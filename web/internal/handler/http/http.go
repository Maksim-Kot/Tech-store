package http

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/Maksim-Kot/Tech-store-web/internal/controller/web"
)

type Handler struct {
	ctrl          *web.Controller
	templateCache map[string]*template.Template
}

func New(ctrl *web.Controller) (*Handler, error) {
	cache, err := newTemplateCache()
	if err != nil {
		return nil, err
	}

	return &Handler{ctrl, cache}, nil
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.NotFound(w)
		return
	}

	data := h.newTemplateData(r)

	h.render(w, http.StatusOK, "home.html", data)
}

func (h *Handler) Catalog(w http.ResponseWriter, r *http.Request) {
	categories, err := h.ctrl.Catalog(r.Context())
	if err != nil {
		h.ServerError(w, err)
		return
	}

	data := h.newTemplateData(r)
	data.Categories = categories

	h.render(w, http.StatusOK, "catalog.html", data)
}

func (h *Handler) ProductsByCategory(w http.ResponseWriter, r *http.Request) {
	id, err := h.getID(r)
	if err != nil || id < 1 {
		h.NotFound(w)
		return
	}

	products, err := h.ctrl.ProductsByCategoryID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, web.ErrNotFound):
			h.NotFound(w)
		default:
			h.ServerError(w, err)
		}
		return
	}

	data := h.newTemplateData(r)
	data.Products = products

	h.render(w, http.StatusOK, "category.html", data)
}

func (h *Handler) Product(w http.ResponseWriter, r *http.Request) {
	id, err := h.getID(r)
	if err != nil || id < 1 {
		h.NotFound(w)
		return
	}

	product, err := h.ctrl.ProductByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, web.ErrNotFound):
			h.NotFound(w)
		default:
			h.ServerError(w, err)
		}
		return
	}

	processedProduct, err := transformProductAttributes(product)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	data := h.newTemplateData(r)
	data.Product = processedProduct

	h.render(w, http.StatusOK, "product.html", data)
}
