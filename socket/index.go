package socket

import (
	"app/constant"
	"app/controller"
	"app/dto/request"
	"app/utils"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
)

func ServerSocker() http.Handler {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	app := chi.NewRouter()

	// A good base middleware stack
	app.Use(middleware.RequestID)
	app.Use(middleware.RealIP)
	app.Use(middleware.Logger)
	app.Use(middleware.Recoverer)

	socketController := controller.NewSocketController()

	app.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Print("upgrade:", err)
			return
		}

		defer c.Close()

		for {
			var messRes string

			mt, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}

			var requestData request.SocketRequest

			errConvert := json.Unmarshal(message, &requestData)
			if errConvert != nil {
				messRes = errConvert.Error()
			}

			jsonData, errJsonData := json.Marshal(requestData.Data)
			if errJsonData != nil {
				messRes = errJsonData.Error()
			}

			if messRes != "" {
				err = c.WriteMessage(mt, utils.ConvertToByte(messRes))
				if err != nil {
					log.Println("write:", err)
					break
				}
				continue
			}

			switch requestData.Type {
			case constant.SEND_FILE_AUTH_FACE:
				var payload request.SendFileAuthFaceReq
				errConvert := json.Unmarshal(jsonData, &payload)

				if errConvert != nil {
					messRes = errConvert.Error()
				} else {
					messRes = socketController.SendFileAuthFace(w, r, payload)
				}
			default:
			}

			err = c.WriteMessage(mt, utils.ConvertToByte(messRes))
			if err != nil {
				log.Println("write:", err)
				break
			}
		}

		log.Printf("Disconnect")
	})

	log.Printf("Socket art-pixel starting success!")

	return app
}
