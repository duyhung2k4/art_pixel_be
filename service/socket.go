package service

import (
	"app/config"
	"app/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type socketService struct {
	psql *gorm.DB
	rdb  *redis.Client
}

type SocketService interface {
	AddFaceEncoding(auth string) ([]model.Face, error)
}

func (s *socketService) AddFaceEncoding(auth string) ([]model.Face, error) {
	var newListFaceEncoding []model.Face

	// Get images in file add model
	path := fmt.Sprintf("file_add_model/%s", auth)
	cmd := exec.Command("python3", "python_code/face_encoding.py", path)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var faceEncoding [][]float64
	if err := json.Unmarshal(output, &faceEncoding); err != nil {
		return nil, err
	}

	// Get profile redis
	profileRedis, err := s.rdb.Get(context.Background(), auth).Result()
	if err != nil {
		return nil, err
	}

	// Convert to profile struct
	var profile *model.Profile
	if err := json.Unmarshal([]byte(profileRedis), &profile); err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, errors.New("profile null")
	}

	// Create list Face
	for _, data := range faceEncoding {
		newListFaceEncoding = append(newListFaceEncoding, model.Face{
			ProfileId:    profile.ID,
			FaceEncoding: data,
		})
	}

	if err := s.psql.Create(&newListFaceEncoding).Error; err != nil {
		return nil, err
	}

	if err := s.psql.Model(&model.Profile{}).
		Where("id = ?", profile.ID).
		Updates(&model.Profile{Active: true}).
		Error; err != nil {
		return nil, err
	}

	return newListFaceEncoding, nil
}

func NewSocketService() SocketService {
	return &socketService{
		psql: config.GetPsql(),
		rdb:  config.GetRedisClient(),
	}
}
