package service

import (
	"app/config"
	"app/dto/request"
	"app/model"

	"gorm.io/gorm"
)

type authService struct {
	psql *gorm.DB
}

type AuthService interface {
	CheckExistProfile(registerReq request.RegisterReq) (bool, error)
	CreateProfilePending(registerReq request.RegisterReq) (*model.Profile, error)
}

func (s *authService) CheckExistProfile(registerReq request.RegisterReq) (bool, error) {
	var profile *model.Profile

	if err := s.psql.
		Model(&model.Profile{}).
		Where("email = ? AND active = ?", registerReq.Email, true).
		First(&profile).Error; err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if profile != nil {
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

func NewAuthService() AuthService {
	return &authService{
		psql: config.GetPsql(),
	}
}
