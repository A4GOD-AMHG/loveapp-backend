package services

import (
	"errors"
	"strings"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/repository"
)

// DeviceService define la lógica de negocio para dispositivos push.
type DeviceService interface {
	RegisterPushToken(userID int64, req *models.RegisterPushTokenRequest) error
	DeletePushToken(userID int64, pushToken string) error
}

type deviceService struct {
	repo repository.DevicePushTokenRepository
}

// NewDeviceService crea un nuevo servicio para tokens push.
func NewDeviceService(repo repository.DevicePushTokenRepository) DeviceService {
	return &deviceService{repo: repo}
}

func (s *deviceService) RegisterPushToken(userID int64, req *models.RegisterPushTokenRequest) error {
	platform := strings.TrimSpace(strings.ToLower(req.Platform))
	deviceName := strings.TrimSpace(req.DeviceName)
	pushToken := strings.TrimSpace(req.PushToken)

	if userID == 0 || platform == "" || deviceName == "" || pushToken == "" {
		return errors.New("invalid push token payload")
	}

	return s.repo.Upsert(&models.DevicePushToken{
		UserID:     userID,
		Platform:   platform,
		PushToken:  pushToken,
		DeviceName: deviceName,
	})
}

func (s *deviceService) DeletePushToken(userID int64, pushToken string) error {
	pushToken = strings.TrimSpace(pushToken)
	if userID == 0 || pushToken == "" {
		return errors.New("invalid push token")
	}

	return s.repo.DeleteByToken(userID, pushToken)
}
