package service

import (
	"app/config"
	"app/model"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

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
	path := fmt.Sprintf("file/file_add_model/%s", auth)

	// Tạo dữ liệu JSON để gửi đến API
	data := map[string]string{
		"directory_path": path,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Gửi yêu cầu POST đến API Flask
	resp, err := http.Post("http://localhost:5000/face_encoding", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Kiểm tra mã trạng thái HTTP
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to call API, status code: %d", resp.StatusCode)
	}

	// Đọc phản hồi từ API
	var response struct {
		Result        string      `json:"result"`
		FaceEncodings [][]float64 `json:"face_encodings"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	// Kiểm tra kết quả từ API
	if response.Result != "success" {
		return nil, fmt.Errorf("API error: %s", response.Result)
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
	for _, data := range response.FaceEncodings {
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
