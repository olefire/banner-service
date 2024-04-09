package middleware

import (
	utils "banner-service/pkg/utils/auth"
	"context"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	publicKey string
}

func NewAuthMiddleware(publicKey string) *AuthMiddleware {
	return &AuthMiddleware{publicKey: publicKey}
}

func (am *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		authFields := strings.Fields(authHeader)

		if len(authFields) == 0 || authFields[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		role, err := utils.ValidateToken(authFields[1], am.publicKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "role", role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
