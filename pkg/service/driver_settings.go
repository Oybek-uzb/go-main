package service

import (
	"abir/models"
	"abir/pkg/repository"
)

type DriverSettingsService struct {
	repo repository.DriverSettings
}

func NewDriverSettingsService(repo repository.DriverSettings) *DriverSettingsService {
	return &DriverSettingsService{repo: repo}
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
