package service

import (
	"abir/models"
	"abir/pkg/repository"
	"abir/pkg/utils"
)

type DriverSettingsService struct {
	repo      repository.DriverSettings
	fcmClient *utils.FCMClient
}

func NewDriverSettingsService(repo repository.DriverSettings, fcmClient *utils.FCMClient) *DriverSettingsService {
	return &DriverSettingsService{repo: repo, fcmClient: fcmClient}
}

func (s *DriverSettingsService) GetTariffs(userId, langId int) ([]models.DriverTariffs, error) {
	return s.repo.GetTariffs(userId, langId)
}

func (s *DriverSettingsService) TariffsEnable(userId, tariffId int, isActive bool) error {
	return s.repo.TariffsEnable(userId, tariffId, isActive)
}

func (s *DriverSettingsService) SetOnline(userId int, isActive int) error {
	return s.repo.SetOnline(userId, isActive)
}
