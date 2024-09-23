package queuepayload

import "app/dto/request"

type DrawPixel struct {
	AccessToken string               `json:"accessToken"`
	Data        request.DrawPixelReq `json:"data"`
}
