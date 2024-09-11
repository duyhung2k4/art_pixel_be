package service

import (
	"app/config"
	"app/model"

	"gorm.io/gorm"
)

type socketService struct {
	db *gorm.DB
}

type SocketService interface {
	AddFaceEncoding(profileId uint, faceEncoding [][]float64) ([]model.Face, error)
}

func (s *socketService) AddFaceEncoding(profileId uint, faceEncoding [][]float64) ([]model.Face, error) {
	var newListFaceEncoding []model.Face

	for _, data := range faceEncoding {
		newListFaceEncoding = append(newListFaceEncoding, model.Face{
			ProfileId:    profileId,
			FaceEncoding: data,
		})

		if err := s.db.Model(&model.Face{}).Create(&newListFaceEncoding).Error; err != nil {
			return []model.Face{}, err
		}
	}

	return newListFaceEncoding, nil
}

func NewSocketService() SocketService {
	return &socketService{
		db: config.GetPsql(),
	}
}
