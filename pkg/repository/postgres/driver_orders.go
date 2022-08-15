package postgres

import (
	"abir/models"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type DriverOrdersPostgres struct {
	db   *sqlx.DB
	dash *sqlx.DB
}

func NewDriverOrdersPostgres(db *sqlx.DB, dash *sqlx.DB) *DriverOrdersPostgres {
	return &DriverOrdersPostgres{db: db, dash: dash}
}

func (r *DriverOrdersPostgres) CreateRide(ride models.Ride, userId int) (int, error) {
	var id int
	createQuery := fmt.Sprintf("INSERT INTO %s (driver_id, from_district_id, to_district_id, departure_date, price, passenger_count, comments, status) SELECT $1,$2,$3,$4,$5,$6,$7,$8 WHERE NOT EXISTS (SELECT id FROM %s WHERE status = $8 AND driver_id = $1) RETURNING id", ridesTable, ridesTable)
	row := r.db.QueryRow(createQuery, userId, ride.FromDistrictId, ride.ToDistrictId, ride.DepartureDate, ride.Price, ride.PassengerCount, ride.Comments, "new")
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("you already have a new trip, cancel or end the trip")
		}
		return 0, err
	}
	return id, nil
}

func (r *DriverOrdersPostgres) RideList(userId int) ([]models.Ride, error) {
	var lists []models.Ride
	listQuery := fmt.Sprintf("SELECT id,driver_id,from_district_id,to_district_id,departure_date,price,passenger_count,comments,status,created_at,view_count FROM %s WHERE driver_id=$1 AND status IN ('cancelled','done') ORDER BY id DESC", ridesTable)
	err := r.db.Select(&lists, listQuery, userId)
	return lists, err
}

func (r *DriverOrdersPostgres) RideSingleActive(userId int) (models.Ride, error) {
	var list models.Ride
	listQuery := fmt.Sprintf("SELECT id,from_district_id,to_district_id,departure_date,price,passenger_count,comments,status,view_count,created_at FROM %s WHERE driver_id=$1 AND status IN ('new','on_the_way') ORDER BY id DESC LIMIT 1", ridesTable)
	err := r.db.Get(&list, listQuery, userId)
	return list, err
}

func (r *DriverOrdersPostgres) RideSingle(id, userId int) (models.Ride, error) {
	var list models.Ride
	listQuery := fmt.Sprintf("SELECT id,from_district_id,to_district_id,departure_date,price,passenger_count,comments,status,view_count,created_at FROM %s WHERE driver_id=$1 AND id=$2", ridesTable)
	err := r.db.Get(&list, listQuery, userId, id)
	if err != nil {
		return models.Ride{}, err
	}
	return list, nil
}

func (r *DriverOrdersPostgres) RideSingleNotifications(id, userId int) ([]models.RideNotification, error) {
	var lists []models.RideNotification
	listQuery := fmt.Sprintf("(SELECT cln.name, o.client_id as client_id, o.id as order_id, 'order' as type, o.created_at FROM %[1]s o LEFT JOIN %[2]s io ON o.order_id = io.id LEFT JOIN %[4]s usr ON o.client_id = usr.id LEFT JOIN %[5]s cln ON usr.client_id = cln.id WHERE io.ride_id=$1 AND o.order_type='interregional' AND o.order_status IN ('new')) UNION ALL (SELECT DISTINCT ON (ch.client_id) cln.name, ch.client_id as client_id, ch.order_id,'message' as type,ch.created_at FROM %[3]s ch LEFT JOIN %[4]s usr ON ch.client_id = usr.id LEFT JOIN %[5]s cln ON usr.client_id = cln.id WHERE ch.ride_id=$1 AND ch.user_type=$2 ORDER BY ch.client_id, ch.created_at DESC) ORDER BY created_at DESC", ordersTable, interregionalOrdersTable, chatMessagesTable, usersTable, clientsTable)
	err := r.db.Select(&lists, listQuery, id, clientType)
	return lists, err
}

func (r *DriverOrdersPostgres) RideSingleOrderList(id int) ([]models.InterregionalOrder, error) {
	var lists []models.InterregionalOrder
	listQuery := fmt.Sprintf("SELECT o.id, o.client_id, o.order_status, io.comments, io.passenger_count, o.created_at FROM %s o LEFT JOIN %s io ON o.order_id = io.id WHERE io.ride_id=$1 AND o.order_status NOT IN ('client_cancelled', 'driver_cancelled')", ordersTable, interregionalOrdersTable)
	err := r.db.Select(&lists, listQuery, id)
	if err != nil {
		return []models.InterregionalOrder{}, err
	}
	return lists, nil
}

func (r *DriverOrdersPostgres) RideSingleOrderView(orderId int) (models.InterregionalOrder, error) {
	var list models.InterregionalOrder
	listQuery := fmt.Sprintf("SELECT o.id, o.client_id, o.order_status, io.comments, io.passenger_count, o.created_at FROM %s o LEFT JOIN %s io ON o.order_id = io.id WHERE o.id=$1", ordersTable, interregionalOrdersTable)
	err := r.db.Get(&list, listQuery, orderId)
	if err != nil {
		return models.InterregionalOrder{}, err
	}
	return list, nil
}

func (r *DriverOrdersPostgres) RideSingleOrderAccept(driverId, orderId int) error {
	updateQuery := fmt.Sprintf(`UPDATE %s SET order_status='driver_accepted', driver_id=$1 WHERE id = $2 AND order_status='new'`, ordersTable)
	res, err := r.db.Exec(updateQuery, driverId, orderId)
	if err != nil {
		return err
	}
	updated, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if updated == 0 {
		return errors.New("order status has changed")
	}
	return nil
}

func (r *DriverOrdersPostgres) RideSingleOrderCancel(driverId, orderId int) error {
	updateQuery := fmt.Sprintf(`UPDATE %s SET order_status='driver_cancelled', driver_id=$1 WHERE id = $2 AND order_status NOT IN ('order_completed', 'client_cancelled')`, ordersTable)
	res, err := r.db.Exec(updateQuery, driverId, orderId)
	if err != nil {
		return err
	}
	updated, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if updated == 0 {
		return errors.New("order status has changed")
	}
	return nil
}

func (r *DriverOrdersPostgres) UpdateRide(ride models.Ride, id, userId int) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1
	if ride.PassengerCount != "" {
		setValues = append(setValues, fmt.Sprintf("passenger_count=$%d", argId))
		args = append(args, ride.PassengerCount)
		argId++
	}
	if ride.Price != "" {
		setValues = append(setValues, fmt.Sprintf("price=$%d", argId))
		args = append(args, ride.Price)
		argId++
	}
	setValues = append(setValues, fmt.Sprintf("comments=$%d", argId))
	if ride.Comments != nil {
		args = append(args, *ride.Comments)
	} else {
		args = append(args, nil)
	}
	argId++
	setQuery := strings.Join(setValues, ", ")
	updateQuery := fmt.Sprintf(`UPDATE %s SET %s WHERE id = $%d AND driver_id = $%d`,
		ridesTable, setQuery, argId, argId+1)
	args = append(args, id, userId)
	_, err := r.db.Exec(updateQuery, args...)
	return err
}

func (r *DriverOrdersPostgres) ChangeRideStatus(id, userId int, status string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	rideStatus := ""
	orderStatus := "driver_accepted"
	if status == "on_the_way" {
		rideStatus = "('new')"
		orderStatus = "trip_started"
	}
	if status == "done" {
		rideStatus = "('on_the_way')"
		orderStatus = "order_completed"
	}
	if status == "cancelled" {
		rideStatus = "('new','on_the_way')"
		orderStatus = "driver_cancelled"
	}
	updateQuery := fmt.Sprintf(`UPDATE %s SET status=$1 WHERE id = $2 AND driver_id = $3 AND status IN %s`,
		ridesTable, rideStatus)
	res, err := tx.Exec(updateQuery, status, id, userId)
	if err != nil {
		tx.Rollback()
		return err
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if cnt == 0 {
		tx.Rollback()
		return errors.New("you can't change the status")
	}
	updateOrdersQuery := fmt.Sprintf(`UPDATE %s SET order_status=$1 FROM %s io WHERE order_id = io.id AND io.ride_id = $2 AND order_status NOT IN ('client_cancelled')`, ordersTable, interregionalOrdersTable)
	_, err = tx.Exec(updateOrdersQuery, orderStatus, id)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (r *DriverOrdersPostgres) ChatFetch(userId, rideId, orderId int) ([]models.ChatMessages, error) {
	var lists []models.ChatMessages
	listQuery := fmt.Sprintf("SELECT user_type,driver_id,client_id,ride_id,order_id,message_type,content,created_at FROM %s WHERE driver_id=$1 AND ride_id=$2 AND order_id=$3 ORDER BY id DESC", chatMessagesTable)
	err := r.db.Select(&lists, listQuery, userId, rideId, orderId)
	return lists, err
}
func (r *DriverOrdersPostgres) CityOrderChangeStatus(req models.CityOrderRequest, cancelOrRate models.CancelOrRateReasons, orderId, userId int, status string) (int, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return 0, err
	}
	orderStatus := ""
	viewUpdateQuery := ""
	uType := ""
	if status == "driver_accepted" {
		orderStatus = "('new')"
		viewUpdateQuery = fmt.Sprintf("UPDATE %s SET order_status=$3,driver_id=$2 WHERE id=$1 AND order_status IN %s", ordersTable, orderStatus)
		uType = "on_the_way"
	}
	if status == "driver_arrived" {
		orderStatus = "('driver_accepted')"
	}
	if status == "trip_started" {
		orderStatus = "('driver_arrived', 'client_going_out')"
	}
	if status == "driver_cancelled" {
		orderStatus = "('driver_accepted', 'driver_arrived', 'client_going_out')"
		uType = "online"
	}
	if status == "order_completed" {
		orderStatus = "('order_completed', 'trip_started')"
		uType = "online"
	}
	if status == "driver_arrived" || status == "trip_started" || status == "driver_cancelled" || status == "order_completed" {
		viewUpdateQuery = fmt.Sprintf("UPDATE %s SET order_status=$3, updated_at=NOW() WHERE id=$1 AND driver_id=$2 AND order_status IN %s", ordersTable, orderStatus)
	}
	if viewUpdateQuery != "" {
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
	}
	var order models.Order
	orderQuery := fmt.Sprintf("SELECT order_id,client_id FROM %s WHERE id=$1", ordersTable)
	err = tx.Get(&order, orderQuery, orderId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if uType != "" {
		updateStatusQuery := fmt.Sprintf("UPDATE %s SET user_id=$1,driver_status=$2 WHERE user_id=$1", driverStatusesTable)
		_, err = tx.Exec(updateStatusQuery, userId, uType)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	if order.ClientId == 0 {
		return 0, errors.New("client not found")
	}
	if status == "order_completed" {
		var rideInfoQuery string
		rideInfo, err := json.Marshal(req)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		if req.DriverLastLocation != nil && req.DriverLastAddress != nil {
			var subOrder models.CityOrder
			subOrderQuery := fmt.Sprintf("SELECT points FROM %s WHERE id=$1", cityOrdersTable)
			err = tx.Get(&subOrder, subOrderQuery, order.OrderId)
			if err != nil {
				tx.Rollback()
				return 0, err
			}
			var pointsArr models.CityOrderPoints
			err = json.Unmarshal([]byte(subOrder.Points), &pointsArr)
			if err != nil {
				return 0, err
			}
			newPointsArrPoints := pointsArr.Points
			newPointsArrPoints = append(newPointsArrPoints, models.PointsArr{Location: *req.DriverLastLocation, Address: *req.DriverLastAddress})
			pointsArrJson, err := json.Marshal(models.CityOrderPoints{Distance: pointsArr.Distance, Points: newPointsArrPoints})
			if err != nil {
				tx.Rollback()
				return 0, err
			}
			rideInfoQuery = fmt.Sprintf("UPDATE %s SET ride_info=$1, price=$3, points=$4 WHERE id=$2", cityOrdersTable)
			_, err = tx.Exec(rideInfoQuery, string(rideInfo), order.OrderId, req.OrderAmount, string(pointsArrJson))
			if err != nil {
				tx.Rollback()
				return 0, err
			}
		} else {
			rideInfoQuery = fmt.Sprintf("UPDATE %s SET ride_info=$1, price=$3 WHERE id=$2", cityOrdersTable)
			_, err = tx.Exec(rideInfoQuery, string(rideInfo), order.OrderId, req.OrderAmount)
			if err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}
	if status == "driver_cancelled" {
		if cancelOrRate.ReasonId != "" {
			comment := ""
			if cancelOrRate.Comments != nil {
				comment = *cancelOrRate.Comments
			}
			var cancelRideReasonId int
			query := fmt.Sprintf("INSERT INTO %s (order_type, user_type,user_id,order_id,comments) SELECT $1,$2,$3,$4,$5 RETURNING id", canceledOrdersTable)
			row := tx.QueryRow(query, "city", "driver", userId, orderId, comment)
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
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return 0, err
	}
	return order.ClientId, nil
}

func (r *DriverOrdersPostgres) CityOrderView(orderId, userId int) (models.CityOrder, error) {
	var order models.Order
	orderQuery := fmt.Sprintf("SELECT order_id,client_id,driver_id,order_status,updated_at as changed_at FROM %s WHERE id=$1 AND driver_id=$2", ordersTable)
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
	subOrder.ClientId = order.ClientId
	subOrder.OrderStatus = order.OrderStatus
	return subOrder, nil
}
func (r *DriverOrdersPostgres) CityTariff(districtId, tariffId int) (models.CityTariffs, error) {
	var list models.CityTariffs
	tariffsQuery := fmt.Sprintf("SELECT starting_price as start_price, per_kilometer as price_per_km, countryside as price_per_km_outer, conditioner as ac_price, expectation FROM %s WHERE tariff_id=$1 AND district_id=$2", routeCityTaxiTable)
	err := r.dash.Get(&list, tariffsQuery, tariffId, districtId)
	return list, err
}

func (r *DriverOrdersPostgres) CityTariffInfo(districtId, tariffId int) (models.TariffInfo, error) {
	var list models.TariffInfo
	tariffsQuery := fmt.Sprintf("SELECT starting_price as start_price, per_kilometer as price_per_km, countryside as price_per_km_outer, conditioner as ac_price, expectation FROM %s WHERE tariff_id=$1 AND district_id=$2", routeCityTaxiTable)
	err := r.dash.Get(&list, tariffsQuery, tariffId, districtId)
	return list, err
}
