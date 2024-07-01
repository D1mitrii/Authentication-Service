package v1

import (
	"encoding/json"
	"fmt"
	"github.com/d1mitrii/authentication-service/internal/controller/http/middlewares"
	"github.com/d1mitrii/authentication-service/internal/models"
	"github.com/d1mitrii/authentication-service/internal/repository/repoerrors"
	"github.com/d1mitrii/authentication-service/internal/services"
	"net/http"
)

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "incorrect request body", http.StatusBadRequest)
		return
	}
	id, err := h.service.Register(r.Context(), user)
	if err != nil {
		if err == repoerrors.ErrAlreadyExist {
			http.Error(w, "user already exist", http.StatusBadRequest)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	type response struct {
		Id int `json:"id"`
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response{id})
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) logIn(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "incorrect request body", http.StatusBadRequest)
		return
	}

	jwt, err := h.service.Login(r.Context(), user)

	if err != nil {
		switch err {
		case services.ErrCannotSignToken:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		case services.ErrIncorrectPassword:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     middlewares.RefreshCookie,
		Value:    jwt.Refresh,
		Path:     "/",
		MaxAge:   int(h.service.JWT.RefreshTTL().Seconds()),
		HttpOnly: true,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jwt)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) logOut(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.Context().Value(middlewares.CtxRefreshToken{}).(string)
	if h.service.Logout(r.Context(), refreshToken) != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     middlewares.RefreshCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.Context().Value(middlewares.CtxRefreshToken{}).(string)
	jwt, err := h.service.RefreshSession(r.Context(), refreshToken)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     middlewares.RefreshCookie,
		Value:    jwt.Refresh,
		Path:     "/",
		HttpOnly: true,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jwt)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) secret(w http.ResponseWriter, r *http.Request) {
	helloStr := fmt.Sprintf("Hi, user: %d", r.Context().Value(middlewares.CtxUserId{}).(int))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(helloStr)
	w.WriteHeader(http.StatusOK)
}
