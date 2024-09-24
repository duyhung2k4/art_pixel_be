package queuepayload

import "app/dto/request"

type DrawPixel struct {
	AccessToken string               `json:"accessToken"`
	ProfileId   string               `json:"profileId"`
	EventId     uint                 `json:"eventId"`
	Data        request.DrawPixelReq `json:"data"`
}
