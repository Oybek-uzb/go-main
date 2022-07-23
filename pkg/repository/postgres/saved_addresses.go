package postgres

import (
	"abir/models"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type SavedAddressesPostgres struct {
	db *sqlx.DB
}

func NewSavedAddressesPostgres(db *sqlx.DB) *SavedAddressesPostgres {
	return &SavedAddressesPostgres{db: db}
}

func (r *SavedAddressesPostgres) Get(userId int) ([]models.SavedAddresses, error) {
	var lists []models.SavedAddresses
	query := fmt.Sprintf("SELECT id, name, place_type, location, address FROM %s WHERE user_id = $1 ORDER BY CASE place_type WHEN 'home' THEN 1 WHEN 'work' THEN 2 WHEN 'custom' THEN 3 END, id ASC", savedAddressesTable)
	err := r.db.Select(&lists, query, userId)
	return lists, err
}
func (r *SavedAddressesPostgres) Store(address models.SavedAddresses, userId int) error {
	query := fmt.Sprintf("INSERT INTO %s (name, place_type, location, address, user_id) SELECT $1,$2,$3,$4,$5 WHERE NOT EXISTS (SELECT id FROM %s WHERE place_type = $2 AND place_type != 'custom' AND user_id = $5)", savedAddressesTable, savedAddressesTable)
	_, err := r.db.Exec(query, address.Name, address.PlaceType, address.Location, address.Address, userId)
	return err
}
func (r *SavedAddressesPostgres) Update(address models.SavedAddresses, addressId, userId int) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1
	if address.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *address.Name)
		argId++
	}
	if address.PlaceType != nil {
		setValues = append(setValues, fmt.Sprintf("place_type=$%d", argId))
		args = append(args, *address.PlaceType)
		argId++
	}
	if address.Location != nil {
		setValues = append(setValues, fmt.Sprintf("location=$%d", argId))
		args = append(args, *address.Location)
		argId++
	}
	if address.Address != nil {
		setValues = append(setValues, fmt.Sprintf("address=$%d", argId))
		args = append(args, *address.Address)
		argId++
	}
	setQuery := strings.Join(setValues, ", ")
	updateQuery := fmt.Sprintf(`UPDATE %s SET %s WHERE id = $%d AND user_id = $%d`,
		savedAddressesTable, setQuery, argId, argId + 1)
	args = append(args, addressId, userId)
	_, err := r.db.Exec(updateQuery, args...)
	return err
}
func (r *SavedAddressesPostgres) Delete(addressId, userId int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1 AND user_id = $2`,
		savedAddressesTable)
	_, err := r.db.Exec(query, addressId, userId)
	return err
}
