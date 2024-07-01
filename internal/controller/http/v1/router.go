package v1

import (
	"github.com/d1mitrii/authentication-service/internal/controller/http/middlewares"
	"github.com/d1mitrii/authentication-service/internal/services"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
	service *services.Services
}

func New(s *services.Services) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.AllowContentType("application/json"))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {})
	r.Post("/signup", h.signUp)
	r.Post("/login", h.logIn)

	auth := middlewares.NewAuthMiddleware(h.service.JWT)

	r.Group(func(r chi.Router) {
		r.Use(auth.RefreshTokenCookie)
		r.Get("/refresh", h.refresh)
		r.Get("/logout", h.logOut)
	})

	r.Group(func(r chi.Router) {
		r.Use(auth.JWT)
		r.Get("/secret", h.secret)
	})

	return r
}
