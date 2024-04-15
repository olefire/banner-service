package middleware

import (
	"banner-service/internal/auth"
	"banner-service/internal/models"
	"fmt"
	"github.com/jellydator/ttlcache/v3"
	"github.com/samber/lo"
	"log"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	auth.TokenProvider
}

func NewAuthMiddleware(tp auth.TokenProvider) *AuthMiddleware {
	return &AuthMiddleware{TokenProvider: tp}
}

func (am *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		authFields := strings.Fields(authHeader)

		if len(authFields) == 0 || authFields[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token := authFields[1]

		cachedResources := am.TokenProvider.Cache.Get(token)

		if cachedResources == nil {
			validateResources, err := am.ValidateToken(token, am.PublicKey)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			log.Println("set cache[token]model.UserResources")
			cachedResources = am.TokenProvider.Cache.Set(token, validateResources, ttlcache.DefaultTTL)
		}

		resources := cachedResources.Value()

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
