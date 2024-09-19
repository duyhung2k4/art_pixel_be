package socket

import (
	"app/controller"

	"github.com/go-chi/chi/v5"
)

func SocketV1(router chi.Router) {
	socketController := controller.NewSocketController()

	router.HandleFunc("/auth", socketController.AuthSocket)
	router.HandleFunc("/login", socketController.FaceLoginSocket)
}
