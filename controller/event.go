package controller

import "net/http"

type eventController struct{}

type EventController interface {
	CreateEvent(w http.ResponseWriter, r *http.Request)
}

func (c *eventController) CreateEvent(w http.ResponseWriter, r *http.Request) {

}

func NewEventController() EventController {
	return &eventController{}
}
