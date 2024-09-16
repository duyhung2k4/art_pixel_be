package socket

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func ServerSocker() http.Handler {
	app := chi.NewRouter()

	// A good base middleware stack
	app.Use(middleware.RequestID)
	app.Use(middleware.RealIP)
	app.Use(middleware.Logger)
	app.Use(middleware.Recoverer)

	app.Route("/socket", func(socket chi.Router) {
		socket.Route("/v1", SocketV1)
	})

	log.Printf("Socket art-pixel starting success!")

	return app
}
