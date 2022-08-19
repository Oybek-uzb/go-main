package postgres

import (
	"abir/models"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type RentCarsPostgres struct {
	db   *sqlx.DB
	dash *sqlx.DB
}

func NewRentCarsPostgres(db *sqlx.DB, dash *sqlx.DB) *RentCarsPostgres {
	return &RentCarsPostgres{db: db, dash: dash}
}

func (r *RentCarsPostgres) GetCategoriesList() ([]models.CarCategory, error) {
	var categoriesList []models.CarCategory

	categoriesListQuery := fmt.Sprintf(`SELECT id, name, photo FROM %s`, carCategoryTable)
	err := r.dash.Select(&categoriesList, categoriesListQuery)

	return categoriesList, err
}

func (r *RentCarsPostgres) GetCarsByCategoryId(categoryId int) ([]models.CarByCategoryId, error) {
	var cars []models.CarByCategoryId

	query := fmt.Sprintf(`SELECT car.id, car.price, car.photo, model.name model_name, company.name company_name FROM %[1]s car LEFT JOIN %[2]s model ON car.car_model_id = model.id LEFT JOIN %[3]s company ON car.rent_car_company_id = company.id`, carsTable, carsModelTable, carsCompanyTable)
	err := r.dash.Select(&cars, query)

	return cars, err
}
