package http

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/Maksim-Kot/Tech-store-web/internal/controller"
	"github.com/Maksim-Kot/Tech-store-web/internal/controller/web"
	"github.com/Maksim-Kot/Tech-store-web/internal/model"
	"github.com/Maksim-Kot/Tech-store-web/internal/session"
	"github.com/Maksim-Kot/Tech-store-web/internal/validator"

	"github.com/go-playground/form/v4"
)

type Handler struct {
	Ctrl           *web.Controller
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	SessionManager session.Manager
}

func New(ctrl *web.Controller, sm session.Manager) (*Handler, error) {
	cache, err := newTemplateCache()
	if err != nil {
		return nil, err
	}

	return &Handler{
		Ctrl:           ctrl,
		templateCache:  cache,
		formDecoder:    form.NewDecoder(),
		SessionManager: sm,
	}, nil
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
	categories, err := h.Ctrl.Catalog.Catalog(r.Context())
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

	products, err := h.Ctrl.Catalog.ProductsByCategoryID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, controller.ErrNotFound):
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

	product, err := h.Ctrl.Catalog.ProductByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, controller.ErrNotFound):
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

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (h *Handler) UserSignup(w http.ResponseWriter, r *http.Request) {
	data := h.newTemplateData(r)
	data.Form = userSignupForm{}

	h.render(w, http.StatusOK, "signup.html", data)
}

func (h *Handler) UserSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	err := h.decodePostForm(r, &form)
	if err != nil {
		h.ClientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := h.newTemplateData(r)
		data.Form = form
		h.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	err = h.Ctrl.User.InsertUser(r.Context(), form.Name, form.Email, form.Password)
	if err != nil {
		switch {
		case errors.Is(err, controller.ErrDuplicateEmail):
			form.AddFieldError("email", "Email address is already in use")

			data := h.newTemplateData(r)
			data.Form = form
			h.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		default:
			h.ServerError(w, err)
		}
		return
	}

	h.SessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (h *Handler) UserLogin(w http.ResponseWriter, r *http.Request) {
	data := h.newTemplateData(r)
	data.Form = userLoginForm{}

	h.render(w, http.StatusOK, "login.html", data)
}

func (h *Handler) UserLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	err := h.decodePostForm(r, &form)
	if err != nil {
		h.ClientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := h.newTemplateData(r)
		data.Form = form
		h.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	id, err := h.Ctrl.User.AuthenticateUser(r.Context(), form.Email, form.Password)
	if err != nil {
		if errors.Is(err, controller.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

			data := h.newTemplateData(r)
			data.Form = form
			h.render(w, http.StatusUnprocessableEntity, "login.html", data)
		} else {
			h.ServerError(w, err)
		}
		return
	}

	err = h.SessionManager.RenewToken(r.Context())
	if err != nil {
		h.ServerError(w, err)
		return
	}

	h.SessionManager.Put(r.Context(), "authenticatedUserID", id)

	path := h.SessionManager.PopString(r.Context(), "redirectPathAfterLogin")
	if path != "" {
		http.Redirect(w, r, path, http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) UserLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := h.SessionManager.RenewToken(r.Context())
	if err != nil {
		h.ServerError(w, err)
		return
	}

	h.SessionManager.Remove(r.Context(), "authenticatedUserID")

	h.SessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) AccountView(w http.ResponseWriter, r *http.Request) {
	id := h.SessionManager.GetInt64(r.Context(), "authenticatedUserID")

	user, err := h.Ctrl.User.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, controller.ErrNotFound) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		} else {
			h.ServerError(w, err)
		}
		return
	}

	data := h.newTemplateData(r)
	data.User = user

	h.render(w, http.StatusOK, "account.html", data)
}

func (h *Handler) ShowCart(w http.ResponseWriter, r *http.Request) {
	var cart model.Cart
	cartData := h.SessionManager.Get(r.Context(), "cart")
	if cartData != nil {
		cart = cartData.(model.Cart)
	}

	data := h.newTemplateData(r)
	data.Cart = &cart

	h.render(w, http.StatusOK, "cart.html", data)
}

func (h *Handler) AddToCart(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	name := r.FormValue("name")
	quantityStr := r.FormValue("quantity")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.ClientError(w, http.StatusBadRequest)
		return
	}

	quantity, err := strconv.ParseInt(quantityStr, 10, 32)
	if err != nil {
		h.ClientError(w, http.StatusBadRequest)
		return
	}

	var cart model.Cart
	cartData := h.SessionManager.Get(r.Context(), "cart")
	if cartData != nil {
		cart = cartData.(model.Cart)
		if cart.Items == nil {
			cart.Items = make(map[int64]model.Item)
		}
	} else {
		cart = model.Cart{
			Items: make(map[int64]model.Item),
		}
	}

	item, exists := cart.Items[id]
	if !exists {
		cart.Items[id] = model.Item{
			ID:       id,
			Name:     name,
			Quantity: int32(quantity),
		}
	} else {
		item.Quantity += int32(quantity)
		cart.Items[id] = item
	}

	h.SessionManager.Put(r.Context(), "cart", cart)

	h.SessionManager.Put(r.Context(), "flash", "Added to cart")

	http.Redirect(w, r, fmt.Sprintf("/product/%s", idStr), http.StatusSeeOther)
}

func (h *Handler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	id, err := h.getID(r)
	if err != nil || id < 1 {
		h.NotFound(w)
		return
	}

	cartData := h.SessionManager.Get(r.Context(), "cart")
	if cartData == nil {
		h.SessionManager.Put(r.Context(), "flash", "Cart is empty")
		http.Redirect(w, r, "/cart", http.StatusSeeOther)
		return
	}

	cart := cartData.(model.Cart)

	if _, exists := cart.Items[id]; exists {
		delete(cart.Items, id)
		h.SessionManager.Put(r.Context(), "flash", "Item removed from cart")
	} else {
		h.SessionManager.Put(r.Context(), "flash", "Item not found in cart")
	}

	h.SessionManager.Put(r.Context(), "cart", cart)

	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}

func (h *Handler) Order(w http.ResponseWriter, r *http.Request) {
	id, err := h.getID(r)
	if err != nil || id < 1 {
		h.NotFound(w)
		return
	}

	purchase, err := h.Ctrl.Orders.OrderByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, controller.ErrNotFound):
			h.NotFound(w)
		default:
			h.ServerError(w, err)
		}
		return
	}

	userID := h.SessionManager.GetInt64(r.Context(), "authenticatedUserID")

	if userID != purchase.UserID {
		h.NotFound(w)
		return
	}

	order := model.Order{
		Status:    purchase.Status,
		CreatedAt: purchase.CreatedAt,
		Price:     purchase.Price,
	}

	for _, item := range purchase.Items {
		product, err := h.Ctrl.Catalog.ProductByID(r.Context(), item.ItemID)
		if err != nil {
			h.ServerError(w, err)
			return
		}

		order.Products = append(order.Products, &model.Product{
			ID:       product.ID,
			Name:     product.Name,
			Quantity: item.Quantity,
		})
	}

	data := h.newTemplateData(r)
	data.Order = &order

	h.render(w, http.StatusOK, "order.html", data)
}

func (h *Handler) OrdersByUser(w http.ResponseWriter, r *http.Request) {
	id := h.SessionManager.GetInt64(r.Context(), "authenticatedUserID")
	if id == 0 {
		h.NotFound(w)
		return
	}

	purchases, err := h.Ctrl.Orders.OrdersByUserID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, controller.ErrNotFound):
			h.NotFound(w)
		default:
			h.ServerError(w, err)
		}
		return
	}

	var orders []*model.Order

	for _, order := range purchases {
		orders = append(orders, &model.Order{
			ID:        order.ID,
			Price:     order.Price,
			Status:    order.Status,
			CreatedAt: order.CreatedAt,
		})
	}

	data := h.newTemplateData(r)
	data.Orders = orders

	h.render(w, http.StatusOK, "orders.html", data)
}
