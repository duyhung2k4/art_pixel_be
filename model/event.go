package model

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model

	StartTime  time.Time `json:"startTime"`
	FinishTime time.Time `json:"finishTime"`
	SizeX      int       `json:"sizeX"`
	SizeY      int       `json:"sizeY"`
}
