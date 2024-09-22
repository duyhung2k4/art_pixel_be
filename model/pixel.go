package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Pixel struct {
	Id              primitive.ObjectID `json:"_id"`
	EventId         uint               `json:"eventId"`
	X               int                `json:"x"`
	Y               int                `json:"y"`
	ProfileIdUpDate uint               `json:"profileIdUpDate"`
}
