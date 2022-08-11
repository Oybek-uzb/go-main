package models

import "gopkg.in/guregu/null.v3"

type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
type PointsRequest struct {
	Points []LatLng `json:"points"`
}
type NewPointsRequest struct {
	TariffId *int         `json:"tariff_id"`
	Points   [][2]float64 `json:"points"`
}

type DriverOrderSocket struct {
	Id       int    `json:"id"`
	DriverId int    `json:"driver_id"`
	Status   string `json:"status"`
}

type ClientOrderSocket struct {
	Id          int      `json:"id"`
	ClientId    int      `json:"client_id"`
	Location    *string  `json:"location"`
	OrderAmount *float64 `json:"order_amount"`
	Status      string   `json:"status"`
}

type Pagination struct {
	CurrentPage int `json:"current_page" db:"current_page"`
	LastPage    int `json:"last_page" db:"last_page"`
	Total       int `json:"total" db:"total"`
	PerPage     int `json:"per_page" db:"per_page"`
	Data        any `json:"data"`
}

type CarMarka struct {
	Id        int    `json:"id" db:"id"`
	Name      string `json:"name"`
	Image     string `json:"image"`
	IsPopular bool   `json:"-"`
}

type CarModel struct {
	Id         int    `json:"id" db:"id"`
	Name       string `json:"name"`
	Image      string `json:"-"`
	IsPopular  bool   `json:"-"`
	CarMarkaId int    `json:"-" db:"carmarka_id"`
}

type Color struct {
	Id      int         `json:"id" db:"id"`
	Name    string      `json:"name"`
	HexCode null.String `json:"hex_code" db:"hex_code"`
}

type Region struct {
	Id      int    `json:"id" db:"id"`
	IsCity  bool   `json:"is_city" db:"is_city"`
	Name    string `json:"name"`
	Polygon string `json:"-"`
}

type District struct {
	Id       int    `json:"id" db:"id"`
	Name     string `json:"name"`
	RegionId int    `json:"-"`
	Polygon  string `json:"-"`
}

type Tariff struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name"`
}

type DriverCancelOrderOptions struct {
	Id      int    `json:"id" db:"id"`
	Options string `json:"options"`
}
type ClientCancelOrderOptions struct {
	Id      int    `json:"id" db:"id"`
	Options string `json:"options"`
	Type    string `json:"type"`
}
type ClientRateOptions struct {
	Id      int    `json:"id" db:"id"`
	Rate    int    `json:"rate" db:"rating"`
	Options string `json:"options"`
	Type    string `json:"type"`
}
