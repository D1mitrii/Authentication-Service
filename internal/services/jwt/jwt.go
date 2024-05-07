package jwt

import (
	"auth/internal/models"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	Id int `json:"id"`
	jwt.RegisteredClaims
}

type JWT struct {
	secret_key  []byte
	access_ttl  time.Duration
	refresh_ttl time.Duration
}

func New(
	secret string,
	access_ttl time.Duration,
	refresh_ttl time.Duration,
) *JWT {
	return &JWT{
		secret_key:  []byte(secret),
		access_ttl:  access_ttl,
		refresh_ttl: refresh_ttl,
	}
}

func (r *JWT) RefreshTTL() time.Duration {
	return r.refresh_ttl
}

func (r *JWT) NewAccessToken(user models.User) (string, error) {
	claims := &TokenClaims{
		user.Id,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(r.access_ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(r.secret_key)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func (r *JWT) NewRefreshToken() (string, error) {
	// Getting random bytes
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", nil
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(r.refresh_ttl)),
	})
	tokenStr, err := token.SignedString(b)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func (r *JWT) Parse(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return r.secret_key, nil
	})

	if err != nil {
		return 0, nil
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return 0, fmt.Errorf("failed to map token")
	}
	return claims.Id, nil
}
