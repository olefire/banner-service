package http

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (ctr *Controller) NewRouter() http.Handler {
	r := chi.NewRouter()
	r.Post("/sign-up", ctr.SignUp)
	r.Post("/sign-in", ctr.SignIn)
	//TODO: add routes for banner

	return r
}
