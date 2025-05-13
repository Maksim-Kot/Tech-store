package http

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/Maksim-Kot/Tech-store-web/internal/controller"
	"github.com/Maksim-Kot/Tech-store-web/internal/controller/web"
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
