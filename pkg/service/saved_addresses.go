package service

import (
	"abir/models"
	"abir/pkg/repository"
)

type SavedAddressesService struct {
	repo repository.SavedAddresses
}

func NewSavedAddressesService(repo repository.SavedAddresses) *SavedAddressesService {
	return &SavedAddressesService{repo: repo}
}

func (s *SavedAddressesService) Get(userId int) ([]models.SavedAddresses, error) {
	return s.repo.Get(userId)
}
func (s *SavedAddressesService) Store(address models.SavedAddresses, userId int) error {
	return s.repo.Store(address, userId)
}
func (s *SavedAddressesService) Update(address models.SavedAddresses, addressId, userId int) error {
	return s.repo.Update(address, addressId, userId)
}
func (s *SavedAddressesService) Delete(addressId, userId int) error {
	return s.repo.Delete(addressId, userId)
}

