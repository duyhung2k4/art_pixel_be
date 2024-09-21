package service

import (
	"app/config"
	queuepayload "app/dto/queue_payload"
	"app/dto/request"
	"app/model"
	"app/utils"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type authService struct {
	psql        *gorm.DB
	redis       *redis.Client
	smtpService SmtpService
}

type AuthService interface {
	CheckExistProfile(registerReq request.RegisterReq) (bool, error)
	CreateProfilePending(registerReq request.RegisterReq) (*model.Profile, error)
	CheckFace(payload queuepayload.SendFileAuthMess) (string, error)
	CreateFileAuthFace(data request.AuthFaceReq) (string, error)
	AuthFace(payload queuepayload.FaceAuth) (bool, error)
	ActiveProfile(auth string) error
	SaveFileAuth(auth string) error
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
	if err != nil {
		log.Println(err)
		return "", err
	}

	fileName := uuid.New().String()

	// Check num image for train
	pathCheckNumFolder := fmt.Sprintf("file/file_add_model/%s", payload.Uuid)
	countFileFolder, err := utils.CheckNumFolder(pathCheckNumFolder)
	if err != nil {
		return "", err
	}
	if countFileFolder >= 10 {
		return "done", nil
	}

	// Tạo file tạm thời từ dữ liệu ảnh
	pathPending := fmt.Sprintf("file/pending_file/%s/%s.png", payload.Uuid, fileName)
	filePending, err := os.Create(pathPending)
	if err != nil {
		return "", err
	}
	defer filePending.Close()

	_, err = filePending.Write(imgData)
	if err != nil {
		return "", err
	}

	payloadDetectFace, err := json.Marshal(map[string]interface{}{
		"input_image_path": pathPending,
	})
	if err != nil {
		return "", err
	}

	// Gọi API để kiểm tra khuôn mặt
	resp, err := http.Post("http://localhost:5000/detect_single_face", "application/json", bytes.NewBuffer(payloadDetectFace))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Kiểm tra mã trạng thái HTTP
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to call API, status code: %d", resp.StatusCode)
	}

	// Đọc phản hồi từ API
	var resultCheckFace struct {
		Result bool `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&resultCheckFace); err != nil {
		return "", err
	}

	if !resultCheckFace.Result {
		if err := os.Remove(pathPending); err != nil {
			return "", err
		}
		return "image not a face!", nil
	}

	// Thêm dữ liệu vào mô hình
	pathAddModel := fmt.Sprintf("file/file_add_model/%s/%s.png", payload.Uuid, fileName)
	fileAddModel, err := os.Create(pathAddModel)
	if err != nil {
		return "", err
	}
	defer fileAddModel.Close()

	_, err = fileAddModel.Write(imgData)
	if err != nil {
		return "", err
	}

	return "not enough data", nil
}

func (s *authService) CreateFileAuthFace(data request.AuthFaceReq) (string, error) {
	base64Data := data.Data
	imgData, err := base64.StdEncoding.DecodeString(base64Data[strings.IndexByte(base64Data, ',')+1:])
	fileName := uuid.New().String()

	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("file/auth_face/%s.png", fileName)
	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	_, err = file.Write(imgData)
	if err != nil {
		return "", err
	}

	return path, nil
}

func (s *authService) AuthFace(payload queuepayload.FaceAuth) (bool, error) {
	var faces []model.Face

	// Lấy danh sách khuôn mặt từ cơ sở dữ liệu
	if err := s.psql.Model(&model.Face{}).Find(&faces).Error; err != nil {
		return false, err
	}

	// Tạo dữ liệu JSON để gửi đến API
	data := map[string]interface{}{
		"faces":            faces,
		"input_image_path": payload.FilePath,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return false, err
	}

	// Gửi yêu cầu POST đến API Flask
	resp, err := http.Post("http://localhost:5000/recognize_faces", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Kiểm tra mã trạng thái HTTP
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("failed to call API, status code: %d", resp.StatusCode)
	}

	// Đọc phản hồi từ API
	var response struct {
		Result string `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return false, err
	}

	// Xử lý kết quả từ API
	if response.Result == "-1" {
		return false, nil // Không tìm thấy khuôn mặt phù hợp
	}

	profileId, err := strconv.Atoi(response.Result)
	if err != nil {
		return false, err
	}

	return profileId >= 0, nil
}

func (s *authService) ActiveProfile(auth string) error {
	var profile model.Profile
	profileJson, err := s.redis.Get(context.Background(), auth).Result()

	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(profileJson), &profile); err != nil {
		return err
	}

	if err := s.psql.
		Model(&model.Profile{}).
		Where("id = ?", profile.ID).
		Updates(&model.Profile{Active: true}).
		Error; err != nil {
		return err
	}

	return nil
}

func (s *authService) SaveFileAuth(auth string) error {
	// Tạo đường dẫn đến thư mục chứa ảnh
	pathFileAddModel := fmt.Sprintf("file/file_add_model/%s", auth)

	// Tạo dữ liệu JSON để gửi đến API
	data := map[string]interface{}{
		"directory_path": pathFileAddModel,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Gửi yêu cầu POST đến API Flask
	resp, err := http.Post("http://localhost:5000/face_encoding", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Kiểm tra mã trạng thái HTTP
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to call API, status code: %d", resp.StatusCode)
	}

	// Đọc phản hồi từ API
	var response struct {
		Result        string      `json:"result"`
		FaceEncodings [][]float64 `json:"face_encodings"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	if response.Result != "success" {
		return fmt.Errorf("API error: %s", response.Result)
	}

	// Lấy thông tin profile từ Redis
	profileJson, err := s.redis.Get(context.Background(), auth).Result()
	if err != nil {
		return err
	}
	var profile model.Profile
	if err := json.Unmarshal([]byte(profileJson), &profile); err != nil {
		return err
	}

	// Thêm khuôn mặt vào danh sách
	var faces []model.Face
	for _, img := range response.FaceEncodings {
		faces = append(faces, model.Face{
			ProfileId:    profile.ID,
			FaceEncoding: img,
		})
	}

	if err := s.psql.Model(&model.Face{}).Create(&faces).Error; err != nil {
		return err
	}

	// Xóa thư mục tạm
	pendingPath := fmt.Sprintf("file/pending_file/%s", auth)
	if err := os.RemoveAll(pendingPath); err != nil {
		return err
	}
	addModelPath := fmt.Sprintf("file/file_add_model/%s", auth)
	if err := os.RemoveAll(addModelPath); err != nil {
		return err
	}

	return nil
}

func NewAuthService() AuthService {
	return &authService{
		redis:       config.GetRedisClient(),
		psql:        config.GetPsql(),
		smtpService: NewSmtpService(),
	}
}
