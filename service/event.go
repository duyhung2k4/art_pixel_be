package service

import (
	"app/config"
	"app/dto/request"
	"app/model"

	"gorm.io/gorm"
)

type eventService struct {
	psql *gorm.DB
}

type EventService interface {
	CreateEvent(payload request.CreateEventReq) (*model.Event, error)
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

	return newEvent, nil
}

func NewEventService() EventService {
	return &eventService{
		psql: config.GetPsql(),
	}
}
