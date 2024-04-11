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
	r.Post("/banner", ctr.CreateBanner)
	r.Patch("/banner/{banner_id}", ctr.PartialUpdateBanner)
	r.Delete("/banner/{banner_id}", ctr.DeleteBanner)
	//TODO: add routes for banner

	return r
}
