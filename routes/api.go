package routes

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/ichtrojan/cook/controllers"
	"github.com/ichtrojan/cook/helpers"
	customMiddleware "github.com/ichtrojan/cook/middlewares"
	"net/http"
	"time"
)

func APIRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*"},
		AllowedMethods:   []string{"POST"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "stripe-signature"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// rate limit by IP Address
	router.Use(httprate.LimitByIP(100, 1*time.Minute))

	// force `content-type` to `application/json`
	router.Use(customMiddleware.AcceptJson)

	// validate json
	router.Use(customMiddleware.ValidateJson)

	// response for 404
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(helpers.PrepareMessage("endpoint not found"))
		return
	})

	router.Group(func(router chi.Router) {
		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(helpers.PrepareMessage("time to cook 🧑‍🍳"))
			return
		})

		// auth
		router.Post("/register", controllers.Register)
		router.Post("/login", controllers.Login)

		// forgot password
		router.Post("/pre-forgot", controllers.PreForgot)
		router.Post("/post-forgot", controllers.PostForgot)
	})

	// verify user
	router.Group(func(router chi.Router) {
		router.Use(customMiddleware.AuthenticateUser)
		router.Use(customMiddleware.ValidateJson)
		router.Post("/verify", controllers.VerifyUser)
	})

	router.Group(func(router chi.Router) {
		router.Use(customMiddleware.AuthenticateUser)
		router.Use(customMiddleware.VerifyUser)

		// get user
		router.Get("/user", controllers.GetUser)
	})

	return router
}
