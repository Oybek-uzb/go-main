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

func (r *RentCarsPostgres) GetCarByCategoryIdCarId(categoryId, carId, langId int) (models.Car, error) {
	var car models.Car

	// per_type.name per_type_name, LEFT JOIN %[6]s per_type ON car.per_car_id = per_type.id LEFT JOIN %[7]s per_type_lang ON per_type.id = per_type_lang.per_car_id WHERE per_type_lang.language_id = $1
	query := fmt.Sprintf(`SELECT car.id, car.price, car.photo, car.conditioner, car.phone_number, per_type.name per_type_name, fc_type_lang.name fc_type_name, model.name model_name, company.name company_name FROM %[1]s car LEFT JOIN %[2]s model ON car.car_model_id = model.id LEFT JOIN %[3]s company ON car.rent_car_company_id = company.id LEFT JOIN %[4]s fc_type ON car.fc_type_id = fc_type.id LEFT JOIN %[5]s fc_type_lang ON fc_type.id = fc_type_lang.fc_type_id LEFT JOIN %[6]s per_type ON car.per_car_id = per_type.id LEFT JOIN %[7]s per_type_lang ON per_type.id = per_type_lang.per_car_id WHERE per_type_lang.language_id = $1 AND fc_type_lang.language_id = $1`, carsTable, carsModelTable, carsCompanyTable, fcTypeTable, fcTypeTableLang, perTypeTable, perTypeTableLang)
	err := r.dash.Get(&car, query, langId)
	if err != nil {
		return models.Car{}, err
	}

	return car, err
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

func NewRentCarsPostgres(db *sqlx.DB, dash *sqlx.DB) *RentCarsPostgres {
	return &RentCarsPostgres{db: db, dash: dash}
}