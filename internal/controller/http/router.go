package http

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (ctr *Controller) NewRouter() http.Handler {
	r := chi.NewRouter()
	//authMiddleware := middleware.NewAuthMiddleware(ctr.publicKey)
	r.Post("/sign-up", ctr.SignUp)
	r.Post("/sign-in", ctr.SignIn)
	r.Get("/user_banner", ctr.GetBanner)
	r.Route("/banner", func(r chi.Router) {
		r.Get("/", ctr.GetFilteredBanners)
		r.Get("/versions/{banner_id}", ctr.GetListOfVersions)
		r.Post("/", ctr.CreateBanner)
		r.Patch("/{banner_id}", ctr.PartialUpdateBanner)
		r.Patch("/choose/{version}", ctr.ChooseBannerVersion)
		r.Delete("/banner/{banner_id}", ctr.DeleteBanner)
	})

	//TODO: add routes for banner

	return r
}
