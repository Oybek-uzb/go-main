package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

const (
	clientsTable   = "clients"
	usersTable     = "users"
	driverTable    = "driver_driver"
	driverCarTable = "driver_car"

	colorsTable                       = "dictionary_color"
	colorsLangTable                   = "dictionary_color_i18n"
	carMarkasTable                    = "dictionary_car_marka"
	carModelsTable                    = "dictionary_car_model"
	regionsTable                      = "dictionary_region"
	regionsLangTable                  = "dictionary_region_i18n"
	districtsTable                    = "dictionary_district"
	districtsLangTable                = "dictionary_district_i18n"
	tariffsTable                      = "dictionary_tariff"
	tariffsLangTable                  = "dictionary_tariff_i18n"
	driverCancelOrderOptionsTable     = "dictionary_driver_cancel_order"
	driverCancelOrderOptionsLangTable = "dictionary_driver_cancel_order_i18n"
	clientCancelOrderOptionsTable     = "dictionary_cancel_order"
	clientCancelOrderOptionsLangTable = "dictionary_cancel_order_i18n"
	driverVerificationsTable          = "driver_deficiencies"
	tariffCarModelTable               = "dictionary_tariff_car_model"
	driverEnabledTariffsTable         = "driver_enabled_tariffs"
	routeCityTaxiTable                = "route_city_taxi"
	clientRateOptionsTable            = "dictionary_driver_assessment"
	clientRateOptionsLangTable        = "dictionary_driver_assessment_i18n"

	savedAddressesTable       = "saved_addresses"
	creditCardsTable          = "credit_cards"
	ridesTable                = "rides"
	ordersTable               = "orders"
	interregionalOrdersTable  = "interregional_orders"
	cityOrdersTable           = "city_orders"
	canceledOrdersTable       = "canceled_orders"
	canceledOrderReasonsTable = "canceled_order_reasons"
	ratedOrdersTable          = "rated_orders"
	ratedOrderReasonsTable    = "rated_order_reasons"
	chatMessagesTable         = "chat_messages"
	rideViewCountsTable       = "ride_view_counts"
	driverStatusesTable       = "driver_statuses"

	clientType             = "client"
	driverType             = "driver"
	orderCityType          = "city"
	orderInterregionalType = "interregional"

	carCategoryTableLang = "dictionary_category_car_i18n"
	carCategoryTable     = "dictionary_category_car"
	carsTable            = "rent_car_rent_car_cars"
	carsModelTable       = "dictionary_car_model"
	carsCompanyTable     = "rent_car_rent_car_company"

	rentCarTable = "rent_car_book_car"

	fcTypeTable     = "dictionary_fc_type"
	fcTypeTableLang = "dictionary_fc_type_i18n"

	perTypeTable     = "dictionary_per_car"
	perTypeTableLang = "dictionary_per_car_i18n"

	transmissionTypeTable     = "dictionary_transmission"
	transmissionTypeTableLang = "dictionary_transmission_i18n"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
	Schema   string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s search_path=%s", cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode, cfg.Schema))
	if err != nil {
		return nil, err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return nil, pingErr
	}

	return db, nil
}
