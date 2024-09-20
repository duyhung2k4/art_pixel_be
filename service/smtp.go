package service

import (
	"app/config"
	"app/model"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/smtp"
	"time"

	"github.com/redis/go-redis/v9"
)

type smtpService struct {
	redisClient *redis.Client
	authSmtp    smtp.Auth
	smtpHost    string
	smtpPort    string
}

type SmtpService interface {
	SendCodeAcceptRegister(auth string) error
}

func (s *smtpService) SendCodeAcceptRegister(auth string) error {
	var profile model.Profile

	profileJson, err := s.redisClient.Get(context.Background(), auth).Result()
	if err != nil {
		return err
	}
	if err = json.Unmarshal([]byte(profileJson), &profile); err != nil {
		return err
	}

	code := fmt.Sprintf("%d", rand.Intn(900000)+100000)

	to := []string{profile.Email}
	msg := []byte(code)

	err = smtp.SendMail(s.smtpHost+":"+s.smtpPort, s.authSmtp, profile.Email, to, msg)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("code_%s", auth)
	err = s.redisClient.SetNX(context.Background(), key, code, 3*60*time.Second).Err()
	if err != nil {
		return err
	}

	return nil
}

func NewSmtpService() SmtpService {
	return &smtpService{
		redisClient: config.GetRedisClient(),
		authSmtp:    config.GetAuthSmtp(),
		smtpHost:    config.GetSmtpHost(),
		smtpPort:    config.GetSmtpPort(),
	}
}
