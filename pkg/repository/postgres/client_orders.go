package postgres

import (
	"abir/models"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
	"strconv"
	"strings"
)

type ClientOrdersPostgres struct {
	db   *sqlx.DB
	dash *sqlx.DB
}

func NewClientOrdersPostgres(db *sqlx.DB, dash *sqlx.DB) *ClientOrdersPostgres {
	return &ClientOrdersPostgres{db: db, dash: dash}
}

func (r *ClientOrdersPostgres) RideList(ride models.Ride, langId, page int) ([]models.ClientRideList, models.Pagination, error) {
	limit, err := strconv.Atoi(viper.GetString("vars.items_limit"))
	if err != nil {
		return []models.ClientRideList{}, models.Pagination{}, err
	}
	offset := limit * (page - 1)
	var pagination models.Pagination
	paginationQuery := fmt.Sprintf("SELECT count(*) AS total, $1 as current_page, CEIL(count(*)::decimal/$2) as last_page,$2 as per_page FROM %s WHERE departure_date::date = $3 AND departure_date::timestamp > NOW() AND from_district_id = $4 AND to_district_id = $5 AND status = $6", ridesTable)
	err = r.db.Get(&pagination, paginationQuery, page, limit, ride.DepartureDate, ride.FromDistrictId, ride.ToDistrictId, "new")
	if err != nil {
		return []models.ClientRideList{}, models.Pagination{}, err
	}
	var lists []models.ClientRideList
	listQuery := fmt.Sprintf("SELECT id as ride_id,driver_id,from_district_id,to_district_id,"+
		"CASE WHEN from_district_id = 0 THEN (select COALESCE(dl.name, d.name) as name FROM dashboard.dictionary_region d LEFT JOIN dashboard.dictionary_region_i18n dl on d.id = dl.region_id WHERE d.id = $6 AND dl.language_id = $5) ELSE (select COALESCE(dl.name, d.name) as name FROM dashboard.dictionary_district d LEFT JOIN dashboard.dictionary_district_i18n dl on d.id = dl.district_id WHERE d.id = from_district_id AND dl.language_id = $5) END as from_district, "+
		"CASE WHEN to_district_id = 0 THEN (select COALESCE(dl.name, d.name) as name FROM dashboard.dictionary_region d LEFT JOIN dashboard.dictionary_region_i18n dl on d.id = dl.region_id WHERE d.id = $6 AND dl.language_id = $5) ELSE (select COALESCE(dl.name, d.name) as name FROM dashboard.dictionary_district d LEFT JOIN dashboard.dictionary_district_i18n dl on d.id = dl.district_id WHERE d.id = to_district_id AND dl.language_id = $5) END as to_district, "+
		"to_char(departure_date, 'HH24:MI') as departure_time,price,passenger_count,comments,status FROM %s WHERE departure_date::date = $1 AND departure_date::timestamp > NOW() AND from_district_id = $2 AND to_district_id = $3 AND status = $4 ORDER BY departure_date DESC LIMIT $7 OFFSET $8", ridesTable)
	err = r.db.Select(&lists, listQuery, ride.DepartureDate, ride.FromDistrictId, ride.ToDistrictId, "new", langId, viper.GetString("vars.capital_id"), limit, offset)
	return lists, pagination, err
}

func (r *ClientOrdersPostgres) RideSingle(langId, id, userId int) (models.ClientRideList, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return models.ClientRideList{}, err
	}
	var list models.ClientRideList
	listQuery := fmt.Sprintf("SELECT id as ride_id,driver_id,from_district_id,to_district_id,"+
		"CASE WHEN from_district_id = 0 THEN (select COALESCE(dl.name, d.name) as name FROM dashboard.dictionary_region d LEFT JOIN dashboard.dictionary_region_i18n dl on d.id = dl.region_id WHERE d.id = $2 AND dl.language_id = $1) ELSE (select COALESCE(dl.name, d.name) as name FROM dashboard.dictionary_district d LEFT JOIN dashboard.dictionary_district_i18n dl on d.id = dl.district_id WHERE d.id = from_district_id AND dl.language_id = $1) END as from_district, "+
		"CASE WHEN to_district_id = 0 THEN (select COALESCE(dl.name, d.name) as name FROM dashboard.dictionary_region d LEFT JOIN dashboard.dictionary_region_i18n dl on d.id = dl.region_id WHERE d.id = $2 AND dl.language_id = $1) ELSE (select COALESCE(dl.name, d.name) as name FROM dashboard.dictionary_district d LEFT JOIN dashboard.dictionary_district_i18n dl on d.id = dl.district_id WHERE d.id = to_district_id AND dl.language_id = $1) END as to_district, "+
		"departure_date as departure_time,price,passenger_count,comments,status FROM %s WHERE id=$3", ridesTable)
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

func (r *ClientOrdersPostgres) RideSingleStatus(rideId, userId int) (models.InterregionalOrder, error) {
	var list models.InterregionalOrder
	listQuery := fmt.Sprintf("SELECT o.id,o.order_status,io.passenger_count,io.comments,o.created_at FROM %[1]s o LEFT JOIN %[2]s io ON o.order_id = io.id WHERE o.client_id=$1 AND io.ride_id=$2", ordersTable, interregionalOrdersTable)
	err := r.db.Get(&list, listQuery, userId, rideId)
	return list, err
}

func (r *ClientOrdersPostgres) RideSingleBook(bookRide models.Ride, rideId, userId int) (int, error) {
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
	Loc     string `json:"loc"`
}
type CityOrderPoints struct {
	Points []CityOrderPoint `json:"points"`
}

func (r *ClientOrdersPostgres) Activity(userId int, page int, activityType, orderType string) ([]models.Activity, models.Pagination, error) {
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
	orderTypeQuery := ""
	switch orderType {
	case "city":
		orderTypeQuery = "AND order_type='city'"
		break
	case "interregional":
		orderTypeQuery = "AND order_type='interregional'"
		break
	}
	switch activityType {
	case "active":
		query = fmt.Sprintf("SELECT id as order_id, order_id as sub_order_id, order_type,order_status as status,created_at as order_time FROM %s WHERE client_id=$1 AND order_status IN('new', 'driver_accepted', 'driver_arrived', 'client_going_out', 'trip_started') AND (CASE WHEN order_status = 'new' THEN created_at > current_timestamp - (4 * interval '1 minute') ELSE true END) ORDER BY id DESC LIMIT $2 OFFSET $3", ordersTable)
		err = r.db.Select(&lists, query, userId, limit, offset)
		break
	case "recently-completed":
		query = fmt.Sprintf("SELECT id as order_id, order_id as sub_order_id, order_type,order_status as status,created_at as order_time FROM %s WHERE client_id=$1 AND order_status IN('client_cancelled','driver_cancelled','order_completed') ORDER BY id DESC LIMIT 2", ordersTable)
		err = r.db.Select(&lists, query, userId)
		break
	case "history":
		query = fmt.Sprintf("SELECT id as order_id, order_id as sub_order_id, order_type,order_status as status,created_at as order_time FROM %s WHERE client_id=$1 AND order_status IN('client_cancelled','driver_cancelled','order_completed') %s ORDER BY id DESC LIMIT $2 OFFSET $3", ordersTable, orderTypeQuery)
		err = r.db.Select(&lists, query, userId, limit, offset)
		paginationQuery := fmt.Sprintf("SELECT count(*) AS total, $1 as current_page, CEIL(count(*)::decimal/$2) as last_page,$2 as per_page FROM %s WHERE client_id=$3 AND order_status IN('client_cancelled','driver_cancelled','order_completed') %s", ordersTable, orderTypeQuery)
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
			if len(cityOrderPoints.Points) > 1 {
				lists[i].To = &cityOrderPoints.Points[len(cityOrderPoints.Points)-1].Address
			}
		} else if list.OrderType == orderInterregionalType {
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

func (r *ClientOrdersPostgres) RideSingleCancel(cancelRide models.CancelOrRateReasons, rideId, orderId, userId int) error {
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
	if cancelRide.ReasonId != "" {
		comment := ""
		if cancelRide.Comments != nil {
			comment = *cancelRide.Comments
		}
		var cancelRideReasonId int
		query := fmt.Sprintf("INSERT INTO %s (order_type, user_type,user_id,order_id,comments) SELECT $1,$2,$3,$4,$5 RETURNING id", canceledOrdersTable)
		row := tx.QueryRow(query, "interregional", "client", userId, orderId, comment)
		if err := row.Scan(&cancelRideReasonId); err != nil {
			tx.Rollback()
			return err
		}
		reasonIds := strings.Split(cancelRide.ReasonId, ",")
		insertValues := make([]string, 0)
		for _, reasonId := range reasonIds {
			insertValues = append(insertValues, fmt.Sprintf("(%v,%v)", cancelRideReasonId, reasonId))
		}
		batchQuery := fmt.Sprintf("INSERT INTO %s (canceled_order_id,reason_id) VALUES %s", canceledOrderReasonsTable, strings.Join(insertValues, ", "))
		_, err = tx.Exec(batchQuery)
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

func (r *ClientOrdersPostgres) ChatFetch(userId, rideId, orderId int) ([]models.ChatMessages, error) {
	var lists []models.ChatMessages
	listQuery := fmt.Sprintf("SELECT user_type,driver_id,client_id,ride_id,order_id,message_type,content,created_at FROM %s WHERE client_id=$1 AND ride_id=$2 AND order_id=$3 ORDER BY id DESC", chatMessagesTable)
	err := r.db.Select(&lists, listQuery, userId, rideId, orderId)
	return lists, err
}

func (r *ClientOrdersPostgres) CityTariffs(districtId, langId int) ([]models.CityTariffs, error) {
	var lists []models.CityTariffs
	tariffsQuery := fmt.Sprintf("SELECT t.id,rt.starting_price as start_price, rt.per_kilometer as price_per_km, rt.countryside as price_per_km_outer, rt.conditioner as ac_price, t.cars as cars, tl.name as tariff_name, tl.description as description, t.image as icon, t.additional as image, rt.expectation  FROM %s t LEFT JOIN %s tl ON t.id = tl.tariff_id LEFT JOIN %s rt ON t.id = rt.tariff_id WHERE tl.language_id=$1 AND rt.district_id=$2 ORDER BY array_position(ARRAY[1,2,3,13,4,5,6,7,8,9,10,11,12]::bigint[], t.id)", tariffsTable, tariffsLangTable, routeCityTaxiTable)
	err := r.dash.Select(&lists, tariffsQuery, langId, districtId)
	return lists, err
}

func (r *ClientOrdersPostgres) CityNewOrder(order models.CityOrder, userId int) (int, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return 0, err
	}
	var cardId any
	if order.CardId != nil {
		cardId = *order.CardId
	} else {
		cardId = nil
	}
	var forAnotherPhone any
	if order.ForAnotherPhone != nil {
		forAnotherPhone = *order.ForAnotherPhone
	} else {
		forAnotherPhone = nil
	}
	var receiverComments any
	if order.ReceiverComments != nil {
		receiverComments = *order.ReceiverComments
	} else {
		receiverComments = nil
	}
	var receiverPhone any
	if order.ReceiverPhone != nil {
		receiverPhone = *order.ReceiverPhone
	} else {
		receiverPhone = nil
	}
	var comments any
	if order.Comments != nil {
		comments = *order.Comments
	} else {
		comments = nil
	}
	var subOrderId int
	query := fmt.Sprintf("INSERT INTO %s (points,tariff_id,cargo_type,payment_type,card_id,has_conditioner,for_another,for_another_phone,receiver_comments,receiver_phone,price,comments) SELECT $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12 RETURNING id", cityOrdersTable)
	row := tx.QueryRow(query, order.Points, order.TariffId, order.CargoType, order.PaymentType, cardId, order.HasConditioner, order.ForAnother, forAnotherPhone, receiverComments, receiverPhone, order.Price, comments)
	if err := row.Scan(&subOrderId); err != nil {
		return 0, err
	}
	subQuery := fmt.Sprintf("INSERT INTO %s (client_id, order_id, order_type) SELECT $1,$2,$3 RETURNING id", ordersTable)
	var orderId int
	subRow := tx.QueryRow(subQuery, userId, subOrderId, orderCityType)
	if err := subRow.Scan(&orderId); err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return 0, err
	}
	return orderId, nil
}

func (r *ClientOrdersPostgres) CityOrderView(orderId, userId int) (models.CityOrder, error) {
	var order models.Order
	orderQuery := fmt.Sprintf("SELECT order_id,driver_id,order_status,updated_at as changed_at FROM %s WHERE id=$1 AND client_id=$2", ordersTable)
	err := r.db.Get(&order, orderQuery, orderId, userId)
	if err != nil {
		return models.CityOrder{}, err
	}
	var subOrder models.CityOrder
	subOrderQuery := fmt.Sprintf("SELECT points,tariff_id,cargo_type,payment_type,has_conditioner,for_another,for_another_phone,receiver_comments,receiver_phone,price,comments,created_at FROM %s WHERE id=$1", cityOrdersTable)
	err = r.db.Get(&subOrder, subOrderQuery, order.OrderId)
	if err != nil {
		return models.CityOrder{}, err
	}
	subOrder.Id = orderId
	subOrder.DriverId = order.DriverId
	subOrder.OrderStatus = order.OrderStatus
	subOrder.ChangedAt = order.ChangedAt
	return subOrder, nil
}

func (r *ClientOrdersPostgres) CityOrderChangeStatus(cancelOrRate models.CancelOrRateReasons, orderId, userId int, status string) (int, error) {
	var ord models.Order
	orderQuery := fmt.Sprintf("SELECT id,order_id,driver_id,order_status FROM %s WHERE id=$1", ordersTable)
	err := r.db.Get(&ord, orderQuery, orderId)
	if err != nil {
		return 0, err
	}
	if ord.DriverId == nil {
		if !slices.Contains([]string{"new", "client_cancelled"}, ord.OrderStatus) {
			return 0, errors.New("driver not found")
		}
	}
	tx, err := r.db.Beginx()
	if err != nil {
		return 0, err
	}
	orderStatus := ""
	if status == "client_cancelled" {
		orderStatus = "('new', 'driver_accepted', 'driver_arrived', 'trip_started')"
	}
	if status == "client_going_out" {
		orderStatus = "('driver_arrived')"
	}
	if status != "client_rate" {
		viewUpdateQuery := fmt.Sprintf("UPDATE %s SET order_status=$3 WHERE id=$1 AND client_id=$2 AND order_status IN %s", ordersTable, orderStatus)
		res, err := tx.Exec(viewUpdateQuery, orderId, userId, status)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		affected, err := res.RowsAffected()
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		if affected == 0 {
			tx.Rollback()
			return 0, errors.New("you can't change the status")
		}
	} else {
		if cancelOrRate.Rate != 0 {
			comment := ""
			if cancelOrRate.Comments != nil {
				comment = *cancelOrRate.Comments
			}
			var cancelRideReasonId int
			query := fmt.Sprintf("INSERT INTO %s (order_type,user_type,user_id,order_id,comments,rate) SELECT $1,$2,$3,$4,$5,$6 RETURNING id", ratedOrdersTable)
			row := tx.QueryRow(query, "city", "client", userId, orderId, comment, cancelOrRate.Rate)
			if err := row.Scan(&cancelRideReasonId); err != nil {
				tx.Rollback()
				return 0, err
			}
			if cancelOrRate.ReasonId != "" {
				reasonIds := strings.Split(cancelOrRate.ReasonId, ",")
				insertValues := make([]string, 0)
				for _, reasonId := range reasonIds {
					insertValues = append(insertValues, fmt.Sprintf("(%v,%v)", cancelRideReasonId, reasonId))
				}
				batchQuery := fmt.Sprintf("INSERT INTO %s (rated_order_id,reason_id) VALUES %s", ratedOrderReasonsTable, strings.Join(insertValues, ", "))
				_, err = tx.Exec(batchQuery)
				if err != nil {
					tx.Rollback()
					return 0, err
				}
			}
		}
	}
	if status == "client_cancelled" {
		if cancelOrRate.ReasonId != "" {
			comment := ""
			if cancelOrRate.Comments != nil {
				comment = *cancelOrRate.Comments
			}
			var cancelRideReasonId int
			query := fmt.Sprintf("INSERT INTO %s (order_type, user_type,user_id,order_id,comments) SELECT $1,$2,$3,$4,$5 RETURNING id", canceledOrdersTable)
			row := tx.QueryRow(query, "city", "client", userId, orderId, comment)
			if err := row.Scan(&cancelRideReasonId); err != nil {
				tx.Rollback()
				return 0, err
			}
			reasonIds := strings.Split(cancelOrRate.ReasonId, ",")
			insertValues := make([]string, 0)
			for _, reasonId := range reasonIds {
				insertValues = append(insertValues, fmt.Sprintf("(%v,%v)", cancelRideReasonId, reasonId))
			}
			batchQuery := fmt.Sprintf("INSERT INTO %s (canceled_order_id,reason_id) VALUES %s", canceledOrderReasonsTable, strings.Join(insertValues, ", "))
			_, err = tx.Exec(batchQuery)
			if err != nil {
				tx.Rollback()
				return 0, err
			}
		}
		if ord.DriverId != nil {
			updateStatusQuery := fmt.Sprintf("UPDATE %s SET user_id=$1,driver_status=$2 WHERE user_id=$1", driverStatusesTable)
			_, err = tx.Exec(updateStatusQuery, *ord.DriverId, "online")
			if err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return 0, err
	}
	if ord.DriverId == nil {
		return 0, nil
	}
	return *ord.DriverId, nil
}
