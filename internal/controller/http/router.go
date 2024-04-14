package http

import (
	"banner-service/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	_ "net/http/pprof"
)

func (ctr *Controller) NewRouter() http.Handler {
	r := chi.NewRouter()
	r.Mount("/debug/pprof", http.DefaultServeMux)
	r.Mount("/metrics", promhttp.Handler())
	r.With(middleware.MetricsMiddleware).Route("/", func(r chi.Router) {
		r.Post("/sign-up", ctr.SignUpEndpoint)
		r.Post("/sign-in", ctr.SignInEndpoint)
		authMiddleware := middleware.NewAuthMiddleware(ctr.publicKey)
		r.With(authMiddleware.Middleware).Route("/", func(r chi.Router) {
			r.Get("/user_banner", ctr.GetBannerEndpoint)
			r.Route("/banner", func(r chi.Router) {
				r.Get("/", ctr.GetFilteredBannersEndpoint)
				r.Get("/versions/{banner_id}", ctr.GetListOfVersionsEndpoint)
				r.Post("/", ctr.CreateBannerEndpoint)
				r.Patch("/{banner_id}", ctr.PartialUpdateBannerEndpoint)
				r.Patch("/{banner_id}/version/{version}", ctr.ChooseBannerVersionEndpoint)
				r.Delete("/{banner_id}", ctr.DeleteBannerEndpoint)
				r.Delete("/", ctr.MarkBannerAsDeletedEndpoint)
			})
		})
	})

	return r
}
