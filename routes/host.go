package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/hostrouter"
	"github.com/ichtrojan/cook/config"
)

func AllRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	hr := hostrouter.New()

	apiHost := func() string {
		if config.DomainConfig.ApiHost == "" {
			return "*"
		}

		return config.DomainConfig.ApiHost
	}()

	hr.Map(apiHost, APIRoutes())

	r.Mount("/", hr)

	return r
}
