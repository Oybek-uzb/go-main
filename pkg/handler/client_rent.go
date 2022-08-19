package handler

import (
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
	lists, err := h.services.RentCars.GetCategoriesList()
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

func (h *Handler) rentMyCompaniesList(c *gin.Context) {
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
