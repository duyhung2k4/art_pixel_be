package router

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func AppRouter() http.Handler {
	app := chi.NewRouter()

	// A good base middleware stack
	app.Use(middleware.RequestID)
	app.Use(middleware.RealIP)
	app.Use(middleware.Logger)
	app.Use(middleware.Recoverer)

	app.Route("/test", func(test chi.Router) {
		test.Get("/v1", func(w http.ResponseWriter, r *http.Request) {
			render.JSON(w, r, map[string]interface{}{
				"mess": "done",
			})
		})
	})

	log.Println("Sevice h-shop-be-shop starting success! URL: http://localhost:10000")

	return app
}
