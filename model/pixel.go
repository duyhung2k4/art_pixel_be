package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Pixel struct {
	Id              primitive.ObjectID `json:"_id" bson:"_id"`
	EventId         uint               `json:"event_id" bson:"event_id"`
	X               int                `json:"x" bson:"x"`
	Y               int                `json:"y" bson:"y"`
	ProfileIdUpDate uint               `json:"profileIdUpDate" bson:"profileIdUpDate"`
}
