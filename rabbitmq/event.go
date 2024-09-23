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

type queueEvent struct {
	mapSocket    map[string]*websocket.Conn
	rabbitmq     *amqp091.Connection
	mutex        *sync.Mutex
	eventService service.EventService
}

type QueueEvent interface {
	InitQueueDrawPixel()
	sendMess(data interface{}, socket *websocket.Conn)
}

func (q *queueEvent) InitQueueDrawPixel() {
	ch, err := q.rabbitmq.Channel()
	if err != nil {
		log.Println("Failed to open a channel:", err)
		return
	}
	defer ch.Close()

	queueName := fmt.Sprint(constant.DRAW_PIXEL_QUEUE)
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

			var dataMess queuepayload.DrawPixel
			if err := json.Unmarshal(msg.Body, &dataMess); err != nil {
				return
			}

			q.eventService.DrawPixel(dataMess, 1)
		}(msg)
	}
}

func (q *queueEvent) sendMess(data interface{}, socket *websocket.Conn) {
	dataByte, _ := json.Marshal(data)
	q.mutex.Lock()
	socket.WriteMessage(websocket.TextMessage, dataByte)
	q.mutex.Unlock()
}

func NewQueueEvent() QueueEvent {
	return &queueEvent{
		mutex:        new(sync.Mutex),
		rabbitmq:     config.GetRabbitmq(),
		mapSocket:    config.GetMapSocket(),
		eventService: service.NewEventService(),
	}
}
