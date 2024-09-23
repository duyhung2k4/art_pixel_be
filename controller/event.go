package controller

import (
	"app/config"
	"app/constant"
	queuepayload "app/dto/queue_payload"
	"app/dto/request"
	"app/service"
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"github.com/rabbitmq/amqp091-go"
)

type eventController struct {
	eventService service.EventService
	rabbitmq     *amqp091.Connection
}

type EventController interface {
	CreateEvent(w http.ResponseWriter, r *http.Request)
	DrawPixel(w http.ResponseWriter, r *http.Request)
}

func (c *eventController) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent request.CreateEventReq
	if err := json.NewDecoder(r.Body).Decode(&newEvent); err != nil {
		badRequest(w, r, err)
		return
	}

	event, err := c.eventService.CreateEvent(newEvent)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	res := Response{
		Data:    event,
		Message: "OK",
		Status:  200,
		Error:   nil,
	}

	render.JSON(w, r, res)
}

func (c *eventController) DrawPixel(w http.ResponseWriter, r *http.Request) {
	var newPixel request.DrawPixelReq
	if err := json.NewDecoder(r.Body).Decode(&newPixel); err != nil {
		badRequest(w, r, err)
		return
	}

	ch, err := c.rabbitmq.Channel()
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	dataMess := queuepayload.DrawPixel{
		AccessToken: "",
		Data:        newPixel,
	}
	dataMessString, err := json.Marshal(dataMess)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	err = ch.PublishWithContext(
		context.Background(),
		"",
		string(constant.DRAW_PIXEL_QUEUE),
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

func NewEventController() EventController {
	return &eventController{
		eventService: service.NewEventService(),
		rabbitmq:     config.GetRabbitmq(),
	}
}
