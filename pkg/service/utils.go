package service

import (
	"abir/models"
	"abir/pkg/repository"
	"context"
)

type UtilsService struct {
	repo repository.Utils
}

func NewUtilsService(repo repository.Utils) *UtilsService {
	return &UtilsService{repo: repo}
}

func (s *UtilsService) GetColors(langId int) ([]models.Color, error) {
	return s.repo.GetColors(langId)
}
func (s *UtilsService) GetCarMarkas() ([]models.CarMarka, error) {
	return s.repo.GetCarMarkas()
}
func (s *UtilsService) GetCarModels(id int) ([]models.CarModel, error) {
	return s.repo.GetCarModels(id)
}
func (s *UtilsService) Test(ctx context.Context) error {
	//key := fmt.Sprintf("Key-%d", 1)
	//msg := kafka.Message{
	//	Key:   []byte(key),
	//	Value: []byte(fmt.Sprint(uuid.New())),
	//}
	//err := r.orderWriter.WriteMessages(ctx, msg)
	//if err != nil {
	//	return err
	//} else {
	//	logrus.Println("produced", key)
	//}
	return nil
}

func (s *UtilsService) GetRegions(langId int) ([]models.Region, error) {
	return s.repo.GetRegions(langId)
}
func (s *UtilsService) GetDistricts(langId int, id int) ([]models.District, error) {
	return s.repo.GetDistricts(langId, id)
}
func (s *UtilsService) GetDistrictsArr(langId int) (map[int]string, error) {
	return s.repo.GetDistrictsArr(langId)
}
func (s *UtilsService) GetTariffs(langId int) (map[string]string, error) {
	return s.repo.GetTariffs(langId)
}

func (s *UtilsService) DriverCancelOrderOptions(langId int) ([]models.DriverCancelOrderOptions, error) {
	return s.repo.DriverCancelOrderOptions(langId)
}
func (s *UtilsService) ClientCancelOrderOptions(langId int, optionType string) ([]models.ClientCancelOrderOptions, error) {
	return s.repo.ClientCancelOrderOptions(langId, optionType)
}

func (s *UtilsService) ClientRateOptions(langId, rate int, optionType string) ([]models.ClientRateOptions, error) {
	return s.repo.ClientRateOptions(langId, rate, optionType)
}
