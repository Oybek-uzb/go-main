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

func (h *Handler) rentCarsByCategoryID(c *gin.Context) {
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
