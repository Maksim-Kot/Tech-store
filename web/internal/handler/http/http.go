package http

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/Maksim-Kot/Tech-store-catalog/pkg/model"
	"github.com/Maksim-Kot/Tech-store-web/internal/controller/web"
)

type Handler struct {
	ctrl *web.Controller
}

func New(ctrl *web.Controller) *Handler {
	return &Handler{ctrl}
}

type templateData struct {
	Categories []*model.Category
	Products   []*model.Product
	Product    *processedProduct
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.notFound(w)
		return
	}

	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/home.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		h.serverError(w, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		h.serverError(w, err)
	}
}

func (h *Handler) Catalog(w http.ResponseWriter, r *http.Request) {
	categories, err := h.ctrl.Catalog(r.Context())
	if err != nil {
		h.serverError(w, err)
		return
	}

	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/catalog.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		h.serverError(w, err)
		return
	}

	data := &templateData{
		Categories: categories,
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		h.serverError(w, err)
	}
}

func (h *Handler) ProductsByCategory(w http.ResponseWriter, r *http.Request) {
	id, err := h.getID(r)
	if err != nil || id < 1 {
		h.notFound(w)
		return
	}

	products, err := h.ctrl.ProductsByCategoryID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, web.ErrNotFound):
			h.notFound(w)
		default:
			h.serverError(w, err)
		}
		return
	}

	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/category.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		h.serverError(w, err)
		return
	}

	data := &templateData{
		Products: products,
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		h.serverError(w, err)
	}
}

func (h *Handler) Product(w http.ResponseWriter, r *http.Request) {
	id, err := h.getID(r)
	if err != nil || id < 1 {
		h.notFound(w)
		return
	}

	product, err := h.ctrl.ProductByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, web.ErrNotFound):
			h.notFound(w)
		default:
			h.serverError(w, err)
		}
		return
	}

	processedProduct, err := transformProductAttributes(product)
	if err != nil {
		h.serverError(w, err)
		return
	}

	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/product.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		h.serverError(w, err)
		return
	}

	data := &templateData{
		Product: processedProduct,
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		h.serverError(w, err)
	}
}
