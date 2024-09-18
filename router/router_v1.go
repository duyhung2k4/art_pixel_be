package router

import (
	"app/controller"

	"github.com/go-chi/chi/v5"
)

func RouterV1(router chi.Router) {
	authController := controller.NewAuthController()

	router.Route("/auth", func(auth chi.Router) {
		auth.Post("/register", authController.Register)
		auth.Post("/send-file-auth", authController.SendFileAuth)
	})
}
