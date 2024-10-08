package controller

import (
	"app/config"
	"app/constant"
	queuepayload "app/dto/queue_payload"
	"app/dto/request"
	"app/service"
	"app/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type authController struct {
	rabbitmq          *amqp091.Connection
	redisClient       *redis.Client
	authService       service.AuthService
	smtpService       service.SmtpService
	mapCheckSendEmail map[string]bool
	mutex             *sync.Mutex
	jwtUtils          utils.JwtUtils
	rdb               *redis.Client
	utils             utils.JwtUtils
}

type AuthController interface {
	Register(w http.ResponseWriter, r *http.Request)
	SendFileAuth(w http.ResponseWriter, r *http.Request)
	AuthFace(w http.ResponseWriter, r *http.Request)
	CreateSocketAuthFace(w http.ResponseWriter, r *http.Request)
	AcceptCode(w http.ResponseWriter, r *http.Request)
	SaveProcess(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)
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

	newFolderPending := fmt.Sprintf("file/pending_file/%s", uuidKey)
	if err := os.Mkdir(newFolderPending, 0777); err != nil {
		internalServerError(w, r, err)
		return
	}
	newFolderAddModel := fmt.Sprintf("file/file_add_model/%s", uuidKey)
	if err := os.Mkdir(newFolderAddModel, 0777); err != nil {
		internalServerError(w, r, err)
		return
	}
	c.mapCheckSendEmail[uuidKey] = false

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

	key := fmt.Sprintf("code_%s", uuid)
	code, err := c.redisClient.Get(r.Context(), key).Result()
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	if code != payload.Code {
		internalServerError(w, r, errors.New("error code"))
		return
	}

	profile, err := c.authService.ActiveProfile(uuid)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	res := Response{
		Data:    profile.PublicKey,
		Message: "OK",
		Status:  200,
		Error:   nil,
	}

	render.JSON(w, r, res)
}

func (c *authController) SaveProcess(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.Header.Get("authorization"), " ")[1]

	if c.mapCheckSendEmail[uuid] {
		internalServerError(w, r, errors.New("email sending"))
		return
	}

	c.mutex.Lock()
	c.mapCheckSendEmail[uuid] = true
	c.mutex.Unlock()

	if err := c.authService.SaveFileAuth(uuid); err != nil {
		internalServerError(w, r, err)
		return
	}
	if err := c.smtpService.SendCodeAcceptRegister(uuid); err != nil {
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

func (c *authController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	tokenString := strings.Split(r.Header.Get("Authorization"), " ")[1]
	mapDataRequest, errMapData := c.utils.JwtDecode(tokenString)

	if errMapData != nil {
		internalServerError(w, r, errMapData)
		return
	}

	profileId := uint(mapDataRequest["profile_id"].(float64))
	profileResponse, errProfile := c.authService.GetProfile(profileId)
	if errProfile != nil {
		internalServerError(w, r, errProfile)
		return
	}

	mapData := map[string]interface{}{
		"profile_id": profileResponse.ID,
		"email":      profileResponse.Email,
	}

	accessData := mapData
	accessData["uuid"] = uuid.New()
	accessData["exp"] = time.Now().Add(3 * time.Hour).Unix()
	accessToken, errAccessToken := c.jwtUtils.JwtEncode(accessData)
	if errAccessToken != nil {
		internalServerError(w, r, errAccessToken)
		return
	}

	refreshData := mapData
	refreshData["uuid"] = uuid.New()
	refreshData["exp"] = time.Now().Add(3 * 3 * time.Hour).Unix()
	refreshToken, errRefreshToken := c.jwtUtils.JwtEncode(refreshData)
	if errRefreshToken != nil {
		internalServerError(w, r, errRefreshToken)
		return
	}

	errSetKeyAccessToken := c.rdb.Set(context.Background(), "access_token:"+strconv.Itoa(int(profileResponse.ID)), accessToken, 24*time.Hour).Err()
	if errSetKeyAccessToken != nil {
		internalServerError(w, r, errSetKeyAccessToken)
		return
	}
	errSetKeyRefreshToken := c.rdb.Set(context.Background(), "refresh_token:"+strconv.Itoa(int(profileResponse.ID)), refreshToken, 3*24*time.Hour).Err()
	if errSetKeyRefreshToken != nil {
		internalServerError(w, r, errSetKeyRefreshToken)
		return
	}

	res := Response{
		Data: map[string]interface{}{
			"profile":      profileResponse,
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		},
	}

	render.JSON(w, r, res)

}

func NewAuthController() AuthController {
	return &authController{
		mutex:             new(sync.Mutex),
		rabbitmq:          config.GetRabbitmq(),
		redisClient:       config.GetRedisClient(),
		authService:       service.NewAuthService(),
		smtpService:       service.NewSmtpService(),
		mapCheckSendEmail: config.GetMapCheckSendEmail(),
		jwtUtils:          utils.NewJwtUtils(),
		rdb:               config.GetRedisClient(),
		utils:             utils.NewJwtUtils(),
	}
}
