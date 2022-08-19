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

func (r *RentCarsPostgres) GetMyCarParkByCompanyId(userId, companyId int, inDiscount bool) ([]models.Car, error) {
	var cars []models.Car

	query := fmt.Sprintf(`SELECT car.id, car.price, car.photo, model.name model_name, car.in_discount, car.discount, company.name company_name FROM %[1]s car LEFT JOIN %[2]s model ON car.car_model_id = model.id LEFT JOIN %[3]s company ON car.rent_car_company_id = company.id WHERE car.rent_car_company_id=$1`, carsTable, carsModelTable, carsCompanyTable)

	if inDiscount {
		var addInDiscount = `AND car.in_discount=true`
		query = query + addInDiscount
	}

	err := r.dash.Select(&cars, query, companyId)

	return cars, err

}

func (r *RentCarsPostgres) GetMyCompanyById(userId, companyId int) (models.CarCompany, error) {
	var company models.CarCompany

	query := fmt.Sprintf(`SELECT id, photo, name, phone_number, description FROM %s WHERE owner_id=$1 AND id=$2`, carsCompanyTable)
	err := r.dash.Get(&company, query, userId, companyId)

	return company, err
}

func (r *RentCarsPostgres) GetMyCompaniesList(userId int) ([]models.CarCompany, error) {
	var companiesList []models.CarCompany

	companiesListQuery := fmt.Sprintf(`SELECT id, name, photo, phone_number FROM %s WHERE owner_id = $1`, carsCompanyTable)
	err := r.dash.Select(&companiesList, companiesListQuery, userId)

	return companiesList, err
}

func (r *RentCarsPostgres) GetCarByCompanyIdCarId(companyId, carId, langId int) (models.Car, error) {
	var car models.Car

	query := fmt.Sprintf(`SELECT car.id, car.price, car.photo, car.conditioner, car.description, car.phone_number, car.in_discount, car.discount, per_type.name per_type_name, fc_type_lang.name fc_type_name, model.name model_name, company.name company_name FROM %[1]s car LEFT JOIN %[2]s model ON car.car_model_id = model.id LEFT JOIN %[3]s company ON car.rent_car_company_id = company.id LEFT JOIN %[4]s fc_type ON car.fc_type_id = fc_type.id LEFT JOIN %[5]s fc_type_lang ON fc_type.id = fc_type_lang.fc_type_id LEFT JOIN %[6]s per_type ON car.per_car_id = per_type.id LEFT JOIN %[7]s per_type_lang ON per_type.id = per_type_lang.per_car_id WHERE per_type_lang.language_id = $3 AND fc_type_lang.language_id = $3 AND car.rent_car_company_id = $1 AND car.id = $2`, carsTable, carsModelTable, carsCompanyTable, fcTypeTable, fcTypeTableLang, perTypeTable, perTypeTableLang)
	err := r.dash.Get(&car, query, companyId, carId, langId)
	if err != nil {
		return models.Car{}, err
	}

	return car, err
}

func (r *RentCarsPostgres) GetCarsByCompanyId(companyId int) (models.CarCompanyDetails, error) {
	var carCompany models.CarCompanyDetails
	var cars []models.CarByCompanyId

	query := fmt.Sprintf(`SELECT id, name, photo, web_site, description, phone_number FROM %s WHERE id=$1`, carsCompanyTable)
	err := r.dash.Get(&carCompany, query, companyId)

	if err != nil {
		return models.CarCompanyDetails{}, err
	}

	query = fmt.Sprintf(`SELECT car.id, car.price, car.photo, car.in_discount, car.discount, model.name model_name, company.name company_name FROM %[1]s car LEFT JOIN %[2]s model ON car.car_model_id = model.id LEFT JOIN %[3]s company ON car.rent_car_company_id = company.id WHERE car.rent_car_company_id=$1`, carsTable, carsModelTable, carsCompanyTable)
	err = r.dash.Select(&cars, query, companyId)

	if err == nil {
		carCompany.Cars = cars
	}

	return carCompany, err
}

func (r *RentCarsPostgres) GetCompaniesList() ([]models.CarCompany, error) {
	var companiesList []models.CarCompany

	categoriesListQuery := fmt.Sprintf(`SELECT id, name, photo FROM %s`, carsCompanyTable)
	err := r.dash.Select(&companiesList, categoriesListQuery)

	return companiesList, err
}

func (r *RentCarsPostgres) GetCarByCategoryIdCarId(categoryId, carId, langId int) (models.Car, error) {
	var car models.Car

	query := fmt.Sprintf(`SELECT car.id, car.price, car.photo, car.conditioner, car.description, car.phone_number, car.discount, car.in_discount, per_type.name per_type_name, fc_type_lang.name fc_type_name, model.name model_name, company.name company_name FROM %[1]s car LEFT JOIN %[2]s model ON car.car_model_id = model.id LEFT JOIN %[3]s company ON car.rent_car_company_id = company.id LEFT JOIN %[4]s fc_type ON car.fc_type_id = fc_type.id LEFT JOIN %[5]s fc_type_lang ON fc_type.id = fc_type_lang.fc_type_id LEFT JOIN %[6]s per_type ON car.per_car_id = per_type.id LEFT JOIN %[7]s per_type_lang ON per_type.id = per_type_lang.per_car_id WHERE per_type_lang.language_id = $3 AND fc_type_lang.language_id = $3 AND car.rent_car_company_id = $1 AND car.id = $2`, carsTable, carsModelTable, carsCompanyTable, fcTypeTable, fcTypeTableLang, perTypeTable, perTypeTableLang)
	err := r.dash.Get(&car, query, categoryId, carId, langId)
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

	query := fmt.Sprintf(`SELECT car.id, car.price, car.photo, car.discount, car.in_discount, model.name model_name, company.name company_name FROM %[1]s car LEFT JOIN %[2]s model ON car.car_model_id = model.id LEFT JOIN %[3]s company ON car.rent_car_company_id = company.id WHERE car.category_car_id=$1`, carsTable, carsModelTable, carsCompanyTable)
	err := r.dash.Select(&cars, query, categoryId)

	return cars, err
}

func NewRentCarsPostgres(db *sqlx.DB, dash *sqlx.DB) *RentCarsPostgres {
	return &RentCarsPostgres{db: db, dash: dash}
}
