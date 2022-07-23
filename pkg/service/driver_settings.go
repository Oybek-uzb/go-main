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

func (s *DriverSettingsService) GetTariffs(userId, langId int) ([]models.DriverTariffs, error){
	return s.repo.GetTariffs(userId, langId)
}
