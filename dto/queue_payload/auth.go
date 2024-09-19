package queuepayload

type SendFileAuthMess struct {
	Uuid string
	Data string
}

type FaceAuth struct {
	Uuid     string
	FilePath string
}
