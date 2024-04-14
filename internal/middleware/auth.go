package middleware

import (
	"banner-service/internal/auth"
	"banner-service/internal/models"
	"fmt"
	"github.com/samber/lo"
	"log"
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

		resources, err := auth.ValidateToken(authFields[1], am.publicKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		log.Println(r.RequestURI)

		ctx := auth.SetRole(r.Context(), resources.Role)
		if resources.Role == models.Admin {
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		if lo.Contains(resources.Resources, buildResource(r)) {
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		w.WriteHeader(http.StatusForbidden)
	})
}

func buildResource(r *http.Request) string {
	return fmt.Sprintf("%s %s", r.Method, r.URL.Path)
}
