package http

import (
	"banner-service/internal/middleware"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (ctr *Controller) NewRouter() http.Handler {
	r := chi.NewRouter()
	authMiddleware := middleware.NewAuthMiddleware(ctr.publicKey)
	r.Post("/sign-up", ctr.SignUpEndpoint)
	r.Post("/sign-in", ctr.SignInEndpoint)
	r.With(authMiddleware.Middleware).Get("/user_banner", ctr.GetBannerEndpoint)
	r.With(authMiddleware.Middleware).Route("/banner", func(r chi.Router) {
		r.Get("/", ctr.GetFilteredBannersEndpoint)
		r.Get("/versions/{banner_id}", ctr.GetListOfVersionsEndpoint)
		r.Post("/", ctr.CreateBannerEndpoint)
		r.Patch("/{banner_id}", ctr.PartialUpdateBannerEndpoint)
		r.Patch("/choose/{version}", ctr.ChooseBannerVersionEndpoint)
		r.Delete("/banner/{banner_id}", ctr.DeleteBannerEndpoint)
	})

	//TODO: add routes for banner

	return r
}
