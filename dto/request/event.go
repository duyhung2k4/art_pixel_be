package request

import "time"

type CreateEventReq struct {
	StartTime  time.Time `json:"startTime"`
	FinishTime time.Time `json:"finishTime"`
	SizeX      int       `json:"sizeX"`
	SizeY      int       `json:"sizeY"`
}
