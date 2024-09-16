package controller

import (
	"app/config"
	"app/dto/request"
	"app/service"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
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
	SendFileAuth(w http.ResponseWriter, r *http.Request)
}

func (c *authController) Register(w http.ResponseWriter, r *http.Request) {
	var registerReq request.RegisterReq

	if err := json.NewDecoder(r.Body).Decode(&registerReq); err != nil {
		badRequest(w, r, err)
		return
	}

	existProfile, err := c.authService.CheckExistProfile(registerReq)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	if existProfile {
		internalServerError(w, r, errors.New("profile exist"))
		return
	}

	newProfile, err := c.authService.CreateProfilePending(registerReq)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	dataJsonString, err := json.Marshal(newProfile)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	uuidKey := uuid.New().String()
	err = c.redisClient.SetNX(context.Background(), uuidKey, dataJsonString, 24*time.Hour).Err()
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	newFolderPending := fmt.Sprintf("pending_file/%s", uuidKey)
	if err := os.Mkdir(newFolderPending, 0075); err != nil {
		internalServerError(w, r, err)
		return
	}
	newFolderAddModel := fmt.Sprintf("file_add_model/%s", uuidKey)
	if err := os.Mkdir(newFolderAddModel, 0075); err != nil {
		internalServerError(w, r, err)
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

func (c *authController) SendFileAuth(w http.ResponseWriter, r *http.Request) {
	var fileReq request.SendFileAuthFaceReq
	err := json.NewDecoder(r.Body).Decode(&fileReq)
	if err != nil {
		badRequest(w, r, err)
		return
	}

	res := Response{
		Data:    nil,
		Message: "OK",
		Status:  200,
		Error:   nil,
	}

	render.JSON(w, r, res)
}

func NewAuthController() AuthController {
	return &authController{
		redisClient: config.GetRedisClient(),
		authService: service.NewAuthService(),
	}
}
