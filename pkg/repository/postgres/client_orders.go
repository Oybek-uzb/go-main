package postgres

import (
	"abir/models"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"strconv"
)

type ClientOrdersPostgres struct {
	db *sqlx.DB
}

func NewClientOrdersPostgres(db *sqlx.DB) *ClientOrdersPostgres {
	return &ClientOrdersPostgres{db: db}
}


func (r *ClientOrdersPostgres) RideList(ride models.Ride, langId, page int) ([]models.ClientRideList, models.Pagination, error){
	limit, err := strconv.Atoi(viper.GetString("vars.items_limit"))
	if err != nil {
		return []models.ClientRideList{}, models.Pagination{}, err
	}
	offset := limit * (page - 1)
	var pagination models.Pagination
	paginationQuery := fmt.Sprintf("SELECT count(*) AS total, $1 as current_page, CEIL(count(*)::decimal/$2) as last_page,$2 as per_page FROM %s WHERE departure_date::date = $3 AND departure_date::timestamp > NOW() AND from_district_id = $4 AND to_district_id = $5 AND status = $6",ridesTable)
	err = r.db.Get(&pagination, paginationQuery, page, limit, ride.DepartureDate,ride.FromDistrictId, ride.ToDistrictId,"new")
	if err != nil {
		return []models.ClientRideList{}, models.Pagination{}, err
	}
	var lists []models.ClientRideList
	listQuery := fmt.Sprintf("SELECT id as ride_id,driver_id,from_district_id,to_district_id," +
		"CASE WHEN from_district_id = 0 THEN (select COALESCE(dl.name, d.name) as name FROM dashboard.dictionary_region d LEFT JOIN dashboard.dictionary_region_i18n dl on d.id = dl.region_id WHERE d.id = $6 AND dl.language_id = $5) ELSE (select COALESCE(dl.name, d.name) as name FROM dashboard.dictionary_district d LEFT JOIN dashboard.dictionary_district_i18n dl on d.id = dl.district_id WHERE d.id = from_district_id AND dl.language_id = $5) END as from_district, " +
		"CASE WHEN to_district_id = 0 THEN (select COALESCE(dl.name, d.name) as name FROM dashboard.dictionary_region d LEFT JOIN dashboard.dictionary_region_i18n dl on d.id = dl.region_id WHERE d.id = $6 AND dl.language_id = $5) ELSE (select COALESCE(dl.name, d.name) as name FROM dashboard.dictionary_district d LEFT JOIN dashboard.dictionary_district_i18n dl on d.id = dl.district_id WHERE d.id = to_district_id AND dl.language_id = $5) END as to_district, " +
		"to_char(departure_date, 'HH24:MI') as departure_time,price,passenger_count,comments,status FROM %s WHERE departure_date::date = $1 AND departure_date::timestamp > NOW() AND from_district_id = $2 AND to_district_id = $3 AND status = $4 ORDER BY departure_date DESC LIMIT $7 OFFSET $8",ridesTable)
	err = r.db.Select(&lists, listQuery, ride.DepartureDate, ride.FromDistrictId, ride.ToDistrictId, "new", langId, viper.GetString("vars.capital_id"), limit, offset)
	return lists, pagination, err
}

func (r *ClientOrdersPostgres) RideSingle(langId, id, userId int) (models.ClientRideList, error){
	tx, err := r.db.Beginx()
	if err != nil {
		return models.ClientRideList{}, err
	}
	var list models.ClientRideList
	listQuery := fmt.Sprintf("SELECT id as ride_id,driver_id,from_district_id,to_district_id," +
		"CASE WHEN from_district_id = 0 THEN (select COALESCE(dl.name, d.name) as name FROM dashboard.dictionary_region d LEFT JOIN dashboard.dictionary_region_i18n dl on d.id = dl.region_id WHERE d.id = $2 AND dl.language_id = $1) ELSE (select COALESCE(dl.name, d.name) as name FROM dashboard.dictionary_district d LEFT JOIN dashboard.dictionary_district_i18n dl on d.id = dl.district_id WHERE d.id = from_district_id AND dl.language_id = $1) END as from_district, " +
		"CASE WHEN to_district_id = 0 THEN (select COALESCE(dl.name, d.name) as name FROM dashboard.dictionary_region d LEFT JOIN dashboard.dictionary_region_i18n dl on d.id = dl.region_id WHERE d.id = $2 AND dl.language_id = $1) ELSE (select COALESCE(dl.name, d.name) as name FROM dashboard.dictionary_district d LEFT JOIN dashboard.dictionary_district_i18n dl on d.id = dl.district_id WHERE d.id = to_district_id AND dl.language_id = $1) END as to_district, " +
		"departure_date as departure_time,price,passenger_count,comments,status FROM %s WHERE id=$3",ridesTable)
	err = tx.Get(&list, listQuery, langId, viper.GetString("vars.capital_id"), id)
	if err != nil {
		return models.ClientRideList{}, err
	}
	viewQuery := fmt.Sprintf("INSERT INTO %[1]s (ride_id, user_id) SELECT $1,$2 WHERE NOT EXISTS (SELECT id FROM %[1]s WHERE ride_id = $1 AND user_id = $2)", rideViewCountsTable)
	view, err := tx.Exec(viewQuery, list.RideId, userId)
	if err != nil {
		tx.Rollback()
		return models.ClientRideList{}, err
	}
	inserted, err := view.RowsAffected()
	if err != nil {
		tx.Rollback()
		return models.ClientRideList{}, err
	}
	if inserted != 0 {
		viewUpdateQuery := fmt.Sprintf("UPDATE %s SET view_count=view_count + 1 WHERE id=$1", ridesTable)
		_, err = r.db.Exec(viewUpdateQuery, list.RideId)
		if err != nil {
			tx.Rollback()
			return models.ClientRideList{}, err
		}
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return models.ClientRideList{}, err
	}
	return list, nil
}

func (r *ClientOrdersPostgres)  RideSingleStatus(rideId, userId int) (models.InterregionalOrder, error){
	var list models.InterregionalOrder
	listQuery := fmt.Sprintf("SELECT o.id,o.order_status,io.passenger_count,io.comments,o.created_at FROM %[1]s o LEFT JOIN %[2]s io ON o.order_id = io.id WHERE o.client_id=$1 AND io.ride_id=$2",ordersTable, interregionalOrdersTable)
	err := r.db.Get(&list, listQuery, userId, rideId)
	return list, err
}

func (r *ClientOrdersPostgres) RideSingleBook(bookRide models.Ride, rideId, userId int) (int, error){
	tx, err := r.db.Beginx()
	if err != nil {
		return 0, err
	}
	var ride models.Ride
	rideQuery := fmt.Sprintf("SELECT from_district_id,to_district_id,price,departure_date FROM %s WHERE id=$1", ridesTable)
	err = tx.Get(&ride, rideQuery, rideId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	var orderCnt int
	checkQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s LEFT JOIN %s ON orders.order_id = interregional_orders.id WHERE interregional_orders.ride_id=$1 AND orders.client_id=$2", ordersTable, interregionalOrdersTable)
	checkRow := tx.QueryRow(checkQuery, rideId, userId)
	if err := checkRow.Scan(&orderCnt); err != nil {
		return 0, err
	}
	if orderCnt > 0 {
		tx.Rollback()
		return 0, errors.New("you cannot order more than once")
	}
	price, err := strconv.Atoi(ride.Price)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	passengerCnt, err := strconv.Atoi(bookRide.PassengerCount)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	var subOrderId int
	query := fmt.Sprintf("INSERT INTO %s (ride_id, from_district_id,to_district_id,price,passenger_count,departure_date,comments) SELECT $1,$2,$3,$4,$5,$6,$7 RETURNING id", interregionalOrdersTable)
	row := tx.QueryRow(query, rideId, ride.FromDistrictId, ride.ToDistrictId, price*passengerCnt, bookRide.PassengerCount, ride.DepartureDate, bookRide.Comments)
	if err := row.Scan(&subOrderId); err != nil {
		return 0, err
	}
	subQuery := fmt.Sprintf("INSERT INTO %s (client_id, order_id, order_type) SELECT $1,$2,$3 RETURNING id", ordersTable)
	var orderId int
	subRow := tx.QueryRow(subQuery, userId, subOrderId, orderInterregionalType)
	if err := subRow.Scan(&orderId); err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return 0, err
	}
	return orderId, nil
}

type CityOrderPoint struct {
	Address string `json:"address"`
	Loc string `json:"loc"`
}
type CityOrderPoints struct{
	Points []CityOrderPoint `json:"points"`
}

func (r *ClientOrdersPostgres) Activity(userId int, page int, activityType string) ([]models.Activity, models.Pagination, error){
	limit, err := strconv.Atoi(viper.GetString("vars.items_limit"))
	if err != nil {
		return []models.Activity{}, models.Pagination{}, err
	}
	offset := limit * (page - 1)
	var lists []models.Activity
	var interregionalOrder models.InterregionalOrder
	var cityOrder models.CityOrder
	var cityOrderPoints CityOrderPoints
	var query string
	var pagination models.Pagination
	switch activityType {
		case "active":
			query = fmt.Sprintf("SELECT id as order_id, order_id as sub_order_id, order_type,order_status as status,created_at as order_time FROM %s WHERE client_id=$1 AND order_status IN('new', 'driver_accepted', 'driver_arrived', 'trip_started') ORDER BY id DESC LIMIT $2 OFFSET $3", ordersTable)
			err = r.db.Select(&lists, query, userId, limit, offset)
			break
		case "recently-completed":
			query = fmt.Sprintf("SELECT id as order_id, order_id as sub_order_id, order_type,order_status as status,created_at as order_time FROM %s WHERE client_id=$1 AND order_status IN('client_cancelled','driver_cancelled','order_completed') ORDER BY id DESC LIMIT 2", ordersTable)
			err = r.db.Select(&lists, query, userId)
			break
		case "history":
			query = fmt.Sprintf("SELECT id as order_id, order_id as sub_order_id, order_type,order_status as status,created_at as order_time FROM %s WHERE client_id=$1 AND order_status IN('client_cancelled','driver_cancelled','order_completed') ORDER BY id DESC LIMIT $2 OFFSET $3", ordersTable)
			err = r.db.Select(&lists, query, userId, limit, offset)
			paginationQuery := fmt.Sprintf("SELECT count(*) AS total, $1 as current_page, CEIL(count(*)::decimal/$2) as last_page,$2 as per_page FROM %s WHERE client_id=$3 AND order_status IN('client_cancelled','driver_cancelled','order_completed')",ordersTable)
			err = r.db.Get(&pagination, paginationQuery, page, limit, userId)
			break
	}
	for i, list := range lists {
		if list.OrderType == orderCityType {
			subQuery := fmt.Sprintf("SELECT points, tariff_id FROM %s WHERE id=$1", cityOrdersTable)
			subErr := r.db.Get(&cityOrder, subQuery, list.SubOrderId)
			if subErr != nil {
				if subErr == sql.ErrNoRows {
					continue
				}
				return []models.Activity{}, models.Pagination{}, subErr
			}
			err = json.Unmarshal([]byte(cityOrder.Points), &cityOrderPoints)
			if err != nil {
				continue
			}
			lists[i].From = cityOrderPoints.Points[0].Address
			lists[i].TariffId = &cityOrder.TariffId
			if len(cityOrderPoints.Points) > 1{
				lists[i].To = &cityOrderPoints.Points[len(cityOrderPoints.Points) - 1].Address
			}
		}else if list.OrderType == orderInterregionalType {
			subQuery := fmt.Sprintf("SELECT ride_id, from_district_id, to_district_id FROM %s WHERE id=$1", interregionalOrdersTable)
			subErr := r.db.Get(&interregionalOrder, subQuery, list.SubOrderId)
			if subErr != nil {
				if subErr == sql.ErrNoRows {
					continue
				}
			}
			lists[i].RideId = interregionalOrder.RideId
			lists[i].From = interregionalOrder.FromDistrictId
			lists[i].To = &interregionalOrder.ToDistrictId
		}
	}
	return lists, pagination, err
}

func (r *ClientOrdersPostgres) RideSingleCancel(cancelRide models.CanceledOrders, rideId, orderId, userId int) error{
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	viewUpdateQuery := fmt.Sprintf("UPDATE %s SET order_status='client_cancelled' WHERE id=$1 AND order_status NOT IN ('driver_cancelled', 'trip_started', 'order_completed')", ordersTable)
	_, err = tx.Exec(viewUpdateQuery, orderId)
	if err != nil {
		tx.Rollback()
		return err
	}
	if cancelRide.ReasonId != nil{
		comment := ""
		if cancelRide.Comments != nil{
			comment = *cancelRide.Comments
		}
		query := fmt.Sprintf("INSERT INTO %s (order_type, user_type,user_id,order_id,reason_id,comments) SELECT $1,$2,$3,$4,$5,$6", canceledOrdersTable)
		_, err = tx.Exec(query, "interregional", "client", userId, orderId, *cancelRide.ReasonId, comment)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (r *ClientOrdersPostgres) ChatFetch(userId, rideId,orderId int) ([]models.ChatMessages, error){
	var lists []models.ChatMessages
	listQuery := fmt.Sprintf("SELECT user_type,driver_id,client_id,ride_id,order_id,message_type,content,created_at FROM %s WHERE client_id=$1 AND ride_id=$2 AND order_id=$3 ORDER BY id DESC", chatMessagesTable)
	err := r.db.Select(&lists, listQuery, userId, rideId, orderId)
	return lists, err
}