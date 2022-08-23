package postgres

import (
	"abir/models"
	"abir/pkg/utils"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type AuthPostgres struct {
	db   *sqlx.DB
	dash *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB, dash *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db, dash: dash}
}
func (r *AuthPostgres) CreateClient(user models.Client, userId int) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	var usr models.User
	usrQuery := fmt.Sprintf("SELECT u.client_id FROM %s u INNER JOIN %s c ON c.id = u.client_id WHERE u.id=$1", usersTable, clientsTable)
	err = tx.Get(&usr, usrQuery, userId)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return err
	}

	var surname, birthdate string
	if user.Surname == nil {
		surname = ""
	}
	if user.Birthdate == nil {
		birthdate = "2006-01-01"
	}

	if usr.ClientId == nil || err == sql.ErrNoRows {
		var id int
		query := fmt.Sprintf("INSERT INTO %s (name, surname, birthdate, gender, avatar) values ($1,$2,$3,$4,$5) RETURNING id", clientsTable)
		row := tx.QueryRow(query, user.Name, surname, birthdate, user.Gender, user.Avatar)
		if err = row.Scan(&id); err != nil {
			tx.Rollback()
			return err
		}
		userQuery := fmt.Sprintf("UPDATE %s SET client_id=$1 WHERE id=$2 RETURNING id", usersTable)
		_, err = tx.Exec(userQuery, id, userId)
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		setValues := make([]string, 0)
		args := make([]interface{}, 0)
		argId := 1
		if user.Name != nil {
			setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
			args = append(args, *user.Name)
			argId++
		}
		if user.Surname != nil {
			setValues = append(setValues, fmt.Sprintf("surname=$%d", argId))
			args = append(args, *user.Surname)
			argId++
		}
		if user.Birthdate != nil {
			setValues = append(setValues, fmt.Sprintf("birthdate=$%d", argId))
			args = append(args, *user.Birthdate)
			argId++
		}
		if user.Gender != nil {
			setValues = append(setValues, fmt.Sprintf("gender=$%d", argId))
			args = append(args, *user.Gender)
			argId++
		}
		if user.Avatar != nil {
			setValues = append(setValues, fmt.Sprintf("avatar=$%d", argId))
			if *user.Avatar != "" {
				args = append(args, user.Avatar)
			} else {
				args = append(args, nil)
			}
			argId++
		}
		setQuery := strings.Join(setValues, ", ")
		updateQuery := fmt.Sprintf(`UPDATE %s SET %s WHERE id = $%d`,
			clientsTable, setQuery, argId)
		args = append(args, *usr.ClientId)
		_, err = tx.Exec(updateQuery, args...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
func (r *AuthPostgres) ClientCheckPhone(phone string) error {
	query := fmt.Sprintf("SELECT id FROM %s WHERE login=$1 AND user_type=$2", usersTable)
	var usr models.User
	err := r.db.Get(&usr, query, phone, clientType)
	if usr.Id == 0 {
		return nil
	} else {
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		return errors.New("client with this number already exists")
	}
}
func (r *AuthPostgres) ClientUpdatePhone(userId int, phone string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	var usr models.User
	usrQuery := fmt.Sprintf("SELECT u.client_id FROM %s u INNER JOIN %s c ON c.id = u.client_id WHERE u.id=$1", usersTable, clientsTable)
	err = tx.Get(&usr, usrQuery, userId)
	if err != nil {
		tx.Rollback()
		return err
	}

	updateUserQuery := fmt.Sprintf(`UPDATE %s SET login = $1 WHERE id = $2`, usersTable)
	_, err = tx.Exec(updateUserQuery, phone, userId)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
func (r *AuthPostgres) DriverCheckPhone(phone string) error {
	query := fmt.Sprintf("SELECT id FROM %s WHERE login=$1 AND user_type=$2", usersTable)
	var usr models.User
	err := r.db.Get(&usr, query, phone, driverType)
	if usr.Id == 0 {
		return nil
	} else {
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		return errors.New("driver with this number already exists")
	}
}
func (r *AuthPostgres) DriverUpdatePhone(userId int, phone string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	var usr models.User
	usrQuery := fmt.Sprintf("SELECT u.driver_id FROM %s u INNER JOIN dashboard.%s c ON c.id = u.driver_id WHERE u.id=$1", usersTable, driverTable)
	err = tx.Get(&usr, usrQuery, userId)
	if err != nil {
		tx.Rollback()
		return err
	}
	updateUserQuery := fmt.Sprintf(`UPDATE %s SET login = $1 WHERE id = $2`, usersTable)
	_, err = tx.Exec(updateUserQuery, phone, userId)
	if err != nil {
		tx.Rollback()
		return err
	}

	updateDriverQuery := fmt.Sprintf(`UPDATE dashboard.%s SET phone = $1 WHERE id = $2`, driverTable)
	_, err = tx.Exec(updateDriverQuery, phone, usr.DriverId)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *AuthPostgres) CreateOrUpdateClient(user models.User) (int, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf("UPDATE %s SET login=$1, password_hash=$2 WHERE login=$3 AND user_type=$4 RETURNING id", usersTable)
	res, err := tx.Exec(query, user.Login, '0', user.Login, clientType)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	var id int
	if cnt == 0 {
		createQuery := fmt.Sprintf("INSERT INTO %s (login, password_hash, user_type) values ($1,$2,$3) RETURNING id", usersTable)
		row := tx.QueryRow(createQuery, user.Login, user.Password, clientType)
		if err := row.Scan(&id); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	return id, tx.Commit()
}

func (r *AuthPostgres) GetUser(login, password, userType string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE login=$1 AND password_hash=$2 AND user_type=$3", usersTable)
	err := r.db.Get(&user, query, login, password, userType)
	return user, err
}

func (r *AuthPostgres) GetClient(userId int) (models.Client, error) {
	var usr models.User
	usrQuery := fmt.Sprintf("SELECT client_id, login FROM %s WHERE id=$1", usersTable)
	err := r.db.Get(&usr, usrQuery, userId)
	if err != nil {
		return models.Client{}, err
	}
	if usr.ClientId == nil {
		return models.Client{}, err
	}
	var client models.Client
	clientQuery := fmt.Sprintf("SELECT id,name,surname,birthdate,gender,avatar FROM %s WHERE id=$1", clientsTable)
	err = r.db.Get(&client, clientQuery, *usr.ClientId)
	if err != nil {
		return models.Client{}, err
	}
	client.Id = userId
	client.Phone = &usr.Login
	if client.Avatar != nil {
		client.Avatar = utils.GetFileUrl(strings.Split(*client.Avatar, "/"))
	}
	return client, nil
}
func (r *AuthPostgres) GetDriverVerification(userId int) ([]models.DriverVerification, error) {
	var usr models.User
	usrQuery := fmt.Sprintf("SELECT driver_id FROM %s WHERE id=$1", usersTable)
	err := r.db.Get(&usr, usrQuery, userId)
	if err != nil {
		return []models.DriverVerification{}, err
	}
	if usr.DriverId == nil {
		return []models.DriverVerification{}, err
	}
	var verification []models.DriverVerification
	verificationQuery := fmt.Sprintf("SELECT DISTINCT ON (status) status, id, description FROM %s WHERE driver_id=$1 ORDER BY status, id desc", driverVerificationsTable)
	err = r.dash.Select(&verification, verificationQuery, *usr.DriverId)
	return verification, err
}
func (r *AuthPostgres) GetDriverId(userId int) (int, error) {
	var usr models.User
	usrQuery := fmt.Sprintf("SELECT driver_id FROM %s WHERE id=$1", usersTable)
	err := r.db.Get(&usr, usrQuery, userId)
	if err != nil {
		return 0, err
	}
	if usr.DriverId == nil {
		return 0, errors.New("driver not found")
	}
	return *usr.DriverId, nil
}
func (r *AuthPostgres) GetDriver(userId int) (models.Driver, error) {
	var usr models.User
	usrQuery := fmt.Sprintf("SELECT driver_id FROM %s WHERE id=$1", usersTable)
	err := r.db.Get(&usr, usrQuery, userId)
	if err != nil {
		return models.Driver{}, err
	}
	if usr.DriverId == nil {
		return models.Driver{}, err
	}
	var driver models.Driver
	driverQuery := fmt.Sprintf("SELECT id,name,surname,date_of_birth,gender,photo,phone,activity,status,passport_copy1,passport_copy2,passport_copy3,driver_license_photo1,driver_license_photo2,driver_license_photo3,driver_license,driver_license_expiration,document_type,passport_serial FROM %s WHERE id=$1", driverTable)
	err = r.dash.Get(&driver, driverQuery, *usr.DriverId)
	if err != nil {
		return models.Driver{}, err
	}
	driver.Id = userId
	var newDriver models.Driver
	newDriverQuery := fmt.Sprintf("SELECT driver_status as status FROM %s WHERE user_id=$1", driverStatusesTable)
	err = r.db.Get(&newDriver, newDriverQuery, userId)
	if err != nil && err != sql.ErrNoRows {
		return models.Driver{}, err
	}
	if newDriver.Status != nil && err != sql.ErrNoRows {
		driver.Status = newDriver.Status
	}
	if driver.Photo != nil {
		driver.Photo = utils.GetFileUrl(strings.Split(*driver.Photo, "/"))
	}
	if driver.PassportCopy1 != nil {
		driver.PassportCopy1 = utils.GetFileUrl(strings.Split(*driver.PassportCopy1, "/"))
	}
	if driver.PassportCopy2 != nil {
		driver.PassportCopy2 = utils.GetFileUrl(strings.Split(*driver.PassportCopy2, "/"))
	}
	if driver.PassportCopy3 != nil {
		driver.PassportCopy3 = utils.GetFileUrl(strings.Split(*driver.PassportCopy3, "/"))
	}
	if driver.DriverLicensePhoto1 != nil {
		driver.DriverLicensePhoto1 = utils.GetFileUrl(strings.Split(*driver.DriverLicensePhoto1, "/"))
	}
	if driver.DriverLicensePhoto2 != nil {
		driver.DriverLicensePhoto2 = utils.GetFileUrl(strings.Split(*driver.DriverLicensePhoto2, "/"))
	}
	if driver.DriverLicensePhoto3 != nil {
		driver.DriverLicensePhoto3 = utils.GetFileUrl(strings.Split(*driver.DriverLicensePhoto3, "/"))
	}
	return driver, nil
}

func (r *AuthPostgres) GetDriverCar(userId int) (models.DriverCar, error) {
	var usr models.User
	usrQuery := fmt.Sprintf("SELECT driver_id FROM %s WHERE id=$1", usersTable)
	err := r.db.Get(&usr, usrQuery, userId)
	if err != nil {
		return models.DriverCar{}, err
	}
	var car models.DriverCar
	carQuery := fmt.Sprintf("SELECT photo_texpasport1,photo_texpasport2,car_number,car_year,car_front,car_left,car_back,car_right,car_front_row,car_front_back,car_baggage,car_color_id,car_marka_id,car_model_id FROM %s WHERE driver_id=$1", driverCarTable)
	err = r.dash.Get(&car, carQuery, *usr.DriverId)
	if err != nil {
		return models.DriverCar{}, err
	}
	if car.PhotoTexpasport1 != nil {
		car.PhotoTexpasport1 = utils.GetFileUrl(strings.Split(*car.PhotoTexpasport1, "/"))
	}
	if car.PhotoTexpasport2 != nil {
		car.PhotoTexpasport2 = utils.GetFileUrl(strings.Split(*car.PhotoTexpasport2, "/"))
	}
	if car.CarFront != nil {
		car.CarFront = utils.GetFileUrl(strings.Split(*car.CarFront, "/"))
	}
	if car.CarLeft != nil {
		car.CarLeft = utils.GetFileUrl(strings.Split(*car.CarLeft, "/"))
	}
	if car.CarRight != nil {
		car.CarRight = utils.GetFileUrl(strings.Split(*car.CarRight, "/"))
	}
	if car.CarBack != nil {
		car.CarBack = utils.GetFileUrl(strings.Split(*car.CarBack, "/"))
	}
	if car.CarFrontRow != nil {
		car.CarFrontRow = utils.GetFileUrl(strings.Split(*car.CarFrontRow, "/"))
	}
	if car.CarFrontBack != nil {
		car.CarFrontBack = utils.GetFileUrl(strings.Split(*car.CarFrontBack, "/"))
	}
	if car.CarBaggage != nil {
		car.CarBaggage = utils.GetFileUrl(strings.Split(*car.CarBaggage, "/"))
	}
	return car, nil
}

func (r *AuthPostgres) ClientSendCode(login, password string) error {
	query := fmt.Sprintf("UPDATE %s SET password_hash=$1 WHERE login=$2 AND user_type=$3", usersTable)
	_, err := r.db.Exec(query, password, login, clientType)
	return err
}

func (r *AuthPostgres) DriverSendCode(login, password string) error {
	query := fmt.Sprintf("UPDATE %s SET password_hash=$1 WHERE login=$2 AND user_type=$3", usersTable)
	_, err := r.db.Exec(query, password, login, driverType)
	return err
}
func (r *AuthPostgres) CreateDriver(user models.Driver, userId int) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	tx1, err := r.dash.Beginx()
	if err != nil {
		return err
	}
	var usr models.User
	var drv models.Driver
	usrQuery := fmt.Sprintf("SELECT driver_id,login FROM %s WHERE id=$1", usersTable)
	err = tx.Get(&usr, usrQuery, userId)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		tx1.Rollback()
		return err
	}
	drvQuery := fmt.Sprintf("SELECT id FROM %s WHERE id=$1", driverTable)
	var drvErr error
	if usr.DriverId == nil {
		drvErr = tx1.Get(&drv, drvQuery, usr.DriverId)
	} else {
		drvErr = tx1.Get(&drv, drvQuery, *usr.DriverId)
	}

	if drvErr != nil && drvErr != sql.ErrNoRows {
		tx.Rollback()
		tx1.Rollback()
		return drvErr
	}
	if usr.DriverId == nil || err == sql.ErrNoRows || drvErr == sql.ErrNoRows {
		var id int
		query := fmt.Sprintf("INSERT INTO %s (name, surname, date_of_birth, gender, photo, phone, activity, status, created_at, updated_at, document_type) values ($1,$2,$3,$4,$5,$6,$7,$8, NOW(), NOW(), $9) RETURNING id", driverTable)
		row := tx1.QueryRow(query, user.Name, user.Surname, user.DateOfBirth, user.Gender, user.Photo, usr.Login, 100, "new", "passport")
		if err = row.Scan(&id); err != nil {
			tx.Rollback()
			tx1.Rollback()
			return err
		}
		carQuery := fmt.Sprintf("INSERT INTO %s (driver_id, created_at, updated_at) values ($1,NOW(),NOW()) RETURNING id", driverCarTable)
		_, err = tx1.Exec(carQuery, id)
		if err != nil {
			tx.Rollback()
			tx1.Rollback()
			return err
		}
		userQuery := fmt.Sprintf("UPDATE %s SET driver_id=$1 WHERE id=$2 RETURNING id", usersTable)
		_, err = tx.Exec(userQuery, id, userId)
		if err != nil {
			tx.Rollback()
			tx1.Rollback()
			return err
		}
	} else {
		setValues := make([]string, 0)
		args := make([]interface{}, 0)
		argId := 1
		if user.Name != nil {
			setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
			args = append(args, *user.Name)
			argId++
		}
		if user.Surname != nil {
			setValues = append(setValues, fmt.Sprintf("surname=$%d", argId))
			args = append(args, *user.Surname)
			argId++
		}
		if user.DateOfBirth != nil {
			setValues = append(setValues, fmt.Sprintf("date_of_birth=$%d", argId))
			args = append(args, *user.DateOfBirth)
			argId++
		}
		if user.Gender != nil {
			setValues = append(setValues, fmt.Sprintf("gender=$%d", argId))
			args = append(args, *user.Gender)
			argId++
		}
		if user.DocumentType != nil {
			setValues = append(setValues, fmt.Sprintf("document_type=$%d", argId))
			args = append(args, *user.DocumentType)
			argId++
		}
		if user.Photo != nil {
			setValues = append(setValues, fmt.Sprintf("photo=$%d", argId))
			if *user.Photo != "" {
				args = append(args, user.Photo)
			} else {
				args = append(args, nil)
			}
			argId++
		}
		setQuery := strings.Join(setValues, ", ")
		updateQuery := fmt.Sprintf(`UPDATE %s SET %s WHERE id = $%d`,
			driverTable, setQuery, argId)
		args = append(args, *usr.DriverId)
		_, err = tx1.Exec(updateQuery, args...)
		if err != nil {
			tx.Rollback()
			tx1.Rollback()
			return err
		}
	}

	if err := tx1.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (r *AuthPostgres) CreateOrUpdateDriver(user models.User) (int, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf("UPDATE %s SET login=$1, password_hash=$2 WHERE login=$1 AND user_type=$3 RETURNING id", usersTable)
	res, err := tx.Exec(query, user.Login, '0', driverType)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	var id int
	if cnt == 0 {
		createQuery := fmt.Sprintf("INSERT INTO %s (login, password_hash, user_type) values ($1,$2,$3) RETURNING id", usersTable)
		row := tx.QueryRow(createQuery, user.Login, user.Password, driverType)
		if err := row.Scan(&id); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	return id, tx.Commit()
}
func (r *AuthPostgres) SendForModerating(userId int) error {
	var usr models.User
	usrQuery := fmt.Sprintf("SELECT driver_id,login FROM %s WHERE id=$1", usersTable)
	err := r.db.Get(&usr, usrQuery, userId)
	if err != nil {
		return err
	}
	userQuery := fmt.Sprintf("UPDATE %s SET status=$1 WHERE id=$2", driverTable)
	_, err = r.dash.Exec(userQuery, "send_for_moderating", usr.DriverId)
	return err
}
func (r *AuthPostgres) UpdateDriver(user models.Driver, userId int) error {
	tx, err := r.dash.Beginx()
	if err != nil {
		return err
	}
	var usr models.User
	usrQuery := fmt.Sprintf("SELECT driver_id,login FROM %s WHERE id=$1", usersTable)
	err = r.db.Get(&usr, usrQuery, userId)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return err
	}
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1
	if user.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *user.Name)
		argId++
	}
	if user.DriverLicenseExpiration != nil {
		setValues = append(setValues, fmt.Sprintf("driver_license_expiration=$%d", argId))
		args = append(args, user.DriverLicenseExpiration)
		argId++
	}
	if user.DriverLicense != nil {
		setValues = append(setValues, fmt.Sprintf("driver_license=$%d", argId))
		args = append(args, user.DriverLicense)
		argId++
	}
	if user.PassportSerial != nil {
		setValues = append(setValues, fmt.Sprintf("passport_serial=$%d", argId))
		args = append(args, user.PassportSerial)
		argId++
	}
	if user.Surname != nil {
		setValues = append(setValues, fmt.Sprintf("surname=$%d", argId))
		args = append(args, *user.Surname)
		argId++
	}
	if user.DateOfBirth != nil {
		setValues = append(setValues, fmt.Sprintf("date_of_birth=$%d", argId))
		args = append(args, *user.DateOfBirth)
		argId++
	}
	if user.Gender != nil {
		setValues = append(setValues, fmt.Sprintf("gender=$%d", argId))
		args = append(args, *user.Gender)
		argId++
	}
	if user.DocumentType != nil {
		setValues = append(setValues, fmt.Sprintf("document_type=$%d", argId))
		args = append(args, *user.DocumentType)
		argId++
	}
	if user.Photo != nil {
		setValues = append(setValues, fmt.Sprintf("photo=$%d", argId))
		if *user.Photo != "" {
			args = append(args, user.Photo)
		} else {
			args = append(args, nil)
		}
		argId++
	}
	if user.PassportCopy1 != nil {
		setValues = append(setValues, fmt.Sprintf("passport_copy1=$%d", argId))
		if *user.PassportCopy1 != "" {
			args = append(args, user.PassportCopy1)
		} else {
			args = append(args, nil)
		}
		argId++
	}
	if user.PassportCopy2 != nil {
		setValues = append(setValues, fmt.Sprintf("passport_copy2=$%d", argId))
		if *user.PassportCopy2 != "" {
			args = append(args, user.PassportCopy2)
		} else {
			args = append(args, nil)
		}
		argId++
	}
	if user.PassportCopy3 != nil {
		setValues = append(setValues, fmt.Sprintf("passport_copy3=$%d", argId))
		if *user.PassportCopy3 != "" {
			args = append(args, user.PassportCopy3)
		} else {
			args = append(args, nil)
		}
		argId++
	}
	if user.DriverLicensePhoto1 != nil {
		setValues = append(setValues, fmt.Sprintf("driver_license_photo1=$%d", argId))
		if *user.DriverLicensePhoto1 != "" {
			args = append(args, user.DriverLicensePhoto1)
		} else {
			args = append(args, nil)
		}
		argId++
	}
	if user.DriverLicensePhoto2 != nil {
		setValues = append(setValues, fmt.Sprintf("driver_license_photo2=$%d", argId))
		if *user.DriverLicensePhoto2 != "" {
			args = append(args, user.DriverLicensePhoto2)
		} else {
			args = append(args, nil)
		}
		argId++
	}
	if user.DriverLicensePhoto3 != nil {
		setValues = append(setValues, fmt.Sprintf("driver_license_photo3=$%d", argId))
		if *user.DriverLicensePhoto3 != "" {
			args = append(args, user.DriverLicensePhoto3)
		} else {
			args = append(args, nil)
		}
		argId++
	}
	setQuery := strings.Join(setValues, ", ")
	updateQuery := fmt.Sprintf(`UPDATE %s SET %s WHERE id = $%d`,
		driverTable, setQuery, argId)
	args = append(args, *usr.DriverId)
	_, err = tx.Exec(updateQuery, args...)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (r *AuthPostgres) UpdateDriverCar(car models.DriverCar, userId int) error {
	tx, err := r.dash.Beginx()
	if err != nil {
		return err
	}
	var usr models.User
	var drvCar models.DriverCar
	usrQuery := fmt.Sprintf("SELECT driver_id FROM %s WHERE id=$1", usersTable)
	err = r.db.Get(&usr, usrQuery, userId)
	if err != nil {
		tx.Rollback()
		return errors.New("user not found")
	}
	if usr.DriverId == nil {
		return errors.New("driver not found")
	}
	carQuery := fmt.Sprintf("SELECT id FROM %s WHERE driver_id=$1", driverCarTable)
	carErr := tx.Get(&drvCar, carQuery, *usr.DriverId)
	if carErr != nil {
		tx.Rollback()
		return carErr
	}
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1
	if car.CarMarkaId != nil {
		setValues = append(setValues, fmt.Sprintf("car_marka_id=$%d", argId))
		args = append(args, car.CarMarkaId)
		argId++
	}
	if car.CarModelId != nil {
		setValues = append(setValues, fmt.Sprintf("car_model_id=$%d", argId))
		args = append(args, car.CarModelId)
		argId++
	}
	if car.CarColorId != nil {
		setValues = append(setValues, fmt.Sprintf("car_color_id=$%d", argId))
		args = append(args, car.CarColorId)
		argId++
	}
	if car.CarYear != nil {
		setValues = append(setValues, fmt.Sprintf("car_year=$%d", argId))
		args = append(args, car.CarYear)
		argId++
	}
	if car.CarNumber != nil {
		setValues = append(setValues, fmt.Sprintf("car_number=$%d", argId))
		args = append(args, car.CarNumber)
		argId++
	}
	if car.PhotoTexpasport1 != nil {
		setValues = append(setValues, fmt.Sprintf("photo_texpasport1=$%d", argId))
		if *car.PhotoTexpasport1 != "" {
			args = append(args, car.PhotoTexpasport1)
		} else {
			args = append(args, nil)
		}
		argId++
	}
	if car.PhotoTexpasport2 != nil {
		setValues = append(setValues, fmt.Sprintf("photo_texpasport2=$%d", argId))
		if *car.PhotoTexpasport2 != "" {
			args = append(args, car.PhotoTexpasport2)
		} else {
			args = append(args, nil)
		}
		argId++
	}
	if car.CarFront != nil {
		setValues = append(setValues, fmt.Sprintf("car_front=$%d", argId))
		if *car.CarFront != "" {
			args = append(args, car.CarFront)
		} else {
			args = append(args, nil)
		}
		argId++
	}
	if car.CarBack != nil {
		setValues = append(setValues, fmt.Sprintf("car_back=$%d", argId))
		if *car.CarBack != "" {
			args = append(args, car.CarBack)
		} else {
			args = append(args, nil)
		}
		argId++
	}
	if car.CarLeft != nil {
		setValues = append(setValues, fmt.Sprintf("car_left=$%d", argId))
		if *car.CarLeft != "" {
			args = append(args, car.CarLeft)
		} else {
			args = append(args, nil)
		}
		argId++
	}
	if car.CarRight != nil {
		setValues = append(setValues, fmt.Sprintf("car_right=$%d", argId))
		if *car.CarRight != "" {
			args = append(args, car.CarRight)
		} else {
			args = append(args, nil)
		}
		argId++
	}
	if car.CarFrontRow != nil {
		setValues = append(setValues, fmt.Sprintf("car_front_row=$%d", argId))
		if *car.CarFrontRow != "" {
			args = append(args, car.CarFrontRow)
		} else {
			args = append(args, nil)
		}
		argId++
	}
	if car.CarFrontBack != nil {
		setValues = append(setValues, fmt.Sprintf("car_front_back=$%d", argId))
		if *car.CarFrontBack != "" {
			args = append(args, car.CarFrontBack)
		} else {
			args = append(args, nil)
		}
		argId++
	}
	if car.CarBaggage != nil {
		setValues = append(setValues, fmt.Sprintf("car_baggage=$%d", argId))
		if *car.CarBaggage != "" {
			args = append(args, car.CarBaggage)
		} else {
			args = append(args, nil)
		}
		argId++
	}
	setQuery := strings.Join(setValues, ", ")
	updateQuery := fmt.Sprintf(`UPDATE %s SET %s WHERE driver_id = $%d`,
		driverCarTable, setQuery, argId)
	args = append(args, *usr.DriverId)
	_, err = tx.Exec(updateQuery, args...)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (r *AuthPostgres) GetDriverCarInfo(langId, userId int) (models.DriverCarInfo, error) {
	var usr models.User
	usrQuery := fmt.Sprintf("SELECT driver_id FROM %s WHERE id=$1", usersTable)
	err := r.db.Get(&usr, usrQuery, userId)
	if err != nil {
		return models.DriverCarInfo{}, err
	}
	var car models.DriverCarInfo
	carQuery := fmt.Sprintf("SELECT cc.name as car_color_name,cmk.name as car_marka_name,cmd.name as car_model_name FROM %s dc INNER JOIN %s cc ON cc.color_id = dc.car_color_id INNER JOIN %s cmk ON cmk.id = dc.car_marka_id INNER JOIN %s cmd ON cmd.id = dc.car_model_id WHERE dc.driver_id=$1 AND cc.language_id=$2", driverCarTable, colorsLangTable, carMarkasTable, carModelsTable)
	err = r.dash.Get(&car, carQuery, *usr.DriverId, langId)
	if err != nil {
		return models.DriverCarInfo{}, err
	}
	return car, nil
}

func (r *AuthPostgres) GetDriverInfo(langId, driverId int) (models.Driver, models.DriverCar, models.DriverCarInfo, error) {
	var usr models.User
	usrQuery := fmt.Sprintf("SELECT driver_id FROM %s WHERE id=$1", usersTable)
	err := r.db.Get(&usr, usrQuery, driverId)
	if err != nil {
		return models.Driver{}, models.DriverCar{}, models.DriverCarInfo{}, err
	}
	var driver models.Driver
	driverQuery := fmt.Sprintf("SELECT name,surname,photo,phone,rating FROM %s WHERE id=$1", driverTable)
	err = r.dash.Get(&driver, driverQuery, *usr.DriverId)
	if err != nil {
		return models.Driver{}, models.DriverCar{}, models.DriverCarInfo{}, err
	}
	if driver.Photo != nil {
		driver.Photo = utils.GetSmallFileUrl(strings.Split(*driver.Photo, "/"))
	}
	var carInfo models.DriverCarInfo
	carInfoQuery := fmt.Sprintf("SELECT cc.name as car_color_name,cmk.name as car_marka_name,cmd.name as car_model_name FROM %s dc INNER JOIN %s cc ON cc.color_id = dc.car_color_id INNER JOIN %s cmk ON cmk.id = dc.car_marka_id INNER JOIN %s cmd ON cmd.id = dc.car_model_id WHERE dc.driver_id=$1 AND cc.language_id=$2", driverCarTable, colorsLangTable, carMarkasTable, carModelsTable)
	err = r.dash.Get(&carInfo, carInfoQuery, *usr.DriverId, langId)
	if err != nil {
		return models.Driver{}, models.DriverCar{}, models.DriverCarInfo{}, err
	}
	var car models.DriverCar
	carQuery := fmt.Sprintf("SELECT car_number FROM %s WHERE driver_id=$1", driverCarTable)
	err = r.dash.Get(&car, carQuery, *usr.DriverId)
	if err != nil {
		return models.Driver{}, models.DriverCar{}, models.DriverCarInfo{}, err
	}
	return driver, car, carInfo, nil
}
