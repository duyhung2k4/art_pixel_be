package response

type SocketErrorRes struct {
	Data  interface{} `json:"data"`
	Error error       `json:"error"`
}
