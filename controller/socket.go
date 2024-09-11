package controller

import (
	"app/dto/request"
	"app/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/google/uuid"
)

type socketController struct{}

type SocketController interface {
	SendFileAuthFace(w http.ResponseWriter, r *http.Request, payload request.SendFileAuthFaceReq) string
}

func (c *socketController) SendFileAuthFace(w http.ResponseWriter, r *http.Request, payload request.SendFileAuthFaceReq) string {
	imgData, err := base64.StdEncoding.DecodeString(payload.Data)
	fileName := uuid.New().String()
	if err != nil {
		return err.Error()
	}

	// Check num image for train
	countFileFolder, errCountFileFolder := utils.CheckNumFolder("file_add_model")
	if errCountFileFolder != nil {
		return errCountFileFolder.Error()
	}
	if countFileFolder == 10 {
		return "done"
	}

	pathPending := fmt.Sprintf("pending_file/%s.png", fileName)
	filePending, err := os.Create(pathPending)
	if err != nil {
		return err.Error()
	}
	_, err = filePending.Write(imgData)
	if err != nil {
		return err.Error()
	}

	// Check face
	cmdCheckFace := exec.Command("python3", "python_code/check_face.py", pathPending)
	outputCheckFace, errCheckFace := cmdCheckFace.Output()
	if errCheckFace != nil {
		return errCheckFace.Error()
	}
	var resultCheckFace bool
	if err := json.Unmarshal(outputCheckFace, &resultCheckFace); err != nil {
		return err.Error()
	}
	if !resultCheckFace {
		if err := os.Remove(pathPending); err != nil {
			return err.Error()
		}

		return "image not a face!"
	}

	// Add data model
	pathAddModel := fmt.Sprintf("file_add_model/%s.png", fileName)
	fileAddModel, err := os.Create(pathAddModel)
	if err != nil {
		return err.Error()
	}

	_, err = fileAddModel.Write(imgData)
	if err != nil {
		return err.Error()
	}

	return "not enough data"
}

func NewSocketController() SocketController {
	return &socketController{}
}
