package controller

import (
	"app/config"
	"app/constant"
	queuepayload "app/dto/queue_payload"
	"app/dto/request"
	"app/service"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type authController struct {
	rabbitmq    *amqp091.Connection
	authService service.AuthService
	redisClient *redis.Client
}

type AuthController interface {
	Register(w http.ResponseWriter, r *http.Request)
	SendFileAuth(w http.ResponseWriter, r *http.Request)
	AuthFace(w http.ResponseWriter, r *http.Request)
	CreateSocketAuthFace(w http.ResponseWriter, r *http.Request)
	AcceptCode(w http.ResponseWriter, r *http.Request)
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
	uuid := strings.Split(r.Header.Get("authorization"), " ")[1]

	if len(uuid) == 0 {
		badRequest(w, r, errors.New("not found uuid"))
		return
	}

	err := json.NewDecoder(r.Body).Decode(&fileReq)
	if err != nil {
		badRequest(w, r, err)
		return
	}

	ch, err := c.rabbitmq.Channel()
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	dataMess := queuepayload.SendFileAuthMess{
		Data: fileReq.Data,
		Uuid: uuid,
	}

	dataMessString, err := json.Marshal(dataMess)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	err = ch.PublishWithContext(
		context.Background(),
		"",
		string(constant.SEND_FILE_AUTH_QUEUE),
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(dataMessString),
		},
	)

	if err != nil {
		internalServerError(w, r, err)
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

func (c *authController) AuthFace(w http.ResponseWriter, r *http.Request) {
	var authFaceReq request.AuthFaceReq
	if err := json.NewDecoder(r.Body).Decode(&authFaceReq); err != nil {
		badRequest(w, r, err)
		return
	}

	uuid := strings.Split(r.Header.Get("authorization"), " ")[1]
	if len(uuid) == 0 {
		badRequest(w, r, errors.New("not found uuid"))
		return
	}

	path, err := c.authService.CreateFileAuthFace(authFaceReq)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	ch, err := c.rabbitmq.Channel()
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	dataMess := queuepayload.FaceAuth{
		FilePath: path,
		Uuid:     uuid,
	}

	dataMessString, err := json.Marshal(dataMess)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	err = ch.PublishWithContext(
		context.Background(),
		"",
		string(constant.FACE_AUTH_QUEUE),
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(dataMessString),
		},
	)

	if err != nil {
		internalServerError(w, r, err)
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

func (c *authController) CreateSocketAuthFace(w http.ResponseWriter, r *http.Request) {
	uuid := uuid.New().String()

	res := Response{
		Data:    uuid,
		Message: "OK",
		Status:  200,
		Error:   nil,
	}

	render.JSON(w, r, res)
}

func (c *authController) AcceptCode(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.Header.Get("authorization"), " ")[1]

	if len(uuid) == 0 {
		badRequest(w, r, errors.New("not found uuid"))
		return
	}

	var payload request.AcceptCodeReq
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		badRequest(w, r, err)
		return
	}

	code, err := c.redisClient.Get(r.Context(), uuid).Result()
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	if code != payload.Code {
		internalServerError(w, r, errors.New("error code"))
		return
	}

	err = c.authService.ActiveProfile(uuid)
	if err != nil {
		internalServerError(w, r, err)
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
		rabbitmq:    config.GetRabbitmq(),
		redisClient: config.GetRedisClient(),
		authService: service.NewAuthService(),
	}
}
