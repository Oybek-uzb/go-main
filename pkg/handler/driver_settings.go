package handler

import (
	"abir/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) driverCityTariffs(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	langId, err := getLangId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	lists, err := h.services.DriverSettings.GetTariffs(userId, langId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, lists)
}

func (h *Handler) driverCityTariffsEnable(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	tariffIdParam := c.PostForm("tariff_id")
	tariffIdFromQuery, err := strconv.Atoi(tariffIdParam)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "tariff_id field is required")
		return
	}
	isActiveParam := c.PostForm("is_active")
	isActiveParamFromQuery, err := strconv.ParseBool(isActiveParam)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "is_active field is required")
		return
	}
	err = h.services.DriverSettings.TariffsEnable(userId, tariffIdFromQuery, isActiveParamFromQuery)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, "ok")
}

func (h *Handler) driverCitySetOnline(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	isOnlineParam := c.PostForm("is_online")
	isOnlineParamFromQuery, err := strconv.Atoi(isOnlineParam)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "is_online field is required")
		return
	}
	if isOnlineParamFromQuery != 1 && isOnlineParamFromQuery != 0 {
		newErrorResponse(c, http.StatusInternalServerError, "is_online field must be 1 or 0")
		return
	}
	err = h.services.DriverSettings.SetOnline(userId, isOnlineParamFromQuery)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, "ok")
}

func (h *Handler) driverStats(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	driverId, err := h.services.Authorization.GetDriverId(userId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	typParam := c.PostForm("type")
	if typParam == "" {
		newErrorResponse(c, http.StatusInternalServerError, "type field is required")
		return
	}
	if !(typParam == "days" || typParam == "months" || typParam == "weeks" || typParam == "years") {
		newErrorResponse(c, http.StatusInternalServerError, "wrong type")
		return
	}
	startParam := c.PostForm("start")
	startFromQuery, err := strconv.Atoi(startParam)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "start field is required")
		return
	}
	resp, err := utils.GetStats(driverId, typParam, startFromQuery)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, resp)
}
func (h *Handler) driverStatOrders(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	langId, err := getLangId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	driverId, err := h.services.Authorization.GetDriverId(userId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := utils.GetStatOrders(driverId, langId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, resp)
}
