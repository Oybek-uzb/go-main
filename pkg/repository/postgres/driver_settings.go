package postgres

import (
	"abir/models"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type DriverSettingsPostgres struct {
	db   *sqlx.DB
	dash *sqlx.DB
}

func NewDriverSettingsPostgres(db *sqlx.DB, dash *sqlx.DB) *DriverSettingsPostgres {
	return &DriverSettingsPostgres{db: db, dash: dash}
}

func (r *DriverSettingsPostgres) GetTariffs(userId, langId int) ([]models.DriverTariffs, error) {
	var usr models.User
	usrQuery := fmt.Sprintf("SELECT driver_id FROM %s WHERE id=$1", usersTable)
	err := r.db.Get(&usr, usrQuery, userId)
	if err != nil {
		return []models.DriverTariffs{}, err
	}
	var car models.DriverCar
	carQuery := fmt.Sprintf("SELECT car_model_id FROM %s WHERE driver_id=$1", driverCarTable)
	err = r.dash.Get(&car, carQuery, *usr.DriverId)
	if err != nil {
		return []models.DriverTariffs{}, err
	}
	var tariffs []models.DriverTariffs
	subQuery := fmt.Sprintf("(SELECT is_active FROM public.%s WHERE tariff_id = t.tariff_id AND user_id = $3 AND is_active=true)", driverEnabledTariffsTable)
	tariffQuery := fmt.Sprintf("SELECT t.tariff_id as id, tl.name as name, CASE WHEN EXISTS %[1]s THEN %[1]s ELSE false END as is_active FROM %[2]s t LEFT JOIN %[3]s tl ON t.tariff_id = tl.tariff_id WHERE t.car_model_id=$1 AND tl.language_id=$2", subQuery, tariffCarModelTable, tariffsLangTable)
	err = r.dash.Select(&tariffs, tariffQuery, *car.CarModelId, langId, userId)
	return tariffs, err
}

func (r *DriverSettingsPostgres) TariffsEnable(userId, tariffId int, isActive bool) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	updateQuery := fmt.Sprintf("UPDATE %s SET user_id=$1,tariff_id=$2,is_active=$3 WHERE user_id=$1 AND tariff_id=$2", driverEnabledTariffsTable)
	res, err := tx.Exec(updateQuery, userId, tariffId, isActive)
	updated, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if updated == 0 {
		insertQuery := fmt.Sprintf("INSERT INTO %s (user_id, tariff_id, is_active) VALUES ($1,$2,$3)", driverEnabledTariffsTable)
		_, err = tx.Exec(insertQuery, userId, tariffId, isActive)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (r *DriverSettingsPostgres) SetOnline(userId int, isActive int) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	sts := "offline"
	if isActive == 1 {
		sts = "online"
	}
	updateQuery := fmt.Sprintf("UPDATE %s SET user_id=$1,driver_status=$2 WHERE user_id=$1", driverStatusesTable)
	res, err := tx.Exec(updateQuery, userId, sts)
	updated, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if updated == 0 {
		insertQuery := fmt.Sprintf("INSERT INTO %s (user_id, driver_status) VALUES ($1,$2)", driverStatusesTable)
		_, err = tx.Exec(insertQuery, userId, sts)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
