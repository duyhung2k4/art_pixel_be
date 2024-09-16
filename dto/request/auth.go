package request

type RegisterReq struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type SendFileAuthFaceReq struct {
	Data string `json:"data"`
}
