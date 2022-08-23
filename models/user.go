package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type User struct {
	Id       int    `json:"-" db:"id"`
	DriverId *int   `json:"driver_id" form:"driver_id" db:"driver_id"`
	ClientId *int   `json:"client_id" form:"client_id" db:"client_id"`
	UserType string `json:"user_type" form:"user_type" db:"user_type"`
	Login    string `json:"login" form:"login" db:"login"`
	Password string `json:"password" form:"password"`
}

type Client struct {
	Id        int     `json:"id" db:"id"`
	Name      *string `json:"name" form:"name"`
	Surname   *string `json:"surname" form:"surname"`
	Birthdate *string `json:"birthdate" form:"birthdate"`
	Gender    *string `json:"gender" form:"gender"`
	Avatar    *string `json:"avatar" form:"avatar"`
	Phone     *string `json:"phone" db:"phone"`
}

func (a Client) ValidateCreate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Name, validation.Required, validation.Length(2, 50)),
		validation.Field(&a.Gender, validation.Required, validation.In("male", "female")),
	)
}
func (a Client) ValidateUpdate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Name, validation.NilOrNotEmpty, validation.Length(2, 50)),
		validation.Field(&a.Surname, validation.NilOrNotEmpty, validation.Length(2, 50)),
		validation.Field(&a.Birthdate, validation.NilOrNotEmpty, validation.Date("2006-01-02")),
		validation.Field(&a.Gender, validation.NilOrNotEmpty, validation.In("male", "female")),
	)
}

type Driver struct {
	Id                      int      `json:"id" db:"id"`
	Name                    *string  `json:"name" form:"name"`
	Surname                 *string  `json:"surname" form:"surname"`
	DateOfBirth             *string  `json:"date_of_birth" form:"date_of_birth" db:"date_of_birth"`
	Gender                  *string  `json:"gender" form:"gender"`
	Phone                   *string  `json:"phone" form:"phone"`
	Activity                *int     `json:"activity" form:"activity"`
	Rating                  *float32 `json:"rating" form:"rating"`
	Status                  *string  `json:"status" form:"status"`
	Photo                   *string  `json:"photo" form:"photo" bson:"photo"`
	Balance                 int      `json:"balance"`
	DocumentType            *string  `json:"document_type" form:"document_type" db:"document_type"`
	PassportSerial          *string  `json:"passport_serial" form:"passport_serial" db:"passport_serial"`
	PassportCopy1           *string  `json:"passport_copy1" form:"passport_copy1" db:"passport_copy1"`
	PassportCopy2           *string  `json:"passport_copy2" form:"passport_copy2" db:"passport_copy2"`
	PassportCopy3           *string  `json:"passport_copy3" form:"passport_copy3" db:"passport_copy3"`
	DriverLicense           *string  `json:"driver_license" form:"driver_license" db:"driver_license"`
	DriverLicenseExpiration *string  `json:"driver_license_expiration" form:"driver_license_expiration" db:"driver_license_expiration"`
	DriverLicensePhoto1     *string  `json:"driver_license_photo1" form:"driver_license_photo1" db:"driver_license_photo1"`
	DriverLicensePhoto2     *string  `json:"driver_license_photo2" form:"driver_license_photo2" db:"driver_license_photo2"`
	DriverLicensePhoto3     *string  `json:"driver_license_photo3" form:"driver_license_photo3" db:"driver_license_photo3"`
}

type DriverVerification struct {
	Id          int    `json:"-"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type DriverCar struct {
	Id               int            `json:"-" db:"id"`
	CarInfo          *DriverCarInfo `json:"car_info" form:"car_info"`
	PhotoTexpasport1 *string        `json:"photo_texpasport1" form:"photo_texpasport1" db:"photo_texpasport1"`
	PhotoTexpasport2 *string        `json:"photo_texpasport2" form:"photo_texpasport2" db:"photo_texpasport2"`
	CarNumber        *string        `json:"car_number" form:"car_number" db:"car_number"`
	CarYear          *int           `json:"car_year" form:"car_year" db:"car_year"`
	CarFront         *string        `json:"car_front" form:"car_front" db:"car_front"`
	CarLeft          *string        `json:"car_left" form:"car_left" db:"car_left"`
	CarBack          *string        `json:"car_back" form:"car_back" db:"car_back"`
	CarRight         *string        `json:"car_right" form:"car_right" db:"car_right"`
	CarFrontRow      *string        `json:"car_front_row" form:"car_front_row" db:"car_front_row"`
	CarFrontBack     *string        `json:"car_front_back" form:"car_front_back" db:"car_front_back"`
	CarBaggage       *string        `json:"car_baggage" form:"car_baggage" db:"car_baggage"`
	CarColorId       *int           `json:"car_color_id" form:"car_color_id" db:"car_color_id"`
	CarMarkaId       *int           `json:"car_marka_id" form:"car_marka_id" db:"car_marka_id"`
	CarModelId       *int           `json:"car_model_id" form:"car_model_id" db:"car_model_id"`
	DriverId         *int           `json:"-" db:"driver_id"`
}

type DriverCarInfo struct {
	CarColorName string `json:"car_color_name" db:"car_color_name"`
	CarMarkaName string `json:"car_marka_name" db:"car_marka_name"`
	CarModelName string `json:"car_model_name" db:"car_model_name"`
}

func (a Driver) ValidateCreate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Name, validation.Required, validation.Length(2, 50)),
		validation.Field(&a.Surname, validation.Required, validation.Length(2, 50)),
		validation.Field(&a.DateOfBirth, validation.Required, validation.Date("2006-01-02")),
		validation.Field(&a.Gender, validation.Required, validation.In("male", "female")),
	)
}

func (a Driver) ValidateUpdate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Name, validation.NilOrNotEmpty, validation.Length(2, 50)),
		validation.Field(&a.Surname, validation.NilOrNotEmpty, validation.Length(2, 50)),
		validation.Field(&a.DateOfBirth, validation.NilOrNotEmpty, validation.Date("2006-01-02")),
		validation.Field(&a.Gender, validation.NilOrNotEmpty, validation.In("male", "female")),
		validation.Field(&a.DocumentType, validation.NilOrNotEmpty, validation.In("passport", "id_card")),
		validation.Field(&a.DriverLicenseExpiration, validation.NilOrNotEmpty, validation.Date("2006-01-02")),
		validation.Field(&a.PassportSerial, validation.NilOrNotEmpty, validation.Length(7, 50)),
	)
}

type SavedAddresses struct {
	Id        int     `json:"id" db:"id"`
	UserId    int     `json:"-" form:"user_id" db:"user_id"`
	Name      *string `json:"name" form:"name" db:"name"`
	PlaceType *string `json:"place_type" form:"place_type" db:"place_type"`
	Location  *string `json:"location" form:"location"`
	Address   *string `json:"address" form:"address"`
}

func (a SavedAddresses) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Name, validation.Required, validation.Length(2, 100)),
		validation.Field(&a.PlaceType, validation.Required, validation.In("work", "home", "custom")),
		validation.Field(&a.Location, validation.Required, validation.Length(2, 100)),
		validation.Field(&a.Address, validation.Required, validation.Length(2, 200)),
	)
}

type CreditCards struct {
	Id             int     `json:"id" db:"id"`
	UserId         int     `json:"-" form:"user_id" db:"user_id"`
	CardInfo       *string `json:"-" form:"card_info" db:"card_info"`
	CardNumber     *string `json:"card_number" form:"card_number"`
	CardExpiration *string `json:"card_expiration" form:"card_expiration"`
	IsActive       *bool   `json:"is_active" form:"is_active" db:"is_active"`
	IsMain         *bool   `json:"is_main" form:"is_main" db:"is_main"`
}

func (a CreditCards) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.CardNumber, validation.Required, validation.Length(16, 16)),
		validation.Field(&a.CardExpiration, validation.Required, validation.Length(4, 4)),
	)
}
