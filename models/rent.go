package models

type Car struct {
	Id          int     `json:"id" db:"id"`
	Conditioner bool    `json:"conditioner" db:"conditioner"`
	Price       int     `json:"price" db:"price"`
	PhoneNumber *string `json:"phone_number" db:"phone_number"`
	Description string  `json:"description" db:"description"`
	Photo       *string `json:"photo" db:"photo"`
	MarkId      *int    `json:"car_mark_id" db:"car_mark_id"`
	ModelId     *int    `json:"car_model_id" db:"car_model_id"`
	CategoryId  *int    `json:"category_car_id" db:"category_car_id"`
	ColorId     *int    `json:"color_id" db:"color_id"`
	DistrictId  *int    `json:"district_id" db:"district_id"`
	RegionId    *int    `json:"region_id" db:"region_id"`
	CompanyId   *int    `json:"company_id" db:"rent_car_company_id"`
	FCTypeId    *int    `json:"fc_type_id" db:"fc_type_id"`
	PerCarId    *int    `json:"per_car_id" db:"per_car_id"`
}

type CarCategory struct {
	Id           int     `json:"id"`
	Photo        *string `json:"photo" db:"photo"`
	CategoryName *string `json:"category_name" db:"name"`
}

type CarByCategoryId struct {
	Id          int     `json:"id" db:"id"`
	ModelName   *string `json:"model_name" db:"model_name"`
	Photo       *string `json:"photo" db:"photo"`
	CompanyName *string `json:"company_name" db:"company_name"`
	Price       int     `json:"price" db:"price"`
}
