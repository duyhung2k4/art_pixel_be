package rabbitmq

import (
	"app/config"
	"app/constant"
	queuepayload "app/dto/queue_payload"
	"app/dto/response"
	"app/model"
	"app/service"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type queueAuth struct {
	mapSocket    map[string]*websocket.Conn
	rabbitmq     *amqp091.Connection
	mutex        *sync.Mutex
	authService  service.AuthService
	eventService service.EventService
	psql         *gorm.DB
}

type QueueAuth interface {
	InitQueueSendFileAuth()
	InitQueueAuthFace()
	sendMess(data interface{}, socket *websocket.Conn)
}

func (q *queueAuth) InitQueueSendFileAuth() {
	ch, err := q.rabbitmq.Channel()
	if err != nil {
		log.Println("Failed to open a channel:", err)
		return
	}
	defer ch.Close()

	queueName := fmt.Sprint(constant.SEND_FILE_AUTH_QUEUE)
	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Println("Failed to declare a queue:", err)
		return
	}

	msgs, err := ch.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Println("Failed to register a consumer:", err)
		return
	}

	var wg sync.WaitGroup
	for msg := range msgs {
		wg.Add(1)
		go func(msg amqp091.Delivery) {
			defer wg.Done()
			msg.Ack(false)

			var dataMess queuepayload.SendFileAuthMess
			if err := json.Unmarshal(msg.Body, &dataMess); err != nil {

				return
			}

			socket := q.mapSocket[dataMess.Uuid]
			if socket == nil {

				return
			}

			result, err := q.authService.CheckFace(dataMess)
			if err != nil {
				socket.WriteMessage(websocket.TextMessage, []byte(err.Error()))

				return
			}

			socket.WriteMessage(websocket.TextMessage, []byte(result))

		}(msg)
	}
}

func (q *queueAuth) InitQueueAuthFace() {
	ch, err := q.rabbitmq.Channel()
	if err != nil {
		log.Println("Failed to open a channel:", err)
		return
	}
	defer ch.Close()

	queueName := fmt.Sprint(constant.FACE_AUTH_QUEUE)
	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Println("Failed to declare a queue:", err)
		return
	}

	msgs, err := ch.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Println("Failed to register a consumer:", err)
		return
	}

	var wg sync.WaitGroup
	for msg := range msgs {
		wg.Add(1)
		go func(msg amqp091.Delivery) {
			defer wg.Done()
			msg.Ack(false)
			var res response.SocketErrorRes

			var dataMess queuepayload.FaceAuth
			if err := json.Unmarshal(msg.Body, &dataMess); err != nil {
				return
			}

			socket := q.mapSocket[dataMess.Uuid]
			if socket == nil {
				return
			}

			result, err := q.authService.AuthFace(dataMess)
			if err != nil {
				res.Error = err
				q.sendMess(res, socket)
				return
			}

			if result <= 0 {
				res.Error = errors.New("profile not found")
				q.sendMess(res, socket)
				return
			}

			accessToken, refreshToken, err := q.authService.CreateToken(uint(result))
			if err != nil {
				res.Error = err
				q.sendMess(res, socket)
				return
			}

			var profile *model.Profile
			if err = q.psql.Model(&model.Profile{}).Where("id = ?", uint(result)).First(&profile).Error; err != nil {
				res.Error = err
				q.sendMess(res, socket)
				return
			}

			res = response.SocketErrorRes{
				Data: map[string]interface{}{
					"accessToken":  accessToken,
					"refreshToken": refreshToken,
					"profileId":    result,
					"profile":      profile,
				},
				Error: nil,
			}
			q.sendMess(res, socket)
		}(msg)
	}
}

func (q *queueAuth) sendMess(data interface{}, socket *websocket.Conn) {
	dataByte, _ := json.Marshal(data)
	q.mutex.Lock()
	socket.WriteMessage(websocket.TextMessage, dataByte)
	q.mutex.Unlock()
}

func NewQueueAuth() QueueAuth {
	return &queueAuth{
		mutex:        new(sync.Mutex),
		rabbitmq:     config.GetRabbitmq(),
		mapSocket:    config.GetMapSocket(),
		authService:  service.NewAuthService(),
		eventService: service.NewEventService(),
		psql:         config.GetPsql(),
	}
}
