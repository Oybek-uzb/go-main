package postgres

import (
	"abir/models"
	"abir/pkg/utils"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
)

type UtilsPostgres struct {
	db   *sqlx.DB
	dash *sqlx.DB
}

func NewUtilsPostgres(db *sqlx.DB, dash *sqlx.DB) *UtilsPostgres {
	return &UtilsPostgres{db: db, dash: dash}
}

func (r *UtilsPostgres) GetColors(langId int) ([]models.Color, error) {
	var lists []models.Color
	query := fmt.Sprintf("SELECT c.id, COALESCE(cl.name, c.name) as name, c.hex_code FROM %s c LEFT JOIN %s cl on c.id = cl.color_id WHERE cl.language_id = $1", colorsTable, colorsLangTable)
	err := r.dash.Select(&lists, query, langId)
	return lists, err
}
func (r *UtilsPostgres) GetCarMarkas() ([]models.CarMarka, error) {
	var lists []models.CarMarka
	query := fmt.Sprintf("SELECT id, name FROM %s", carMarkasTable)
	err := r.dash.Select(&lists, query)
	return lists, err
}
func (r *UtilsPostgres) GetCarModels(id int) ([]models.CarModel, error) {
	var lists []models.CarModel
	query := fmt.Sprintf("SELECT id, name FROM %s WHERE carmarka_id=$1", carModelsTable)
	err := r.dash.Select(&lists, query, id)
	return lists, err
}

func (r *UtilsPostgres) GetRegions(langId int) ([]models.Region, error) {
	var lists []models.Region
	query := fmt.Sprintf("SELECT r.id, COALESCE(rl.name, r.name) as name, is_city FROM %s r LEFT JOIN %s rl on r.id = rl.region_id WHERE rl.language_id = $1", regionsTable, regionsLangTable)
	err := r.dash.Select(&lists, query, langId)
	return lists, err
}
func (r *UtilsPostgres) GetDistricts(langId int, id int) ([]models.District, error) {
	var lists []models.District
	query := fmt.Sprintf("SELECT d.id, COALESCE(dl.name, d.name) as name FROM %s d LEFT JOIN %s dl on d.id = dl.district_id WHERE dl.language_id = $1 AND d.region_id = $2", districtsTable, districtsLangTable)
	err := r.dash.Select(&lists, query, langId, id)
	return lists, err
}
func (r *UtilsPostgres) GetDistrictsArr(langId int) (map[int]string, error) {
	var lists []models.District
	districts := make(map[int]string)
	query := fmt.Sprintf("SELECT d.id, COALESCE(dl.name, d.name) as name FROM %s d LEFT JOIN %s dl on d.id = dl.district_id WHERE dl.language_id = $1", districtsTable, districtsLangTable)
	err := r.dash.Select(&lists, query, langId)
	for _, list := range lists {
		districts[list.Id] = list.Name
	}
	districts[0] = utils.Translation["capital"][langId]
	return districts, err
}

func (r *UtilsPostgres) GetTariffs(langId int) (map[string]string, error) {
	var lists []models.Tariff
	tariffs := make(map[string]string)
	query := fmt.Sprintf("SELECT t.id, COALESCE(tl.name, t.name) as name FROM %s t LEFT JOIN %s tl on t.id = tl.tariff_id WHERE tl.language_id = $1", tariffsTable, tariffsLangTable)
	err := r.dash.Select(&lists, query, langId)
	for _, list := range lists {
		tariffs[strconv.Itoa(list.Id)] = list.Name
	}
	return tariffs, err
}

func (r *UtilsPostgres) DriverCancelOrderOptions(langId int) ([]models.DriverCancelOrderOptions, error) {
	var lists []models.DriverCancelOrderOptions
	query := fmt.Sprintf("SELECT o.id, COALESCE(ol.options, o.options) as options FROM %s o LEFT JOIN %s ol on o.id = ol.driver_cancel_order_id WHERE ol.language_id = $1", driverCancelOrderOptionsTable, driverCancelOrderOptionsLangTable)
	err := r.dash.Select(&lists, query, langId)
	return lists, err
}
func (r *UtilsPostgres) ClientCancelOrderOptions(langId int, optionType string) ([]models.ClientCancelOrderOptions, error) {
	var lists []models.ClientCancelOrderOptions
	query := fmt.Sprintf("SELECT o.id, COALESCE(ol.options, o.options) as options, type FROM %s o LEFT JOIN %s ol on o.id = ol.cancel_order_id WHERE ol.language_id = $1 AND o.type = $2", clientCancelOrderOptionsTable, clientCancelOrderOptionsLangTable)
	err := r.dash.Select(&lists, query, langId, optionType)
	return lists, err
}

func (r *UtilsPostgres) ClientRateOptions(langId, rate int, optionType string) ([]models.ClientRateOptions, error) {
	var lists []models.ClientRateOptions
	query := fmt.Sprintf("SELECT o.id, o.icon, COALESCE(ol.options, o.options) as options, type "+
		"FROM %s o LEFT JOIN %s ol on o.id = ol.driver_assessment_id "+
		"WHERE ol.language_id = $1 AND o.type = $2 AND o.rating = $3",
		clientRateOptionsTable,
		clientRateOptionsLangTable)
	err := r.dash.Select(&lists, query, langId, optionType, rate)

	for i, list := range lists {
		if list.Icon != nil {
			lists[i].Icon = utils.GetFileUrl(strings.Split(*list.Icon, "/"))
		}
	}

	return lists, err
}
