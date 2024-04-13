package middleware

import (
	"banner-service/internal/models"
	authUtils "banner-service/pkg/utils/auth"
	contextUtils "banner-service/pkg/utils/context"
	"fmt"
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

		resources, err := authUtils.ValidateToken(authFields[1], am.publicKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		mapResources := resources.(map[string]interface{})

		if role, ok := mapResources["role"]; ok {
			if role == models.Admin {
				ctx := contextUtils.SetPayload(r.Context(), models.Admin)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		if resources, ok := mapResources["resources"]; ok {
			log.Print(resources.([]interface{}), r.Method, r.URL.Path)
			if findResourceIndex(resources.([]interface{}), fmt.Sprintf("%s %s", r.Method, r.URL.Path)) != -1 {
				log.Print(findResourceIndex(resources.([]interface{}), fmt.Sprintf("%s %s", r.Method, r.URL.Path)), r.URL.Path, r.Method)
				ctx := contextUtils.SetPayload(r.Context(), mapResources["role"])
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		} else {
			w.WriteHeader(http.StatusForbidden)
		}
	})
}

func findResourceIndex(resources []interface{}, resource string) int {
	for i, res := range resources {
		if res == resource {
			return i
		}
	}
	return -1
}
