package models

import "gopkg.in/guregu/null.v3"

type Pagination struct {
	CurrentPage int `json:"current_page" db:"current_page"`
	LastPage int `json:"last_page" db:"last_page"`
	Total int `json:"total" db:"total"`
	PerPage int `json:"per_page" db:"per_page"`
	Data any `json:"data"`
}

type CarMarka struct {
	Id int `json:"id" db:"id"`
	Name string `json:"name"`
	Image string `json:"image"`
	IsPopular bool `json:"-"`
}

type CarModel struct {
	Id int `json:"id" db:"id"`
	Name string `json:"name"`
	Image string `json:"-"`
	IsPopular bool `json:"-"`
	CarMarkaId int `json:"-" db:"carmarka_id"`
}

type Color struct {
	Id int `json:"id" db:"id"`
	Name string `json:"name"`
	HexCode null.String `json:"hex_code" db:"hex_code"`
}

type Region struct {
	Id int `json:"id" db:"id"`
	IsCity bool `json:"is_city" db:"is_city"`
	Name string `json:"name"`
	Polygon string `json:"-"`
}

type District struct {
	Id int `json:"id" db:"id"`
	Name string `json:"name"`
	RegionId int `json:"-"`
	Polygon string `json:"-"`
}

type Tariff struct {
	Id int `json:"id" db:"id"`
	Name string `json:"name"`
}

type DriverCancelOrderOptions struct {
	Id int `json:"id" db:"id"`
	Options string `json:"options"`
}
type ClientCancelOrderOptions struct {
	Id int `json:"id" db:"id"`
	Options string `json:"options"`
	Type string `json:"type"`
}