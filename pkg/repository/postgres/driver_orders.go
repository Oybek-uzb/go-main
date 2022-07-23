package postgres

import (
	"abir/models"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type DriverOrdersPostgres struct {
	db *sqlx.DB
}

func NewDriverOrdersPostgres(db *sqlx.DB) *DriverOrdersPostgres {
	return &DriverOrdersPostgres{db: db}
}

func (r *DriverOrdersPostgres) CreateRide(ride models.Ride, userId int) (int, error){
	var id int
	createQuery := fmt.Sprintf("INSERT INTO %s (driver_id, from_district_id, to_district_id, departure_date, price, passenger_count, comments, status) SELECT $1,$2,$3,$4,$5,$6,$7,$8 WHERE NOT EXISTS (SELECT id FROM %s WHERE status = $8 AND driver_id = $1) RETURNING id", ridesTable, ridesTable)
	row := r.db.QueryRow(createQuery, userId, ride.FromDistrictId, ride.ToDistrictId, ride.DepartureDate, ride.Price, ride.PassengerCount, ride.Comments, "new")
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows{
			return 0, errors.New("you already have a new trip, cancel or end the trip")
		}
		return 0, err
	}
	return id, nil
}

func (r *DriverOrdersPostgres) RideList(userId int) ([]models.Ride, error){
	var lists []models.Ride
	listQuery := fmt.Sprintf("SELECT id,driver_id,from_district_id,to_district_id,departure_date,price,passenger_count,comments,status,created_at,view_count FROM %s WHERE driver_id=$1 AND status IN ('cancelled','done') ORDER BY id DESC", ridesTable)
	err := r.db.Select(&lists, listQuery, userId)
	return lists, err
}

func (r *DriverOrdersPostgres) RideSingleActive(userId int) (models.Ride, error){
	var list models.Ride
	listQuery := fmt.Sprintf("SELECT id,from_district_id,to_district_id,departure_date,price,passenger_count,comments,status,view_count,created_at FROM %s WHERE driver_id=$1 AND status IN ('new','on_the_way') ORDER BY id DESC LIMIT 1", ridesTable)
	err := r.db.Get(&list, listQuery, userId)
	return list, err
}


func (r *DriverOrdersPostgres) RideSingle(id, userId int) (models.Ride, error){
	var list models.Ride
	listQuery := fmt.Sprintf("SELECT id,from_district_id,to_district_id,departure_date,price,passenger_count,comments,status,view_count,created_at FROM %s WHERE driver_id=$1 AND id=$2",ridesTable)
	err := r.db.Get(&list, listQuery, userId, id)
	if err != nil {
		return models.Ride{}, err
	}
	return list, nil
}

func (r *DriverOrdersPostgres) RideSingleNotifications(id, userId int) ([]models.RideNotification, error){
	var lists []models.RideNotification
	listQuery := fmt.Sprintf("(SELECT cln.name, o.id as order_id, 'order' as type, o.created_at FROM %[1]s o LEFT JOIN %[2]s io ON o.order_id = io.id LEFT JOIN %[4]s usr ON o.client_id = usr.id LEFT JOIN %[5]s cln ON usr.client_id = cln.id WHERE io.ride_id=$1 AND o.order_status IN ('new')) UNION ALL (SELECT cln.name, ch.order_id,'message' as type,ch.created_at FROM %[3]s ch LEFT JOIN %[4]s usr ON ch.client_id = usr.id LEFT JOIN %[5]s cln ON usr.client_id = cln.id WHERE ch.ride_id=$1 AND ch.user_type=$2 ORDER BY ch.created_at DESC LIMIT 1) ORDER BY created_at DESC", ordersTable, interregionalOrdersTable, chatMessagesTable, usersTable, clientsTable)
	err := r.db.Select(&lists, listQuery, id, clientType)
	return lists, err
}

func (r *DriverOrdersPostgres) RideSingleOrderList(id int) ([]models.InterregionalOrder, error){
	var lists []models.InterregionalOrder
	listQuery := fmt.Sprintf("SELECT o.id, o.client_id, o.order_status, io.comments, io.passenger_count, o.created_at FROM %s o LEFT JOIN %s io ON o.order_id = io.id WHERE io.ride_id=$1 AND o.order_status NOT IN ('client_cancelled', 'driver_cancelled')",ordersTable,interregionalOrdersTable)
	err := r.db.Select(&lists, listQuery, id)
	if err != nil {
		return []models.InterregionalOrder{}, err
	}
	return lists, nil
}

func (r *DriverOrdersPostgres) RideSingleOrderView(orderId int) (models.InterregionalOrder, error){
	var list models.InterregionalOrder
	listQuery := fmt.Sprintf("SELECT o.id, o.client_id, o.order_status, io.comments, io.passenger_count, o.created_at FROM %s o LEFT JOIN %s io ON o.order_id = io.id WHERE o.id=$1",ordersTable,interregionalOrdersTable)
	err := r.db.Get(&list, listQuery, orderId)
	if err != nil {
		return models.InterregionalOrder{}, err
	}
	return list, nil
}

func (r *DriverOrdersPostgres) RideSingleOrderAccept(driverId, orderId int) error{
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

func (r *DriverOrdersPostgres) RideSingleOrderCancel(driverId, orderId int) error{
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
	if ride.Comments != nil{
		args = append(args, *ride.Comments)
	}else{
		args = append(args, nil)
	}
	argId++
	setQuery := strings.Join(setValues, ", ")
	updateQuery := fmt.Sprintf(`UPDATE %s SET %s WHERE id = $%d AND driver_id = $%d`,
		ridesTable, setQuery, argId, argId + 1)
	args = append(args, id, userId)
	_, err := r.db.Exec(updateQuery, args...)
	return err
}

func (r *DriverOrdersPostgres) ChangeRideStatus(id, userId int, status string) error{
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	rideStatus := ""
	orderStatus := "driver_accepted"
	if status == "on_the_way"{
		rideStatus = "('new')"
		orderStatus = "trip_started"
	}
	if status == "done"{
		rideStatus = "('on_the_way')"
		orderStatus = "order_completed"
	}
	if status == "cancelled"{
		rideStatus = "('new','on_the_way')"
		orderStatus = "driver_cancelled"
	}
	updateQuery := fmt.Sprintf(`UPDATE %s SET status=$1 WHERE id = $2 AND driver_id = $3 AND status IN %s`,
		ridesTable, rideStatus)
	res, err := tx.Exec(updateQuery, status, id, userId)
	if err != nil{
		tx.Rollback()
		return err
	}
	cnt, err := res.RowsAffected()
	if err != nil{
		tx.Rollback()
		return err
	}
	if cnt == 0{
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

func (r *DriverOrdersPostgres) ChatFetch(userId, rideId,orderId int) ([]models.ChatMessages, error){
	var lists []models.ChatMessages
	listQuery := fmt.Sprintf("SELECT user_type,driver_id,client_id,ride_id,order_id,message_type,content,created_at FROM %s WHERE driver_id=$1 AND ride_id=$2 AND order_id=$3 ORDER BY id DESC", chatMessagesTable)
	err := r.db.Select(&lists, listQuery, userId, rideId, orderId)
	return lists, err
}
