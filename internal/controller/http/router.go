package http

import (
	"banner-service/internal/middleware"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (ctr *Controller) NewRouter() http.Handler {
	r := chi.NewRouter()
	authMiddleware := middleware.NewAuthMiddleware(ctr.publicKey)
	r.Post("/sign-up", ctr.SignUp)
	r.Post("/sign-in", ctr.SignIn)
	r.With(authMiddleware.Middleware).Post("/user_banner", ctr.GetBanner)
	//TODO: add routes for banner

	return r
}
