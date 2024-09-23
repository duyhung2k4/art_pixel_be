package constant

type QUEUE string

const (
	SEND_FILE_AUTH_QUEUE QUEUE = "send_file_auth_queue"
	FACE_AUTH_QUEUE      QUEUE = "face_auth_queue"
	DRAW_PIXEL_QUEUE     QUEUE = "draw_pixel_queue"
)
