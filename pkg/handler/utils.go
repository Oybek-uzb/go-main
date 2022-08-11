package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) utilsColor(c *gin.Context) {
	langId, err := getLangId(c)
	if err != nil {
		return
	}
	lists, err := h.services.Utils.GetColors(langId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, lists)
}

func (h *Handler) utilsCarMarka(c *gin.Context) {
	lists, err := h.services.Utils.GetCarMarkas()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, lists)
}

func (h *Handler) utilsCarModel(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	lists, err := h.services.Utils.GetCarModels(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, lists)
}

func (h *Handler) test(c *gin.Context) {
	err := h.services.Utils.Test(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, "ok")
}

func (h *Handler) utilsRegion(c *gin.Context) {
	langId, err := getLangId(c)
	if err != nil {
		return
	}
	lists, err := h.services.Utils.GetRegions(langId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, lists)
}

func (h *Handler) utilsDistrict(c *gin.Context) {
	langId, err := getLangId(c)
	if err != nil {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	lists, err := h.services.Utils.GetDistricts(langId, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, lists)
}

func (h *Handler) utilsDriverCancelOrderOptions(c *gin.Context) {
	langId, err := getLangId(c)
	if err != nil {
		return
	}
	lists, err := h.services.Utils.DriverCancelOrderOptions(langId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, lists)
}

func (h *Handler) utilsClientCancelOrderOptions(c *gin.Context) {
	langId, err := getLangId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	optionType := ""
	optionTypeParam, ok := c.GetQuery("type")
	if ok {
		optionType = optionTypeParam
	}
	lists, err := h.services.Utils.ClientCancelOrderOptions(langId, optionType)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, lists)
}

func (h *Handler) utilsClientRateOptions(c *gin.Context) {
	langId, err := getLangId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	optionType := ""
	optionTypeParam, ok := c.GetQuery("type")
	if ok {
		optionType = optionTypeParam
	}
	rate := 0
	rateParam, ok := c.GetQuery("rate")
	if ok {
		rateParamInt, err := strconv.Atoi(rateParam)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		rate = rateParamInt
	}
	lists, err := h.services.Utils.ClientRateOptions(langId, rate, optionType)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, lists)
}
