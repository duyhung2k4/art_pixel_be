package socket

import (
	"app/constant"
	"app/controller"
	"app/dto/request"
	"app/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func AuthSocket(w http.ResponseWriter, r *http.Request, upgrader websocket.Upgrader) {
	socketController := controller.NewSocketController()
	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		var messRes string

		// Convert data
		var requestData request.SocketRequest
		err = json.Unmarshal(message, &requestData)
		if err != nil {
			err = c.WriteMessage(mt, utils.ConvertToByte(err.Error()))
			if err != nil {
				log.Println("write:", err)
				break
			}
			continue
		}

		jsonData, err := json.Marshal(requestData.Data)
		if err != nil {
			err = c.WriteMessage(mt, utils.ConvertToByte(err.Error()))
			if err != nil {
				log.Println("write:", err)
				break
			}
			continue
		}

		// Handle message
		switch requestData.Type {
		case constant.SEND_FILE_AUTH_FACE:
			messRes = socketController.SendFileAuthFace(w, r, jsonData, requestData.Auth)
		case constant.SEND_MESS:
			messRes = fmt.Sprintf("rep: %d", time.Now().Unix())
		default:
			messRes = "null"
			log.Println(r.Cookies())
		}

		err = c.WriteMessage(mt, utils.ConvertToByte(messRes))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}

	log.Printf("Disconnect")
}
