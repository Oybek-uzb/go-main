package models

type Car struct {
	Id           int     `json:"id" db:"id"`
	Conditioner  bool    `json:"conditioner" db:"conditioner"`
	Price        int     `json:"price" db:"price"`
	PhoneNumber  *string `json:"phone_number" db:"phone_number"`
	Description  string  `json:"description" db:"description"`
	Photo        *string `json:"photo" db:"photo"`
	MarkName     *string `json:"mark_name" db:"mark_name"`
	ModelName    *string `json:"model_name" db:"model_name"`
	CategoryName *string `json:"category_name" db:"category_name"`
	ColorName    *string `json:"color_name" db:"color_name"`
	DistrictName *string `json:"district_name" db:"district_name"`
	RegionName   *string `json:"region_name" db:"region_name"`
	CompanyName  *string `json:"company_name" db:"company_name"`
	FCTypeName   *string `json:"fc_type_name" db:"fc_type_name"`
	PerCarName   *string `json:"per_type_name" db:"per_type_name"`
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