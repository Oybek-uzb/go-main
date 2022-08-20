package handler

import (
	"abir/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
