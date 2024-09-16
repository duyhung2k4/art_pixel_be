package controller

// import (
// 	"app/dto/request"
// 	"app/service"
// 	"app/utils"
// 	"encoding/base64"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"os/exec"

// 	"github.com/google/uuid"
// )

// type socketController struct {
// 	socketService service.SocketService
// }

// type SocketController interface {
// 	SendFileAuthFace(w http.ResponseWriter, r *http.Request, jsonData []byte, auth string) string
// }

// func (c *socketController) SendFileAuthFace(w http.ResponseWriter, r *http.Request, jsonData []byte, auth string) string {
// 	var payload request.SendFileAuthFaceReq

// 	if err := json.Unmarshal(jsonData, &payload); err != nil {
// 		return err.Error()
// 	}

// 	imgData, err := base64.StdEncoding.DecodeString(payload.Data)
// 	fileName := uuid.New().String()
// 	if err != nil {
// 		return err.Error()
// 	}

// 	// Check num image for train
// 	pathCheckNumFolder := fmt.Sprintf("file_add_model/%s", auth)
// 	countFileFolder, err := utils.CheckNumFolder(pathCheckNumFolder)
// 	if err != nil {
// 		return err.Error()
// 	}
// 	if countFileFolder == 10 {
// 		if _, err := c.socketService.AddFaceEncoding(auth); err != nil {
// 			return err.Error()
// 		}

// 		pendingPath := fmt.Sprintf("pending_file/%s", auth)
// 		if err := os.RemoveAll(pendingPath); err != nil {
// 			return err.Error()
// 		}
// 		addModelPath := fmt.Sprintf("file_add_model/%s", auth)
// 		if err := os.RemoveAll(addModelPath); err != nil {
// 			return err.Error()
// 		}

// 		return "done"
// 	}

// 	pathPending := fmt.Sprintf("pending_file/%s/%s.png", auth, fileName)
// 	filePending, err := os.Create(pathPending)
// 	if err != nil {
// 		return err.Error()
// 	}
// 	_, err = filePending.Write(imgData)
// 	if err != nil {
// 		return err.Error()
// 	}

// 	// Check face
// 	cmdCheckFace := exec.Command("python3", "python_code/check_face.py", pathPending)
// 	outputCheckFace, err := cmdCheckFace.Output()
// 	if err != nil {
// 		return err.Error()
// 	}
// 	var resultCheckFace bool
// 	if err := json.Unmarshal(outputCheckFace, &resultCheckFace); err != nil {
// 		return err.Error()
// 	}
// 	if !resultCheckFace {
// 		if err := os.Remove(pathPending); err != nil {
// 			return err.Error()
// 		}

// 		return "image not a face!"
// 	}

// 	// Add data model
// 	pathAddModel := fmt.Sprintf("file_add_model/%s/%s.png", auth, fileName)
// 	fileAddModel, err := os.Create(pathAddModel)
// 	if err != nil {
// 		return err.Error()
// 	}

// 	_, err = fileAddModel.Write(imgData)
// 	if err != nil {
// 		return err.Error()
// 	}

// 	return "not enough data"
// }

// func NewSocketController() SocketController {
// 	return &socketController{
// 		socketService: service.NewSocketService(),
// 	}
// }
