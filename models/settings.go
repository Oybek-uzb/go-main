package models

type DriverTariffs struct {
	Id int `json:"id"`
	Name string `json:"name"`
	IsActive bool `json:"is_active" db:"is_active"`
}
