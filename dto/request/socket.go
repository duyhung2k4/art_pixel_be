package request

import (
	"app/constant"
)

type SocketRequest struct {
	Type constant.SOCKET_MESS   `json:"type"`
	Data map[string]interface{} `json:"data"`
}

type SendFileAuthFaceReq struct {
	TypeFile string `json:"typeFile"`
	Data     string `json:"data"`
}
