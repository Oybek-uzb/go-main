package service

import (
	"abir/models"
	"abir/pkg/repository"
	"abir/pkg/utils"
	"github.com/streadway/amqp"
	"strings"
)

type ClientRentService struct {
	repo repository.RentCars
	ch   *amqp.Channel
}

func (c *ClientRentService) GetCarByCategoryIdCarId(categoryId, carId, langId int) (models.Car, error) {
	car, err := c.repo.GetCarByCategoryIdCarId(categoryId, carId, langId)
	if err != nil {
		return models.Car{}, err
	}

	if car.Photo != nil {
		car.Photo = utils.GetFileUrl(strings.Split(*car.Photo, "/"))
	}

	return car, nil
}

func (c *ClientRentService) GetCarsByCategoryId(categoryId int) ([]models.CarByCategoryId, error) {
	cars, err := c.repo.GetCarsByCategoryId(categoryId)
	if err != nil {
		return nil, err
	}

	for i, car := range cars {
		if car.Photo != nil {
			cars[i].Photo = utils.GetFileUrl(strings.Split(*car.Photo, "/"))
		}
	}

	return cars, nil
}

func (c *ClientRentService) GetCategoriesList() ([]models.CarCategory, error) {
	carCategories, err := c.repo.GetCategoriesList()
	if err != nil {
		return nil, err
	}

	for i, category := range carCategories {
		if category.Photo != nil {
			carCategories[i].Photo = utils.GetFileUrl(strings.Split(*category.Photo, "/"))
		}
	}

	return carCategories, nil
}

func NewClientRentService(repo repository.RentCars, ch *amqp.Channel) *ClientRentService {
	return &ClientRentService{
		repo: repo,
		ch:   ch,
	}
}
