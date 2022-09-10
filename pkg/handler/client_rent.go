package handler

import (
	"abir/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) rentCategoriesList(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	langId, err := getLangId(c)
	if err != nil {
		return
	}

	lists, err := h.services.RentCars.GetCategoriesList(langId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, lists)
}

func (h *Handler) rentForEventsCategories(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	langId, err := getLangId(c)
	if err != nil {
		return
	}

	list, err := h.services.RentCars.GetCategoriesForEvents(langId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, list)
}

func (h *Handler) rentCarsByCategoryId(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	categoryId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid category_id param")
		return
	}

	cars, err := h.services.RentCars.GetCarsByCategoryId(categoryId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, cars)
}

func (h *Handler) rentCarByCategoryIdCarId(c *gin.Context) {
	langId, err := getLangId(c)
	if err != nil {
		return
	}

	_, err = getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	categoryId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid category_id param")
		return
	}

	carId, err := strconv.Atoi(c.Param("car_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid car_id param")
		return
	}

	car, err := h.services.RentCars.GetCarByCategoryIdCarId(categoryId, carId, langId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, car)
}

func (h *Handler) rentCarFromCategoryCreate(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var rentDetails models.RentCarDetails
	if err = c.Bind(&rentDetails); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = rentDetails.ValidateRentCarDetails()
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	carId, err := strconv.Atoi(c.Param("car_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid car_id param")
		return
	}

	car, err := h.services.RentCars.PostRentCarByCarId(userId, carId, rentDetails)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, car)
}

func (h *Handler) rentCarFromCompanyCreate(c *gin.Context) {
	h.rentCarFromCategoryCreate(c)
}

func (h *Handler) rentCompaniesList(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	companies, err := h.services.RentCars.GetCompaniesList()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, companies)
}

func (h *Handler) rentCompanyById(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	companyId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid company_id param")
		return
	}

	carCompany, err := h.services.RentCars.GetCarsByCompanyId(companyId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, carCompany)
}

func (h *Handler) rentCarByCompanyIdCarId(c *gin.Context) {
	langId, err := getLangId(c)
	if err != nil {
		return
	}

	_, err = getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	companyId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid company_id param")
		return
	}

	carId, err := strconv.Atoi(c.Param("car_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid car_id param")
		return
	}

	car, err := h.services.RentCars.GetCarByCompanyIdCarId(companyId, carId, langId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, car)
}

func (h *Handler) myCompaniesList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	myCompanies, err := h.services.RentCars.GetMyCompaniesList(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, myCompanies)
}

func (h *Handler) rentForEventsMyCars(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	langId, err := getLangId(c)
	if err != nil {
		return
	}

	myCars, err := h.services.RentCars.GetMyCarsForEvents(userId, langId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, myCars)
}

func (h *Handler) rentForEventsMyCarCreate(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var car models.CarCreate
	if err := c.Bind(&car); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = car.ValidateCarCreate()
	if err != nil {
		fmt.Println("validation")
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}
	carId, err := h.services.RentCars.PostMyCarForEvents(c.Request.Context(), userId, car)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessResponse(c, http.StatusOK, carId)
}

func (h *Handler) myCompanyById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	myCompanyId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid company_id param")
		return
	}

	carCompany, err := h.services.RentCars.GetMyCompanyById(userId, myCompanyId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, carCompany)
}

func (h *Handler) rentMyCompaniesCreate(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var company models.RentMyCompanyCreate
	if err := c.Bind(&company); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = company.ValidateCompanyCreate()
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}
	companyId, err := h.services.RentCars.PostMyCompany(c.Request.Context(), userId, company)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessResponse(c, http.StatusOK, companyId)
}

func (h *Handler) rentAnnouncementCreate(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	myCompanyId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid company_id param")
		return
	}

	var car models.CarCreate
	if err := c.Bind(&car); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = car.ValidateCarCreate()
	if err != nil {
		fmt.Println("validation")
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}
	carId, err := h.services.RentCars.PostMyCar(c.Request.Context(), userId, myCompanyId, car)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessResponse(c, http.StatusOK, carId)
}

func (h *Handler) rentAnnouncementUpdate(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	myCompanyId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid company_id param")
		return
	}

	carId, err := strconv.Atoi(c.Param("car_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid car_id param")
		return
	}

	var car models.CarCreate
	if err := c.Bind(&car); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = car.ValidateCarCreate()
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}
	carId, err = h.services.RentCars.PutMyCar(c.Request.Context(), userId, carId, myCompanyId, car)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessResponse(c, http.StatusOK, carId)
}

func (h *Handler) rentAnnouncementDelete(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	myCompanyId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid company_id param")
		return
	}

	carId, err := strconv.Atoi(c.Param("car_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid car_id param")
		return
	}

	carId, err = h.services.RentCars.DeleteMyCar(userId, carId, myCompanyId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessResponse(c, http.StatusOK, carId)
}

func (h *Handler) myCarPark(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	inDiscount := false
	indicator, ok := c.GetQuery("in_discount")
	if ok {
		parsedIndicator, err := strconv.ParseBool(indicator)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		inDiscount = parsedIndicator
	}

	companyId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid company_id param")
		return
	}

	cars, err := h.services.RentCars.GetMyCarParkByCompanyId(userId, companyId, inDiscount)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, cars)
}

func (h *Handler) myCarByCompanyId(c *gin.Context) {
	h.rentCarByCompanyIdCarId(c)
}
