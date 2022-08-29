package postgres

import (
	"abir/models"
	"abir/pkg/utils"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type RentCarsPostgres struct {
	db   *sqlx.DB
	dash *sqlx.DB
}

func (r *RentCarsPostgres) DeleteMyCar(userId, carId, companyId int) (int, error) {
	var deletedCarId int

	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1 RETURNING id`, carsTable)
	row := r.dash.QueryRow(query, carId)
	err := row.Scan(&deletedCarId)
	if err != nil {
		return 0, err
	}

	return deletedCarId, nil
}

func (r *RentCarsPostgres) PostMyCompany(userId int, company models.RentMyCompanyCreate) (companyId int, err error) {
	details := utils.CheckForNil(company)
	query := fmt.Sprintf("INSERT INTO %s (name, photo, description, web_site, phone_number, status, owner_id) SELECT $1,$2,$3,$4,$5,$6,$7 RETURNING id", carsCompanyTable)

	row := r.dash.QueryRow(query, details["CompanyName"], details["Photo"], details["Description"], details["WebSite"], details["PhoneNumber"], "new", userId)
	err = row.Scan(&companyId)
	if err != nil {
		return 0, err
	}

	return
}

func (r *RentCarsPostgres) PostMyCar(userId, companyId int, car models.CarCreate) (int, error) {
	var regionId int
	var carId int

	details := utils.CheckForNil(car)

	q := fmt.Sprintf(`SELECT region_id FROM %s WHERE id=$1`, districtsTable)
	qRow := r.dash.QueryRow(q, details["DistrictId"])
	err := qRow.Scan(&regionId)
	if err != nil {
		return 0, err
	}

	query := fmt.Sprintf("INSERT INTO %s (conditioner, photo, description, price, phone_number, status, car_marka_id, car_model_id, category_car_id, color_id, district_id, fc_type_id, per_car_id, rent_car_company_id, discount, in_discount, transmission_id, consumption_fuel, region_id) SELECT $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19 RETURNING id", carsTable)

	row := r.dash.QueryRow(query, details["Conditioner"], details["Photo"], details["Description"], details["Price"], details["PhoneNumber"], true, details["MarkId"], details["ModelId"], details["CategoryId"], details["ColorId"], details["DistrictId"], details["FCTypeId"], details["PerCarId"], companyId, 0, false, details["TransmissionId"], details["ConsumptionFuel"], regionId)
	err = row.Scan(&carId)
	if err != nil {
		return 0, err
	}

	return carId, nil
}

func (r *RentCarsPostgres) PutMyCar(userId, carId, companyId int, car models.CarCreate) (int, error) {
	var regionId int
	var updatedCarId int

	details := utils.CheckForNil(car)

	q := fmt.Sprintf(`SELECT region_id FROM %s WHERE id=$1`, districtsTable)
	qRow := r.dash.QueryRow(q, details["DistrictId"])
	err := qRow.Scan(&regionId)
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf("UPDATE %s SET conditioner=$1, photo=$2, description=$3, price=$4, phone_number=$5, status=$6, car_marka_id=$7, car_model_id=$8, category_car_id=$9, color_id=$10, district_id=$11, fc_type_id=$12, per_car_id=$13, rent_car_company_id=$14, discount=$15, in_discount=$16, transmission_id=$17, consumption_fuel=$18, region_id=$19 WHERE id=$20 RETURNING id", carsTable)

	row := r.dash.QueryRow(query, details["Conditioner"], details["Photo"], details["Description"], details["Price"], details["PhoneNumber"], true, details["MarkId"], details["ModelId"], details["CategoryId"], details["ColorId"], details["DistrictId"], details["FCTypeId"], details["PerCarId"], companyId, 0, false, details["TransmissionId"], details["ConsumptionFuel"], regionId, carId)
	err = row.Scan(&updatedCarId)
	if err != nil {
		return 0, err
	}

	return updatedCarId, nil
}

func (r *RentCarsPostgres) PostRentCarByCarId(userId, carId int, rentCarDetails models.RentCarDetails) (rentId int, err error) {
	details := utils.CheckForNil(rentCarDetails)

	query := fmt.Sprintf("INSERT INTO %s (user_id, rent_cars_id, description, start_time, end_time) SELECT $1,$2,$3,$4,$5 RETURNING id", rentCarTable)

	row := r.dash.QueryRow(query, userId, carId, details["Description"], details["FromDate"], details["ToDate"])
	err = row.Scan(&rentId)
	if err != nil {
		return 0, err
	}

	return
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

	query := fmt.Sprintf(`SELECT car.id, car.price, car.photo, car.conditioner, car.description, car.phone_number, car.in_discount, car.discount, per_type_lang.name per_type_name, fc_type_lang.name fc_type_name, model.name model_name, company.name company_name, transmission_type_lang.name transmission_name, color_lang.name color_name, color.hex_code hex_code, car.consumption_fuel FROM %[1]s car LEFT JOIN %[2]s model ON car.car_model_id = model.id LEFT JOIN %[3]s company ON car.rent_car_company_id = company.id LEFT JOIN %[4]s fc_type ON car.fc_type_id = fc_type.id LEFT JOIN %[5]s fc_type_lang ON fc_type.id = fc_type_lang.fc_type_id LEFT JOIN %[6]s per_type ON car.per_car_id = per_type.id LEFT JOIN %[7]s per_type_lang ON per_type.id = per_type_lang.per_car_id LEFT JOIN %[8]s transmission_type ON car.transmission_id = transmission_type.id LEFT JOIN %[9]s transmission_type_lang ON transmission_type.id=transmission_type_lang.transmission_id LEFT JOIN %[10]s color ON car.color_id = color.id LEFT JOIN %[11]s color_lang ON color.id=color_lang.color_id WHERE per_type_lang.language_id = $3 AND fc_type_lang.language_id = $3 AND car.rent_car_company_id = $1 AND car.id = $2 AND color_lang.language_id = $3 AND transmission_type_lang.language_id=$3`, carsTable, carsModelTable, carsCompanyTable, fcTypeTable, fcTypeTableLang, perTypeTable, perTypeTableLang, transmissionTypeTable, transmissionTypeTableLang, colorsTable, colorsLangTable)
	err := r.dash.Get(&car, query, companyId, carId, langId)
	if err != nil {
		return models.Car{}, err
	}
	//query := fmt.Sprintf(`SELECT car.id, car.price, car.photo, car.conditioner, car.description, car.phone_number, car.discount, car.in_discount, per_type_lang.name per_type_name, fc_type_lang.name fc_type_name, model.name model_name, company.name company_name, transmission_type_lang.name transmission_name, color_lang.name color_name, color.hex_code hex_code FROM %[1]s car LEFT JOIN %[2]s model ON car.car_model_id = model.id LEFT JOIN %[3]s company ON car.rent_car_company_id = company.id LEFT JOIN %[4]s fc_type ON car.fc_type_id = fc_type.id LEFT JOIN %[5]s fc_type_lang ON fc_type.id = fc_type_lang.fc_type_id LEFT JOIN %[6]s per_type ON car.per_car_id = per_type.id LEFT JOIN %[7]s per_type_lang ON per_type.id = per_type_lang.per_car_id LEFT JOIN %[8]s transmission_type ON car.transmission_id = transmission_type.id LEFT JOIN %[9]s transmission_type_lang ON transmission_type.id=transmission_type_lang.transmission_id LEFT JOIN %[10]s color ON car.color_id = color.id LEFT JOIN %[11]s color_lang ON color.id=color_lang.color_id WHERE per_type_lang.language_id = $3 AND fc_type_lang.language_id = $3 AND car.category_car_id = $1 AND car.id = $2 AND color_lang.language_id = $3 AND transmission_type_lang.language_id=$3`, carsTable, carsModelTable, carsCompanyTable, fcTypeTable, fcTypeTableLang, perTypeTable, perTypeTableLang, transmissionTypeTable, transmissionTypeTableLang, colorsTable, colorsLangTable)

	if car.ConsumptionFuel != nil {
		var res string
		conFuel := strings.Split(*car.ConsumptionFuel, "/")

		if langId == 2 {
			res = conFuel[0] + "л/" + conFuel[1] + "км"
		} else {
			res = conFuel[0] + "l/" + conFuel[1] + "km"
		}

		car.ConsumptionFuel = &res
	}

	return car, err
}

func (r *RentCarsPostgres) GetCarsByCompanyId(companyId int) (models.CarCompanyDetails, error) {
	var carCompany models.CarCompanyDetails
	var cars []models.CarByCompanyId

	query := fmt.Sprintf(`SELECT id, name, photo, web_site, description, phone_number FROM %s WHERE id=$1 AND status='checked'`, carsCompanyTable)
	err := r.dash.Get(&carCompany, query, companyId)

	if err != nil {
		return models.CarCompanyDetails{}, err
	}

	query = fmt.Sprintf(`SELECT car.id, car.price, car.photo, car.in_discount, car.discount, model.name model_name, company.name company_name FROM %[1]s car LEFT JOIN %[2]s model ON car.car_model_id = model.id LEFT JOIN %[3]s company ON car.rent_car_company_id = company.id WHERE car.rent_car_company_id=$1 AND company.status='checked'`, carsTable, carsModelTable, carsCompanyTable)
	err = r.dash.Select(&cars, query, companyId)

	if err == nil {
		carCompany.Cars = cars
	}

	return carCompany, err
}

func (r *RentCarsPostgres) GetCompaniesList() ([]models.CarCompany, error) {
	var companiesList []models.CarCompany

	categoriesListQuery := fmt.Sprintf(`SELECT id, name, photo FROM %s WHERE status='checked'`, carsCompanyTable)
	err := r.dash.Select(&companiesList, categoriesListQuery)

	return companiesList, err
}

func (r *RentCarsPostgres) GetCarByCategoryIdCarId(categoryId, carId, langId int) (models.Car, error) {
	var car models.Car

	query := fmt.Sprintf(`SELECT car.id, car.price, car.photo, car.conditioner, car.description, car.phone_number, car.discount, car.in_discount, per_type_lang.name per_type_name, fc_type_lang.name fc_type_name, model.name model_name, company.name company_name, transmission_type_lang.name transmission_name, color_lang.name color_name, car.consumption_fuel, color.hex_code hex_code FROM %[1]s car LEFT JOIN %[2]s model ON car.car_model_id = model.id LEFT JOIN %[3]s company ON car.rent_car_company_id = company.id LEFT JOIN %[4]s fc_type ON car.fc_type_id = fc_type.id LEFT JOIN %[5]s fc_type_lang ON fc_type.id = fc_type_lang.fc_type_id LEFT JOIN %[6]s per_type ON car.per_car_id = per_type.id LEFT JOIN %[7]s per_type_lang ON per_type.id = per_type_lang.per_car_id LEFT JOIN %[8]s transmission_type ON car.transmission_id = transmission_type.id LEFT JOIN %[9]s transmission_type_lang ON transmission_type.id=transmission_type_lang.transmission_id LEFT JOIN %[10]s color ON car.color_id = color.id LEFT JOIN %[11]s color_lang ON color.id=color_lang.color_id WHERE per_type_lang.language_id = $3 AND fc_type_lang.language_id = $3 AND car.category_car_id = $1 AND car.id = $2 AND color_lang.language_id = $3 AND transmission_type_lang.language_id=$3`, carsTable, carsModelTable, carsCompanyTable, fcTypeTable, fcTypeTableLang, perTypeTable, perTypeTableLang, transmissionTypeTable, transmissionTypeTableLang, colorsTable, colorsLangTable)
	err := r.dash.Get(&car, query, categoryId, carId, langId)
	if err != nil {
		return models.Car{}, err
	}

	if car.ConsumptionFuel != nil {
		var res string
		conFuel := strings.Split(*car.ConsumptionFuel, "/")

		if langId == 2 {
			res = conFuel[0] + "л/" + conFuel[1] + "км"
		} else {
			res = conFuel[0] + "l/" + conFuel[1] + "km"
		}

		car.ConsumptionFuel = &res
	}

	return car, err
}

func (r *RentCarsPostgres) GetCategoriesList(langId int) ([]models.CarCategory, error) {
	var categoriesList []models.CarCategory

	categoriesListQuery := fmt.Sprintf(`SELECT car_category.id, car_category.photo, car_category_lang.name FROM %[1]s car_category LEFT JOIN %[2]s car_category_lang ON car_category.id = car_category_lang.category_car_id WHERE car_category_lang.language_id=$1 AND car_category.car_type='for_events'`, carCategoryTable, carCategoryTableLang)
	err := r.dash.Select(&categoriesList, categoriesListQuery, langId)

	return categoriesList, err
}

func (r *RentCarsPostgres) GetCategoriesForEvents(langId int) ([]models.CarCategory, error) {
	var categories []models.CarCategory

	categoriesListQuery := fmt.Sprintf(`SELECT car_category.id, car_category.photo, car_category_lang.name FROM %[1]s car_category LEFT JOIN %[2]s car_category_lang ON car_category.id = car_category_lang.category_car_id WHERE car_category_lang.language_id=$1`, carCategoryTable, carCategoryTableLang)
	err := r.dash.Select(&categories, categoriesListQuery, langId)

	return categories, err
}

func (r *RentCarsPostgres) GetCarsByCategoryId(categoryId int) ([]models.CarByCategoryId, error) {
	var cars []models.CarByCategoryId

	query := fmt.Sprintf(`SELECT car.id, car.price, car.photo, car.discount, car.in_discount, model.name model_name, company.name company_name FROM %[1]s car LEFT JOIN %[2]s model ON car.car_model_id = model.id LEFT JOIN %[3]s company ON car.rent_car_company_id = company.id WHERE car.category_car_id=$1 AND company.status='checked'`, carsTable, carsModelTable, carsCompanyTable)
	err := r.dash.Select(&cars, query, categoryId)

	return cars, err
}

func NewRentCarsPostgres(db *sqlx.DB, dash *sqlx.DB) *RentCarsPostgres {
	return &RentCarsPostgres{db: db, dash: dash}
}
