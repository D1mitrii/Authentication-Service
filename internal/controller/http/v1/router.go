package v1

import (
	"auth/internal/controller/http/middlewares"
	"auth/internal/services"

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
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.Heartbeat("/health"))

	r.Post("/signup", h.signUp)
	r.Post("/login", h.logIn)

	auth := middlewares.New(h.service)

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
