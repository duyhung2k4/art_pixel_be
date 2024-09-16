package request

import (
	"app/constant"
)

type SocketRequest struct {
	Type constant.SOCKET_MESS   `json:"type"`
	Auth string                 `json:"auth"`
	Data map[string]interface{} `json:"data"`
}

// type SendFileAuthFaceReq struct {
// 	Data string `json:"data"`
// }
