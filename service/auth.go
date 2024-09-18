package service

import (
	"app/config"
	queuepayload "app/dto/queue_payload"
	"app/dto/request"
	"app/model"
	"app/utils"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type authService struct {
	psql  *gorm.DB
	redis *redis.Client
}

type AuthService interface {
	CheckExistProfile(registerReq request.RegisterReq) (bool, error)
	CreateProfilePending(registerReq request.RegisterReq) (*model.Profile, error)
	CheckFace(payload queuepayload.SendFileAuthMess) (string, error)

	saveFileAuth(auth string) error
}

func (s *authService) CheckExistProfile(registerReq request.RegisterReq) (bool, error) {
	var profile *model.Profile

	if err := s.psql.
		Model(&model.Profile{}).
		Where("email = ? AND active = ?", registerReq.Email, true).
		First(&profile).Error; err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if profile.ID != 0 {
		return true, nil
	}

	return false, nil
}

func (s *authService) CreateProfilePending(registerReq request.RegisterReq) (*model.Profile, error) {
	var newProfile = model.Profile{
		FirstName: registerReq.FirstName,
		LastName:  registerReq.LastName,
		Email:     registerReq.Email,
		Active:    false,
	}

	if err := s.psql.Model(&model.Profile{}).Create(&newProfile).Error; err != nil {
		return nil, err
	}

	return &newProfile, nil
}

func (s *authService) CheckFace(payload queuepayload.SendFileAuthMess) (string, error) {
	base64Data := payload.Data
	imgData, err := base64.StdEncoding.DecodeString(base64Data[strings.IndexByte(base64Data, ',')+1:])
	fileName := uuid.New().String()
	if err != nil {
		log.Println(err)
		return "", err
	}

	// Check num image for train
	pathCheckNumFolder := fmt.Sprintf("file_add_model/%s", payload.Uuid)
	countFileFolder, err := utils.CheckNumFolder(pathCheckNumFolder)
	if err != nil {
		return "", err
	}
	if countFileFolder >= 10 {
		if err := s.saveFileAuth(payload.Uuid); err != nil {
			return "", err
		}

		return "done", nil
	}

	pathPending := fmt.Sprintf("pending_file/%s/%s.png", payload.Uuid, fileName)
	filePending, err := os.Create(pathPending)
	if err != nil {
		return "", err
	}
	_, err = filePending.Write(imgData)
	if err != nil {
		return "", err
	}

	// Check face
	cmdCheckFace := exec.Command("python3", "python_code/check_face.py", pathPending)
	outputCheckFace, err := cmdCheckFace.Output()
	if err != nil {
		return "", err
	}
	var resultCheckFace bool
	if err := json.Unmarshal(outputCheckFace, &resultCheckFace); err != nil {
		return "", err
	}
	if !resultCheckFace {
		if err := os.Remove(pathPending); err != nil {
			return "", err
		}

		return "image not a face!", nil
	}

	// Add data model
	pathAddModel := fmt.Sprintf("file_add_model/%s/%s.png", payload.Uuid, fileName)
	fileAddModel, err := os.Create(pathAddModel)
	if err != nil {
		return "", err
	}

	_, err = fileAddModel.Write(imgData)
	if err != nil {
		return "", err
	}

	return "not enough data", nil
}

func (s *authService) saveFileAuth(auth string) error {
	// convert data file add model
	pathFileAddModel := fmt.Sprintf("file_add_model/%s", auth)
	cmdFaceEndcoding := exec.Command("python3", "python_code/face_encoding.py", pathFileAddModel)
	outputCheckFace, err := cmdFaceEndcoding.Output()

	if err != nil {
		return err
	}

	var listImages [][]float64
	if err := json.Unmarshal(outputCheckFace, &listImages); err != nil {
		return err
	}

	// get profile auth in redis
	profileJson, err := s.redis.Get(context.Background(), auth).Result()
	if err != nil {
		return err
	}
	var profile model.Profile
	if err := json.Unmarshal([]byte(profileJson), &profile); err != nil {
		return err
	}

	// add faces
	var faces []model.Face
	for _, img := range listImages {
		faces = append(faces, model.Face{
			ProfileId:    profile.ID,
			FaceEncoding: img,
		})
	}

	if err := s.psql.Model(&model.Face{}).Create(&faces).Error; err != nil {
		return err
	}

	if err := s.psql.
		Model(&model.Profile{}).
		Where("id = ?", profile.ID).
		Updates(&model.Profile{Active: true}).Error; err != nil {
		return err
	}

	// delete pending file
	pendingPath := fmt.Sprintf("pending_file/%s", auth)
	if err := os.RemoveAll(pendingPath); err != nil {
		return err
	}
	addModelPath := fmt.Sprintf("file_add_model/%s", auth)
	if err := os.RemoveAll(addModelPath); err != nil {
		return err
	}

	return nil
}

func NewAuthService() AuthService {
	return &authService{
		redis: config.GetRedisClient(),
		psql:  config.GetPsql(),
	}
}
