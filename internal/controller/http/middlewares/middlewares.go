package middlewares

import (
	"auth/internal/services"
	"context"
	"net/http"
	"strings"
)

type CtxUserId struct{}
type CtxRefreshToken struct{}

const (
	RefreshCookie string = "refresh-token"
)

type AuthMiddleware struct {
	service *services.Services
}

func New(services *services.Services) *AuthMiddleware {
	return &AuthMiddleware{
		service: services,
	}
}

func (m *AuthMiddleware) JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := getBearerToken(r.Header.Get("Authorization"))
		if !ok {
			http.Error(w, "incorrect authorization header", http.StatusUnauthorized)
			return
		}
		id, err := m.service.JWT.Parse(token)
		if err != nil {
			http.Error(w, "incorrect access token", http.StatusForbidden)
			return
		}
		ctx := context.WithValue(r.Context(), CtxUserId{}, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getBearerToken(header string) (string, bool) {
	splitHeader := strings.Split(header, "Bearer ")
	if len(splitHeader) != 2 {
		return "", false
	}
	return splitHeader[1], true
}

func (m *AuthMiddleware) RefreshTokenCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(RefreshCookie)
		if err != nil || len(cookie.Value) == 0 {
			http.Error(w, "refresh token cookie didn't provided", http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), CtxRefreshToken{}, cookie.Value)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
