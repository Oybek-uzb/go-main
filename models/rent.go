package models

import validation "github.com/go-ozzo/ozzo-validation"

type Car struct {
	Id               int     `json:"id" db:"id"`
	Conditioner      bool    `json:"conditioner" db:"conditioner"`
	Price            int     `json:"price" db:"price"`
	PhoneNumber      *string `json:"phone_number" db:"phone_number"`
	Description      string  `json:"description" db:"description"`
	Photo            *string `json:"photo" db:"photo"`
	MarkName         *string `json:"mark_name" db:"mark_name"`
	ModelName        *string `json:"model_name" db:"model_name"`
	CategoryName     *string `json:"category_name" db:"category_name"`
	ColorName        *string `json:"color_name" db:"color_name"`
	DistrictName     *string `json:"district_name" db:"district_name"`
	RegionName       *string `json:"region_name" db:"region_name"`
	CompanyName      *string `json:"company_name" db:"company_name"`
	FCTypeName       *string `json:"fc_type_name" db:"fc_type_name"`
	PerCarName       *string `json:"per_type_name" db:"per_type_name"`
	InDiscount       bool    `json:"in_discount" db:"in_discount"`
	Discount         *int    `json:"discount" db:"discount"`
	TransmissionName *string `json:"transmission_name" db:"transmission_name"`
}

type CarCategory struct {
	Id    int     `json:"id"`
	Photo *string `json:"photo" db:"photo"`
	Name  *string `json:"category_name" db:"name"`
}

type CarByCategoryId struct {
	Id          int     `json:"id" db:"id"`
	ModelName   *string `json:"model_name" db:"model_name"`
	Photo       *string `json:"photo" db:"photo"`
	CompanyName *string `json:"company_name" db:"company_name"`
	Price       int     `json:"price" db:"price"`
	InDiscount  bool    `json:"in_discount" db:"in_discount"`
	Discount    *int    `json:"discount" db:"discount"`
}

type CarCompany struct {
	Id          int     `json:"id"`
	Photo       *string `json:"photo" db:"photo"`
	Name        *string `json:"company_name" db:"name"`
	PhoneNumber *string `json:"phone_number" db:"phone_number"`
	Description *string `json:"description" db:"description"`
}

type CarCompanyDetails struct {
	Id          int              `json:"id"`
	Photo       *string          `json:"photo" db:"photo"`
	Name        *string          `json:"company_name" db:"name"`
	Description *string          `json:"description" db:"description"`
	WebSite     *string          `json:"web_site" db:"web_site"`
	PhoneNumber *string          `json:"phone_number" db:"phone_number"`
	Cars        []CarByCompanyId `json:"company_cars"`
}

type CarByCompanyId struct {
	Id          int     `json:"id" db:"id"`
	ModelName   *string `json:"model_name" db:"model_name"`
	Photo       *string `json:"photo" db:"photo"`
	CompanyName *string `json:"company_name" db:"company_name"`
	Price       int     `json:"price" db:"price"`
	InDiscount  bool    `json:"in_discount" db:"in_discount"`
	Discount    *int    `json:"discount" db:"discount"`
}

type RentMyCompanyCreate struct {
	Id          int     `json:"id" db:"id"`
	Photo       *string `json:"photo" db:"photo" form:"photo"`
	CompanyName *string `json:"company_name" db:"name" form:"company_name"`
	Description *string `json:"description" db:"description" form:"description"`
	WebSite     *string `json:"web_site" db:"web_site" form:"web_site"`
	PhoneNumber *string `json:"phone_number" db:"phone_number" form:"phone_number"`
	Status      *string `json:"status" db:"status" form:"status"`
	OwnerId     int     `json:"owner_id" db:"owner_id"`
}

type RentCarDetails struct {
	FromDate    *string `json:"from_date" db:"start_time"`
	ToDate      *string `json:"to_date" db:"end_time"`
	Description *string `json:"description" db:"description"`
}

func (a RentCarDetails) ValidateRentCarDetails() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.FromDate, validation.Required, validation.Length(2, 100)),
		validation.Field(&a.ToDate, validation.Required, validation.Length(2, 100)),
	)
}

func (a RentMyCompanyCreate) ValidateCompanyCreate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.CompanyName, validation.Required, validation.Length(2, 100)),
		validation.Field(&a.Status, validation.In("new", "moderating", "cancelled", "blocked", "checked")),
		validation.Field(&a.PhoneNumber, validation.Required, validation.Length(9, 13)),
		validation.Field(&a.Description, validation.Required),
	)
}
