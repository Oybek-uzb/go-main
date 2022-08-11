package models

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"strings"
)

type CityTariffs struct {
	Id              int     `json:"id"`
	Price           *int    `json:"price"`
	TariffName      *string `json:"tariff_name" db:"tariff_name"`
	StartPrice      *int    `json:"start_price" db:"start_price"`
	PricePerKm      *int    `json:"price_per_km" db:"price_per_km"`
	PricePerKmOuter *int    `json:"price_per_km_outer" db:"price_per_km_outer"`
	ACPrice         *int    `json:"ac_price" db:"ac_price"`
	Expectation     *int    `json:"expectation" db:"expectation"`
	Cars            *string `json:"cars"`
	Description     *string `json:"description"`
	IsRushHour      bool    `json:"is_rush_hour" db:"is_rush_hour"`
	Icon            *string `json:"icon"`
	Image           *string `json:"image"`
}

type TariffInfo struct {
	StartPrice      *int `json:"start_price" db:"start_price"`
	PricePerKm      *int `json:"price_per_km" db:"price_per_km"`
	PricePerKmOuter *int `json:"price_per_km_outer" db:"price_per_km_outer"`
	ACPrice         *int `json:"ac_price" db:"ac_price"`
	Expectation     *int `json:"expectation" db:"expectation"`
}

type Ride struct {
	Id             int                   `json:"id" db:"id"`
	DriverId       int                   `json:"driver_id,omitempty" db:"driver_id"`
	FromDistrictId string                `json:"from_district_id" form:"from_district_id" db:"from_district_id"`
	ToDistrictId   string                `json:"to_district_id" form:"to_district_id" db:"to_district_id"`
	From           *string               `json:"from,omitempty"`
	To             *string               `json:"to,omitempty"`
	DepartureDate  string                `json:"departure_date" form:"departure_date" db:"departure_date"`
	Price          string                `json:"price" form:"price" db:"price"`
	PassengerCount string                `json:"passenger_count" form:"passenger_count" db:"passenger_count"`
	Comments       *string               `json:"comments" form:"comments" db:"comments"`
	ViewCount      int                   `json:"view_count" db:"view_count"`
	Status         string                `json:"status" db:"status"`
	Notifications  *[]RideNotification   `json:"notifications,omitempty"`
	OrderList      *[]InterregionalOrder `json:"order_list,omitempty"`
	CreatedAt      string                `json:"created_at" db:"created_at"`
}
type RideNotification struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	OrderId   int    `json:"order_id" db:"order_id"`
	CreatedAt string `json:"created_at" db:"created_at"`
}

type ClientRideList struct {
	RideId         int                     `json:"ride_id" db:"ride_id"`
	Driver         *map[string]interface{} `json:"driver" db:"driver"`
	DriverCar      *map[string]interface{} `json:"driver_car" db:"driver_car"`
	DriverCarInfo  *DriverCarInfo          `json:"driver_car_info" db:"driver_car_info"`
	DriverId       int                     `json:"driver_id" db:"driver_id"`
	FromDistrictId string                  `json:"from_district_id" db:"from_district_id"`
	ToDistrictId   string                  `json:"to_district_id" db:"to_district_id"`
	FromDistrict   *string                 `json:"from_district" db:"from_district"`
	ToDistrict     *string                 `json:"to_district" db:"to_district"`
	DepartureTime  string                  `json:"departure_time" db:"departure_time"`
	Price          string                  `json:"price" db:"price"`
	PassengerCount string                  `json:"passenger_count" db:"passenger_count"`
	Comments       *string                 `json:"comments" db:"comments"`
	Status         string                  `json:"status" db:"status"`
}

type Activity struct {
	OrderType  string  `json:"order_type" db:"order_type"`
	OrderId    int     `json:"order_id" db:"order_id"`
	SubOrderId int     `json:"-" db:"sub_order_id"`
	RideId     int     `json:"ride_id" db:"ride_id"`
	Direction  *string `json:"direction"`
	OrderTime  string  `json:"order_time" db:"order_time"`
	From       string  `json:"-"`
	To         *string `json:"-"`
	TariffId   *string `json:"tariff_id,omitempty" db:"tariff_id"`
	Status     *string `json:"status"`
}

type Order struct {
	Id          int    `json:"id" db:"id"`
	DriverId    *int   `json:"driver_id" db:"driver_id"`
	ClientId    int    `json:"client_id" db:"client_id"`
	OrderId     int    `json:"order_id" db:"order_id"`
	OrderType   string `json:"order_type" db:"order_type"`
	OrderStatus string `json:"order_status" db:"order_status"`
	CreatedAt   string `json:"created_at" db:"created_at"`
}
type CancelOrRateReasons struct {
	ReasonId string  `json:"reason_id" form:"reason_id"`
	Rate     int     `json:"rate" form:"rate"`
	Comments *string `json:"comments" form:"comments"`
}
type CanceledOrders struct {
	Id        int     `json:"id" db:"id"`
	OrderType string  `json:"order_type" db:"order_type"`
	UserType  string  `json:"user_type" db:"user_type"`
	UserId    int     `json:"user_id" db:"user_id"`
	OrderId   int     `json:"order_id" db:"order_id"`
	Comments  *string `json:"comments" form:"comments"`
	CreatedAt string  `json:"created_at" db:"created_at"`
}
type CanceledOrderReasons struct {
	Id              int  `json:"id" db:"id"`
	CanceledOrderId int  `json:"canceled_order_id" db:"canceled_order_id"`
	ReasonId        *int `json:"reason_id" db:"reason_id"`
}

type RatedOrders struct {
	Id        int     `json:"id" db:"id"`
	OrderType string  `json:"order_type" db:"order_type"`
	UserType  string  `json:"user_type" db:"user_type"`
	Rate      int     `json:"rate"`
	UserId    int     `json:"user_id" db:"user_id"`
	OrderId   int     `json:"order_id" db:"order_id"`
	Comments  *string `json:"comments" form:"comments"`
	CreatedAt string  `json:"created_at" db:"created_at"`
}
type RatedOrderReasons struct {
	Id           int  `json:"id" db:"id"`
	RatedOrderId int  `json:"rated_order_id" db:"rated_order_id"`
	ReasonId     *int `json:"reason_id" db:"reason_id"`
}

type InterregionalOrder struct {
	Id             int     `json:"id" db:"id"`
	Client         *Client `json:"client,omitempty"`
	ClientId       int     `json:"client_id,omitempty" db:"client_id"`
	RideId         int     `json:"ride_id,omitempty" db:"ride_id"`
	FromDistrictId string  `json:"from_district_id,omitempty" db:"from_district_id"`
	ToDistrictId   string  `json:"to_district_id,omitempty" db:"to_district_id"`
	Price          float32 `json:"price,omitempty" db:"price"`
	PassengerCount int     `json:"passenger_count" db:"passenger_count"`
	DepartureDate  string  `json:"departure_date,omitempty" db:"departure_date"`
	OrderStatus    *string `json:"order_status,omitempty" db:"order_status"`
	Comments       *string `json:"comments" db:"comments"`
	CreatedAt      string  `json:"created_at" db:"created_at"`
}

type PointsArr struct {
	Address  string `json:"address"`
	Location string `json:"location"`
}
type CityOrderPoints struct {
	Distance int         `json:"distance"`
	Points   []PointsArr `json:"points"`
}
type CityOrderRequest struct {
	DriverLastLocation *string  `json:"driver_last_location" form:"driver_last_location"`
	OrderAmount        *float64 `json:"order_amount" form:"order_amount"`
	WaitTime           *int     `json:"wait_time" form:"wait_time"`
	WaitTimeAmount     *float64 `json:"wait_time_amount" form:"wait_time_amount"`
	RideDistance       *int     `json:"ride_distance" form:"ride_distance"`
	RideAmount         *float64 `json:"ride_amount" form:"ride_amount"`
	RideTime           *int     `json:"ride_time" form:"ride_time"`
	Commission         *float64 `json:"commission" form:"commission"`
}
type CityOrder struct {
	Id               int                     `json:"id" db:"id"`
	TariffInfo       *TariffInfo             `json:"tariff_info"`
	DriverId         *int                    `json:"driver_id" db:"driver_id"`
	ClientId         int                     `json:"client_id" db:"client_id"`
	Client           *Client                 `json:"client,omitempty"`
	Driver           *map[string]interface{} `json:"driver" db:"driver"`
	DriverCar        *map[string]interface{} `json:"driver_car" db:"driver_car"`
	DriverCarInfo    *DriverCarInfo          `json:"driver_car_info" db:"driver_car_info"`
	Points           string                  `json:"-" db:"points" form:"points"`
	PointsArr        *CityOrderPoints        `json:"points_arr,omitempty"`
	From             *string                 `json:"from,omitempty"`
	To               *string                 `json:"to,omitempty"`
	TariffId         string                  `json:"tariff_id" db:"tariff_id" form:"tariff_id"`
	CargoType        string                  `json:"cargo_type" db:"cargo_type" form:"cargo_type"`
	PaymentType      string                  `json:"payment_type" db:"payment_type" form:"payment_type"`
	CardId           *string                 `json:"card_id,omitempty" db:"card_id" form:"card_id"`
	HasConditioner   string                  `json:"has_conditioner" db:"has_conditioner" form:"has_conditioner"`
	ForAnother       string                  `json:"for_another" db:"for_another" form:"for_another"`
	ForAnotherPhone  *string                 `json:"for_another_phone" db:"for_another_phone" form:"for_another_phone"`
	ReceiverComments *string                 `json:"receiver_comments" db:"receiver_comments" form:"receiver_comments"`
	ReceiverPhone    *string                 `json:"receiver_phone" db:"receiver_phone" form:"receiver_phone"`
	Price            string                  `json:"price" db:"price" form:"price"`
	RideInfo         *string                 `json:"-" db:"ride_info"`
	RideInfoArr      *CityOrderRequest       `json:"ride_info_arr,omitempty"`
	OrderStatus      string                  `json:"order_status,omitempty" db:"order_status" form:"order_status"`
	Comments         *string                 `json:"comments" db:"comments" form:"comments"`
	CreatedAt        *string                 `json:"created_at" db:"created_at"`
}

type ChatMessages struct {
	Id          int    `json:"-" db:"id"`
	UserType    string `json:"from" db:"user_type"`
	DriverId    int    `json:"driver_id" db:"driver_id"`
	ClientId    int    `json:"client_id" db:"client_id"`
	RideId      int    `json:"ride_id" db:"ride_id"`
	OrderId     int    `json:"order_id" db:"order_id"`
	MessageType string `json:"type" db:"message_type"`
	Content     string `json:"content"`
	CreatedAt   string `json:"created_at" db:"created_at"`
}

func (a Ride) ValidateCreate() error {
	a.DepartureDate = fmt.Sprintf("%s:00+0500", strings.Replace(a.DepartureDate, " ", "T", -1))
	return validation.ValidateStruct(&a,
		validation.Field(&a.FromDistrictId, validation.Required, is.Digit),
		validation.Field(&a.ToDistrictId, validation.Required, is.Digit),
		validation.Field(&a.DepartureDate, validation.Required, validation.Date("2006-01-02T15:04:05-0700")),
		validation.Field(&a.Price, validation.Required, is.Digit),
		validation.Field(&a.PassengerCount, validation.Required, is.Int),
		validation.Field(&a.Comments, validation.NilOrNotEmpty, validation.Length(3, 200)),
	)
}

func (a Ride) ValidateUpdate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Price, validation.Required, is.Digit),
		validation.Field(&a.PassengerCount, validation.Required, is.Digit),
	)
}

func (a Ride) ValidateSearch() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.FromDistrictId, validation.Required, is.Digit),
		validation.Field(&a.ToDistrictId, validation.Required, is.Digit),
		validation.Field(&a.DepartureDate, validation.Required, validation.Date("2006-01-02")),
	)
}

func (a Ride) ValidateBook() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.PassengerCount, validation.Required, is.Int),
		validation.Field(&a.Comments, validation.NilOrNotEmpty, validation.Length(3, 200)),
	)
}

func (a CityOrder) ValidateOrder() error {
	var fieldRules []*validation.FieldRules
	fieldRules = append(fieldRules, validation.Field(&a.Points, validation.Required, is.JSON),
		validation.Field(&a.TariffId, validation.Required, is.Int),
		validation.Field(&a.CargoType, validation.Required, validation.In("no", "small", "medium", "large")),
		validation.Field(&a.PaymentType, validation.Required, validation.In("cash", "card")),
		validation.Field(&a.HasConditioner, validation.Required, validation.In("0", "1"), is.Digit),
		validation.Field(&a.ForAnother, validation.Required, validation.In("0", "1"), is.Digit),
		validation.Field(&a.Price, validation.Required, is.Digit),
	)
	if a.PaymentType == "card" {
		fieldRules = append(fieldRules, validation.Field(&a.CardId, validation.Required, is.Int))
	}
	if a.ForAnother == "1" {
		fieldRules = append(fieldRules, validation.Field(&a.ForAnotherPhone, validation.Required, validation.Length(3, 100)))
	}
	if a.TariffId == "8" {
		fieldRules = append(fieldRules, validation.Field(&a.ReceiverPhone, validation.Required, validation.Length(3, 100)))
	}

	return validation.ValidateStruct(&a, fieldRules...)
}
