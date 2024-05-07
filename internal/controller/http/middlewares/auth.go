package middlewares

import (
	"context"
	"net/http"
	"strings"
)

type CtxUserId struct{}
type CtxRefreshToken struct{}

const (
	RefreshCookie string = "refresh-token"
)

type JWT interface {
	Parse(token string) (int, error)
}

type AuthMiddleware struct {
	jwt JWT
}

func NewAuthMiddleware(jwt JWT) *AuthMiddleware {
	return &AuthMiddleware{
		jwt: jwt,
	}
}

func (m *AuthMiddleware) JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := getBearerToken(r.Header.Get("Authorization"))
		if !ok {
			http.Error(w, "incorrect authorization header", http.StatusUnauthorized)
			return
		}
		id, err := m.jwt.Parse(token)
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
