package controller

import (
	"app/config"
	"app/dto/request"
	"app/service"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type authController struct {
	authService service.AuthService
	redisClient *redis.Client
}

type AuthController interface {
	Register(w http.ResponseWriter, r *http.Request)
}

func (c *authController) Register(w http.ResponseWriter, r *http.Request) {
	var registerReq request.RegisterReq

	if err := json.NewDecoder(r.Body).Decode(&registerReq); err != nil {
		badRequest(w, r, err)
		return
	}

	newProfile, errNewProfile := c.authService.CreateProfilePending(registerReq)
	if errNewProfile != nil {
		internalServerError(w, r, errNewProfile)
		return
	}

	dataJsonString, errJsonString := json.Marshal(newProfile)
	if errJsonString != nil {
		internalServerError(w, r, errJsonString)
		return
	}

	uuidKey := uuid.New().String()
	errSetData := c.redisClient.SetNX(context.Background(), uuidKey, dataJsonString, 2*time.Minute).Err()
	if errSetData != nil {
		internalServerError(w, r, errSetData)
		return
	}

	res := Response{
		Data:    uuidKey,
		Message: "OK",
		Error:   nil,
		Status:  200,
	}

	render.JSON(w, r, res)
}

func NewAuthController() AuthController {
	return &authController{
		redisClient: config.GetRedisClient(),
		authService: service.NewAuthService(),
	}
}
