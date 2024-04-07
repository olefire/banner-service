package middleware

import (
	utils "banner-service/pkg/utils/auth"
	"context"
	"net/http"
	"strings"
)

func AuthorizationMiddleware(next http.Handler, publicKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		authFields := strings.Fields(authHeader)

		if len(authFields) == 0 || authFields[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		_, err := utils.ValidateToken(authFields[1], publicKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func AdminMiddleware(next http.Handler, publicKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		authFields := strings.Fields(authHeader)

		if len(authFields) == 0 || authFields[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		isAdmin, err := utils.ValidateToken(authFields[1], publicKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "isAdmin", isAdmin.(bool))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
