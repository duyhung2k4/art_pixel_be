package socket

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
)

func ServerSocker() http.Handler {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	app := chi.NewRouter()

	// A good base middleware stack
	app.Use(middleware.RequestID)
	app.Use(middleware.RealIP)
	app.Use(middleware.Logger)
	app.Use(middleware.Recoverer)

	app.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) { AuthSocket(w, r, upgrader) })

	log.Printf("Socket art-pixel starting success!")

	return app
}
