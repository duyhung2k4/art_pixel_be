package router

import (
	"app/config"
	"app/controller"
	middlewares "app/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

func RouterV1(router chi.Router) {
	authController := controller.NewAuthController()
	eventController := controller.NewEventController()
	middlewares := middlewares.NewMiddlewares()

	router.Route("/auth", func(auth chi.Router) {
		auth.Post("/register", authController.Register)
		auth.Post("/auth-face", authController.AuthFace)
		auth.Post("/accept-code", authController.AcceptCode)
		auth.Post("/save-process", authController.SaveProcess)
		auth.Post("/send-file-auth", authController.SendFileAuth)
		auth.Post("/create-socket-auth-face", authController.CreateSocketAuthFace)
	})

	router.Route("/event", func(event chi.Router) {
		event.Use(jwtauth.Verifier(config.GetJWT()))
		event.Use(middlewares.ValidateExpAccessToken())

		event.Get("/all", eventController.GetAllEvent)
		event.Post("/draw", eventController.DrawPixel)
		event.Post("/new-event", eventController.CreateEvent)
	})
}
