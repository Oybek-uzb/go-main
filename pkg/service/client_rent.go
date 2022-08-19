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

func (c *ClientRentService) GetMyCarParkByCompanyId(userId, companyId int, inDiscount bool) ([]models.Car, error) {
	cars, err := c.repo.GetMyCarParkByCompanyId(userId, companyId, inDiscount)
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

func (c *ClientRentService) GetMyCompanyById(userId, companyId int) (models.CarCompany, error) {
	company, err := c.repo.GetMyCompanyById(userId, companyId)
	if err != nil {
		return models.CarCompany{}, err
	}

	if company.Photo != nil {
		company.Photo = utils.GetFileUrl(strings.Split(*company.Photo, "/"))
	}

	return company, nil
}

func (c *ClientRentService) GetMyCompaniesList(userId int) ([]models.CarCompany, error) {
	myCompanies, err := c.repo.GetMyCompaniesList(userId)
	if err != nil {
		return nil, err
	}

	for i, company := range myCompanies {
		if company.Photo != nil {
			myCompanies[i].Photo = utils.GetFileUrl(strings.Split(*company.Photo, "/"))
		}
	}

	return myCompanies, nil
}

func (c *ClientRentService) GetCarByCompanyIdCarId(companyId, carId, langId int) (models.Car, error) {
	car, err := c.repo.GetCarByCompanyIdCarId(companyId, carId, langId)
	if err != nil {
		return models.Car{}, err
	}

	if car.Photo != nil {
		car.Photo = utils.GetFileUrl(strings.Split(*car.Photo, "/"))
	}

	return car, nil
}

func (c *ClientRentService) GetCarsByCompanyId(companyId int) (models.CarCompanyDetails, error) {
	carCompany, err := c.repo.GetCarsByCompanyId(companyId)
	if err != nil {
		return models.CarCompanyDetails{}, err
	}

	if carCompany.Photo != nil {
		carCompany.Photo = utils.GetFileUrl(strings.Split(*carCompany.Photo, "/"))
	}

	for i, car := range carCompany.Cars {
		if car.Photo != nil {
			carCompany.Cars[i].Photo = utils.GetFileUrl(strings.Split(*car.Photo, "/"))
		}
	}

	return carCompany, nil
}

func (c *ClientRentService) GetCompaniesList() ([]models.CarCompany, error) {
	carCompanies, err := c.repo.GetCompaniesList()
	if err != nil {
		return nil, err
	}

	for i, company := range carCompanies {
		if company.Photo != nil {
			carCompanies[i].Photo = utils.GetFileUrl(strings.Split(*company.Photo, "/"))
		}
	}

	return carCompanies, nil
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

func (c *ClientRentService) GetCategoriesList(langId int) ([]models.CarCategory, error) {
	carCategories, err := c.repo.GetCategoriesList(langId)
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
