package rabbitmq

import (
	"app/config"
	"app/constant"
	queuepayload "app/dto/queue_payload"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/rabbitmq/amqp091-go"
)

type queueAuth struct {
	rabbitmq  *amqp091.Connection
	mapSocket map[string]*websocket.Conn
}

type QueueAuth interface {
	InitQueueSendFileAuth()
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
		queueName, // Tên queue
		"",        // Consumer name (để trống để RabbitMQ tự tạo)
		false,     // Auto-Ack (đặt là false để dùng thủ công acknowledgment)
		false,     // Exclusive (chỉ được dùng cho connection hiện tại)
		false,     // No-local (chỉ dành cho các message local)
		false,     // No-wait (không chờ RabbitMQ trả lời)
		nil,       // Thêm các option khác (nếu cần)
	)

	if err != nil {
		log.Println("Failed to register a consumer:", err)
		return
	}

	for msg := range msgs {
		var dataMess queuepayload.SendFileAuthMess
		if err := json.Unmarshal(msg.Body, &dataMess); err != nil {
			msg.Ack(true)
			continue
		}

		socket := q.mapSocket[dataMess.Uuid]
		if socket == nil {
			msg.Ack(true)
			continue
		}

		socket.WriteMessage(websocket.TextMessage, []byte("done cc"))
		msg.Ack(false)
	}
}

func NewQueueAuth() QueueAuth {
	return &queueAuth{
		rabbitmq:  config.GetRabbitmq(),
		mapSocket: config.GetMapSocket(),
	}
}
