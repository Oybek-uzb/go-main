package service

import (
	"abir/models"
	"abir/pkg/repository"
	"abir/pkg/storage"
	"abir/pkg/utils"
	"context"
	"strings"

	"github.com/streadway/amqp"
)

type ClientRentService struct {
	repo        repository.RentCars
	ch          *amqp.Channel
	fileStorage storage.Storage
}

func (c *ClientRentService) DeleteMyCar(userId, carId, myCompanyId int) (int, error) {
	carId, err := c.repo.DeleteMyCar(userId, carId, myCompanyId)
	if err != nil {
		return 0, err
	}

	return carId, nil
}

func (c *ClientRentService) PostRentCarByCarId(userId, carId int, rentCarDetails models.RentCarDetails) (int, error) {
	carId, err := c.repo.PostRentCarByCarId(userId, carId, rentCarDetails)
	if err != nil {
		return 0, err
	}

	return carId, nil
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

func (c *ClientRentService) PostMyCompany(ctx context.Context, userId int, company models.RentMyCompanyCreate) (int, error) {
	if company.Photo != nil && *company.Photo != "" {
		fileName, err := utils.GenerateFileName()
		if err != nil {
			return 0, err
		}
		uploadedFileName, err := c.fileStorage.Upload(ctx, storage.UploadInput{
			File:   *company.Photo,
			Name:   fileName,
			Folder: "companies",
		})
		if err != nil {
			return 0, err
		}
		company.Photo = &uploadedFileName
	}

	companyId, err := c.repo.PostMyCompany(userId, company)
	if err != nil {
		return 0, err
	}

	return companyId, nil
}

func (c *ClientRentService) PostMyCar(ctx context.Context, userId, companyId int, car models.CarCreate) (int, error) {
	if car.Photo != nil && *car.Photo != "" {
		fileName, err := utils.GenerateFileName()
		if err != nil {
			return 0, err
		}
		uploadedFileName, err := c.fileStorage.Upload(ctx, storage.UploadInput{
			File:   *car.Photo,
			Name:   fileName,
			Folder: "companies",
		})
		if err != nil {
			return 0, err
		}
		car.Photo = &uploadedFileName
	}

	carId, err := c.repo.PostMyCar(userId, companyId, car)
	if err != nil {
		return 0, err
	}

	return carId, nil
}

func (c *ClientRentService) PutMyCar(ctx context.Context, userId, carId, companyId int, car models.CarCreate) (int, error) {
	if car.Photo != nil && *car.Photo != "" {
		fileName, err := utils.GenerateFileName()
		if err != nil {
			return 0, err
		}
		uploadedFileName, err := c.fileStorage.Upload(ctx, storage.UploadInput{
			File:   *car.Photo,
			Name:   fileName,
			Folder: "companies",
		})
		if err != nil {
			return 0, err
		}
		car.Photo = &uploadedFileName
	}

	carId, err := c.repo.PutMyCar(userId, carId, companyId, car)
	if err != nil {
		return 0, err
	}

	return carId, nil
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

func (c *ClientRentService) GetMyCarsForEvents(userId, langId int) ([]models.MyCarForEvents, error) {
	myCars, err := c.repo.GetMyCarsForEvents(userId, langId)
	if err != nil {
		return nil, err
	}

	for i, company := range myCars {
		if company.Photo != nil {
			myCars[i].Photo = utils.GetFileUrl(strings.Split(*company.Photo, "/"))
		}
	}

	return myCars, nil
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

func (c *ClientRentService) GetCategoriesForEvents(langId int) ([]models.CarCategory, error) {
	carCategories, err := c.repo.GetCategoriesForEvents(langId)
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

func NewClientRentService(repo repository.RentCars, ch *amqp.Channel, s storage.Storage) *ClientRentService {
	return &ClientRentService{
		repo:        repo,
		ch:          ch,
		fileStorage: s,
	}
}
