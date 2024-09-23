package request

import "time"

type CreateEventReq struct {
	StartTime  time.Time `json:"startTime"`
	FinishTime time.Time `json:"finishTime"`
	SizeX      int       `json:"sizeX"`
	SizeY      int       `json:"sizeY"`
}

type DrawPixelReq struct {
	EventId   uint `json:"eventId"`
	X         int  `json:"x"`
	Y         int  `json:"y"`
	ProfileId uint `json:"profileId"`
}
