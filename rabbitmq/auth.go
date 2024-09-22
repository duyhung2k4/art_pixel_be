package rabbitmq

import (
	"app/config"
	"app/constant"
	queuepayload "app/dto/queue_payload"
	"app/service"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rabbitmq/amqp091-go"
)

type queueAuth struct {
	rabbitmq    *amqp091.Connection
	mapSocket   map[string]*websocket.Conn
	authService service.AuthService
	mutex       *sync.Mutex
}

type QueueAuth interface {
	InitQueueSendFileAuth()
	InitQueueAuthFace()
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
				q.mutex.Lock()
				socket.WriteMessage(websocket.TextMessage, []byte(err.Error()))
				q.mutex.Unlock()
				return
			}
			q.mutex.Lock()
			socket.WriteMessage(websocket.TextMessage, []byte(fmt.Sprint(result)))
			q.mutex.Unlock()
		}(msg)
	}
}

func NewQueueAuth() QueueAuth {
	return &queueAuth{
		mutex:       new(sync.Mutex),
		rabbitmq:    config.GetRabbitmq(),
		mapSocket:   config.GetMapSocket(),
		authService: service.NewAuthService(),
	}
}
