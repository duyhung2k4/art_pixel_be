package service

import (
	"app/config"
	"app/constant"
	queuepayload "app/dto/queue_payload"
	"app/dto/request"
	"app/model"
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type eventService struct {
	mutex          *sync.Mutex
	psql           *gorm.DB
	mongodb        *mongo.Database
	mapSocketEvent map[string]map[string]*websocket.Conn
}

type EventService interface {
	GetAllEvent() ([]model.Event, error)
	CreateEvent(payload request.CreateEventReq) (*model.Event, error)
	DrawPixel(payload queuepayload.DrawPixel, profileId uint) (*model.Pixel, error)
}

func (s *eventService) GetAllEvent() ([]model.Event, error) {
	var events []model.Event

	if err := s.psql.Model(&model.Event{}).Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}

func (s *eventService) CreateEvent(payload request.CreateEventReq) (*model.Event, error) {
	var newEvent *model.Event = &model.Event{
		StartTime:  payload.StartTime,
		FinishTime: payload.FinishTime,
		SizeX:      payload.SizeX,
		SizeY:      payload.SizeY,
	}

	if err := s.psql.Model(&model.Event{}).Create(&newEvent).Error; err != nil {
		return nil, err
	}

	listPixelInteface := []interface{}{}

	for x := 0; x < newEvent.SizeX; x++ {
		for y := 0; y < newEvent.SizeY; y++ {
			pixel := PixelInsert{
				X:               x,
				Y:               y,
				EventId:         newEvent.ID,
				Rgb:             nil,
				ProfileIdUpDate: 0,
			}
			listPixelInteface = append(listPixelInteface, pixel)
		}
	}

	_, err := s.mongodb.Collection(string(constant.PIXEL)).InsertMany(context.TODO(), listPixelInteface)
	if err != nil {
		return nil, err
	}

	return newEvent, nil
}

func (s *eventService) DrawPixel(payload queuepayload.DrawPixel, profileId uint) (*model.Pixel, error) {
	filter := bson.M{
		"x":       payload.Data.X,
		"y":       payload.Data.Y,
		"eventId": payload.Data.EventId,
	}

	update := bson.M{
		"$set": bson.M{
			"profileIdUpDate": profileId,
			"rgb":             payload.Data.Rgb,
		},
	}

	result, err := s.mongodb.Collection(string(constant.PIXEL)).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	if result.MatchedCount == 0 {
		return nil, errors.New("pixel not found")
	}

	eventIdString := fmt.Sprint(payload.EventId)

	for _, c := range s.mapSocketEvent[eventIdString] {
		s.mutex.Lock()
		c.WriteMessage(websocket.TextMessage, []byte("cc"))
		s.mutex.Unlock()
	}

	return nil, nil
}

func NewEventService() EventService {
	return &eventService{
		mutex:          new(sync.Mutex),
		psql:           config.GetPsql(),
		mongodb:        config.GetMongoDB(),
		mapSocketEvent: config.GetSocketEvent(),
	}
}

type PixelInsert struct {
	EventId         uint    `json:"eventId" bson:"eventId"`
	X               int     `json:"x" bson:"x"`
	Y               int     `json:"y" bson:"y"`
	Rgb             *string `json:"rgb" bson:"rgb"`
	ProfileIdUpDate uint    `json:"profileIdUpDate" bson:"profileIdUpDate"`
}
