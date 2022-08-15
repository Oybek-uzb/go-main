package handler

import (
	"abir/models"
	"abir/pkg/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) clientOrdersActivityActive(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	langId, err := getLangId(c)
	if err != nil {
		return
	}
	page := 1
	pageParam, ok := c.GetQuery("page")
	if ok {
		pageFromQuery, err := strconv.Atoi(pageParam)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		page = pageFromQuery
	}
	orderType := ""
	orderTypeParam, ok := c.GetQuery("order_type")
	if ok {
		orderType = orderTypeParam
	}
	districts, err := h.services.Utils.GetDistrictsArr(langId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	tariffs, err := h.services.Utils.GetTariffs(langId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	lists, _, err := h.services.ClientOrders.Activity(userId, page, "active", orderType)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	for idx, list := range lists {
		var direction string
		var directionWithTariff string
		if list.OrderType == "city" {
			if list.To == nil {
				direction = fmt.Sprintf("%s -> %s", list.From, utils.Translation["not_set"][langId])
			} else {
				direction = fmt.Sprintf("%s -> %s", list.From, *list.To)
			}
		} else {
			fromId, err := strconv.Atoi(list.From)
			if err != nil {
				continue
			}
			if list.To != nil {
				toId, err := strconv.Atoi(*list.To)
				if err != nil {
					continue
				}
				direction = fmt.Sprintf("%s -> %s", districts[fromId], districts[toId])
			}
		}
		if list.TariffId != nil {
			directionWithTariff = fmt.Sprintf("%s, %s", tariffs[*list.TariffId], direction)
			lists[idx].Direction = &directionWithTariff
		} else {
			lists[idx].Direction = &direction
		}
	}
	newSuccessResponse(c, http.StatusOK, lists)
}
func (h *Handler) clientOrdersActivityRecentlyCompleted(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	langId, err := getLangId(c)
	if err != nil {
		return
	}
	districts, err := h.services.Utils.GetDistrictsArr(langId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	orderType := ""
	orderTypeParam, ok := c.GetQuery("order_type")
	if ok {
		orderType = orderTypeParam
	}
	tariffs, err := h.services.Utils.GetTariffs(langId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	lists, _, err := h.services.ClientOrders.Activity(userId, 1, "recently-completed", orderType)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	for i, list := range lists {
		var direction string
		var directionWithTariff string
		if list.OrderType == "city" {
			if list.To == nil {
				direction = fmt.Sprintf("%s -> %s", list.From, utils.Translation["not_set"][langId])
			} else {
				direction = fmt.Sprintf("%s -> %s", list.From, *list.To)
			}

		} else {
			fromId, err := strconv.Atoi(list.From)
			if err != nil {
				continue
			}
			if list.To != nil {
				toId, err := strconv.Atoi(*list.To)
				if err != nil {
					continue
				}
				direction = fmt.Sprintf("%s -> %s", districts[fromId], districts[toId])
			}
		}
		if list.TariffId != nil {
			directionWithTariff = fmt.Sprintf("%s, %s", tariffs[*list.TariffId], direction)
			lists[i].Direction = &directionWithTariff
		} else {
			lists[i].Direction = &direction
		}
	}
	newSuccessResponse(c, http.StatusOK, lists)
}
func (h *Handler) clientOrdersActivityHistory(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	langId, err := getLangId(c)
	if err != nil {
		return
	}
	districts, err := h.services.Utils.GetDistrictsArr(langId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	tariffs, err := h.services.Utils.GetTariffs(langId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	page := 1
	pageParam, ok := c.GetQuery("page")
	if ok {
		pageFromQuery, err := strconv.Atoi(pageParam)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		page = pageFromQuery
	}
	orderType := ""
	orderTypeParam, ok := c.GetQuery("order_type")
	if ok {
		orderType = orderTypeParam
	}
	lists, pagination, err := h.services.ClientOrders.Activity(userId, page, "history", orderType)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	for i, list := range lists {
		var direction string
		var directionWithTariff string
		if list.OrderType == "city" {
			if list.To == nil {
				direction = fmt.Sprintf("%s -> %s", list.From, utils.Translation["not_set"][langId])
			} else {
				direction = fmt.Sprintf("%s -> %s", list.From, *list.To)
			}

		} else {
			fromId, err := strconv.Atoi(list.From)
			if err != nil {
				continue
			}
			if list.To != nil {
				toId, err := strconv.Atoi(*list.To)
				if err != nil {
					continue
				}
				direction = fmt.Sprintf("%s -> %s", districts[fromId], districts[toId])
			}
		}
		if list.TariffId != nil {
			directionWithTariff = fmt.Sprintf("%s, %s", tariffs[*list.TariffId], direction)
			lists[i].Direction = &directionWithTariff
		} else {
			lists[i].Direction = &direction
		}
	}
	pagination.Data = lists
	newSuccessResponse(c, http.StatusOK, pagination)
}
func (h *Handler) clientOrdersRideList(c *gin.Context) {
	langId, err := getLangId(c)
	if err != nil {
		return
	}
	var ride models.Ride

	if err := c.Bind(&ride); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = ride.ValidateSearch()
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}
	page := 1
	pageParam, ok := c.GetQuery("page")
	if ok {
		pageFromQuery, err := strconv.Atoi(pageParam)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		page = pageFromQuery
	}
	list, pagination, err := h.services.ClientOrders.RideList(ride, langId, page)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	for i, rideList := range list {
		driver, car, carInfo, err := h.services.Authorization.GetDriverInfo(langId, rideList.DriverId)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		list[i].DriverCarInfo = &carInfo
		list[i].DriverCar = &map[string]interface{}{
			"car_number": car.CarNumber,
		}
		list[i].Driver = &map[string]interface{}{
			"name":    driver.Name,
			"surname": driver.Surname,
			"photo":   driver.Photo,
			"phone":   driver.Phone,
			"rating":  driver.Rating,
		}
	}
	pagination.Data = list
	newSuccessResponse(c, http.StatusOK, pagination)
}
func (h *Handler) clientOrdersRideSingle(c *gin.Context) {
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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	list, err := h.services.ClientOrders.RideSingle(langId, id, userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	driver, car, carInfo, err := h.services.Authorization.GetDriverInfo(langId, list.DriverId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	list.DriverCarInfo = &carInfo
	list.DriverCar = &map[string]interface{}{
		"car_number": car.CarNumber,
	}
	list.Driver = &map[string]interface{}{
		"name":    driver.Name,
		"surname": driver.Surname,
		"photo":   driver.Photo,
		"phone":   driver.Phone,
		"rating":  driver.Rating,
	}
	newSuccessResponse(c, http.StatusOK, list)
}

func (h *Handler) clientOrdersRideSingleBook(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	rideId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	var bookRide models.Ride
	if err = c.Bind(&bookRide); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = bookRide.ValidateBook()
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}
	orderId, err := h.services.ClientOrders.RideSingleBook(bookRide, rideId, userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, map[string]int{
		"id": orderId,
	})
}
func (h *Handler) clientOrdersRideSingleCancel(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	rideId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	orderId, err := strconv.Atoi(c.Param("order_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid order_id param")
		return
	}
	var cancelRide models.CancelOrRateReasons
	if err = c.Bind(&cancelRide); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.services.ClientOrders.RideSingleCancel(cancelRide, rideId, orderId, userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, "ok")
}
func (h *Handler) clientOrdersRideSingleStatus(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	rideId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	status, err := h.services.ClientOrders.RideSingleStatus(rideId, userId)
	if err != nil {
		if err != sql.ErrNoRows {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, "Ride not booked")
		}
		return
	}
	newSuccessResponse(c, http.StatusOK, status)
}

func (h *Handler) clientChatFetch(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	rideId := 0
	rideIdParam, ok := c.GetQuery("ride_id")
	if ok {
		rideIdFromQuery, err := strconv.Atoi(rideIdParam)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		rideId = rideIdFromQuery
	}
	orderId := 0
	orderIdParam, ok := c.GetQuery("order_id")
	if ok {
		orderIdFromQuery, err := strconv.Atoi(orderIdParam)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		orderId = orderIdFromQuery
	}
	lists, err := h.services.ClientOrders.ChatFetch(userId, rideId, orderId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, lists)
}

func (h *Handler) clientCityTariffs(c *gin.Context) {
	langId, err := getLangId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var points models.NewPointsRequest
	if err := c.Bind(&points); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if len(points.Points) == 0 {
		newErrorResponse(c, http.StatusBadRequest, "points field is required")
		return
	}
	lists, err := h.services.ClientOrders.CityTariffs(points.Points, langId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, lists)
}

func (h *Handler) clientCityOrder(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var order models.CityOrder
	if err = c.Bind(&order); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = order.ValidateOrder()
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}
	orderId, err := h.services.ClientOrders.CityNewOrder(order, userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	err = utils.SearchTaxi(orderId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, map[string]int{
		"id": orderId,
	})
}
func (h *Handler) clientCityOrderView(c *gin.Context) {
	langId, err := getLangId(c)
	if err != nil {
		return
	}
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	orderId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid order_id param")
		return
	}
	order, err := h.services.ClientOrders.CityOrderView(orderId, userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if order.DriverId != nil {
		driver, car, carInfo, err := h.services.Authorization.GetDriverInfo(langId, *order.DriverId)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		order.DriverCarInfo = &carInfo
		order.DriverCar = &map[string]interface{}{
			"car_number": car.CarNumber,
		}
		order.Driver = &map[string]interface{}{
			"name":    driver.Name,
			"surname": driver.Surname,
			"photo":   driver.Photo,
			"phone":   driver.Phone,
			"rating":  driver.Rating,
		}
	}
	var pointsArr models.CityOrderPoints
	err = json.Unmarshal([]byte(order.Points), &pointsArr)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	order.PointsArr = &pointsArr
	tariffId, err := strconv.Atoi(order.TariffId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	tariffInfo, err := h.services.DriverOrders.CityTariffInfo(pointsArr.Points[0].Location, tariffId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	order.TariffInfo = &tariffInfo
	if order.RideInfo != nil {
		var reqArr models.CityOrderRequest
		err = json.Unmarshal([]byte(*order.RideInfo), &reqArr)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		order.RideInfoArr = &reqArr
	}

	newSuccessResponse(c, http.StatusOK, order)
}
func (h *Handler) clientCityOrderCancel(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	orderId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid order_id param")
		return
	}
	var cancelOrRate models.CancelOrRateReasons
	if err = c.Bind(&cancelOrRate); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.ClientOrders.CityOrderChangeStatus(cancelOrRate, orderId, userId, "client_cancelled")
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	err = utils.CancelTaxi(orderId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, "ok")
}
func (h *Handler) clientCityOrderGoingOut(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	orderId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid order_id param")
		return
	}
	err = h.services.ClientOrders.CityOrderChangeStatus(models.CancelOrRateReasons{}, orderId, userId, "client_going_out")
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, "ok")
}
func (h *Handler) clientCityOrderRate(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	orderId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid order_id param")
		return
	}
	var cancelOrRate models.CancelOrRateReasons
	if err = c.Bind(&cancelOrRate); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.ClientOrders.CityOrderChangeStatus(cancelOrRate, orderId, userId, "client_rate")
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, "ok")
}
