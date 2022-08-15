package handler

import (
	"abir/models"
	"abir/pkg/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) driverOrdersCreateRide(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var ride models.Ride

	if err = c.Bind(&ride); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = ride.ValidateCreate()
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}
	id, err := h.services.DriverOrders.CreateRide(ride, userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, map[string]int{
		"id": id,
	})
}

func (h *Handler) driverOrdersRideList(c *gin.Context) {
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
	districts, err := h.services.Utils.GetDistrictsArr(langId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	lists, err := h.services.DriverOrders.RideList(userId)
	for i, list := range lists {
		fromId, err := strconv.Atoi(list.FromDistrictId)
		if err != nil {
			continue
		}
		toId, err := strconv.Atoi(list.ToDistrictId)
		if err != nil {
			continue
		}
		from := districts[fromId]
		lists[i].From = &from
		to := districts[toId]
		lists[i].To = &to
	}
	newSuccessResponse(c, http.StatusOK, lists)
}
func (h *Handler) driverOrdersSingleRideActive(c *gin.Context) {
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
	districts, err := h.services.Utils.GetDistrictsArr(langId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	list, err := h.services.DriverOrders.RideSingleActive(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	fromId, err := strconv.Atoi(list.FromDistrictId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	toId, err := strconv.Atoi(list.ToDistrictId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	from := districts[fromId]
	list.From = &from
	to := districts[toId]
	list.To = &to
	notifications, err := h.services.DriverOrders.RideSingleNotifications(list.Id, userId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	list.Notifications = &notifications
	newSuccessResponse(c, http.StatusOK, list)
}
func (h *Handler) driverOrdersSingleRide(c *gin.Context) {
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
	districts, err := h.services.Utils.GetDistrictsArr(langId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	list, err := h.services.DriverOrders.RideSingle(id, userId)
	fromId, err := strconv.Atoi(list.FromDistrictId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	toId, err := strconv.Atoi(list.ToDistrictId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	orderList, err := h.services.DriverOrders.RideSingleOrderList(list.Id)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	for i, list := range orderList {
		client, err := h.services.Authorization.GetClient(list.ClientId)
		if err != nil {
			continue
		}
		orderList[i].Client = &client
	}
	from := districts[fromId]
	list.From = &from
	to := districts[toId]
	list.To = &to
	list.OrderList = &orderList
	newSuccessResponse(c, http.StatusOK, list)
}
func (h *Handler) driverOrdersSingleRideOrderView(c *gin.Context) {
	orderId, err := strconv.Atoi(c.Param("order_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid order_id param")
		return
	}
	list, err := h.services.DriverOrders.RideSingleOrderView(orderId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	client, err := h.services.Authorization.GetClient(list.ClientId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	list.Client = &client
	newSuccessResponse(c, http.StatusOK, list)
}
func (h *Handler) driverOrdersSingleRideOrderAccept(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	orderId, err := strconv.Atoi(c.Param("order_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid order_id param")
		return
	}
	err = h.services.DriverOrders.RideSingleOrderAccept(userId, orderId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, "ok")
}
func (h *Handler) driverOrdersSingleRideOrderCancel(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	orderId, err := strconv.Atoi(c.Param("order_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid order_id param")
		return
	}
	err = h.services.DriverOrders.RideSingleOrderCancel(userId, orderId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, "ok")
}
func (h *Handler) driverOrdersUpdateRide(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var ride models.Ride

	if err = c.Bind(&ride); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = ride.ValidateUpdate()
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	err = h.services.DriverOrders.UpdateRide(ride, id, userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, "ok")
}

func (h *Handler) driverOrdersStartRide(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	err = h.services.DriverOrders.ChangeRideStatus(id, userId, "on_the_way")
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, "ok")
}

func (h *Handler) driverOrdersCancelRide(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	err = h.services.DriverOrders.ChangeRideStatus(id, userId, "cancelled")
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, "ok")
}

func (h *Handler) driverOrdersCompleteRide(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	err = h.services.DriverOrders.ChangeRideStatus(id, userId, "done")
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, "ok")
}

func (h *Handler) driverChatFetch(c *gin.Context) {
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
	lists, err := h.services.DriverOrders.ChatFetch(userId, rideId, orderId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, lists)
}
func (h *Handler) driverCityOrderView(c *gin.Context) {
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
	order, err := h.services.DriverOrders.CityOrderView(orderId, userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
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
	client, err := h.services.Authorization.GetClient(order.ClientId)
	order.Client = &client
	newSuccessResponse(c, http.StatusOK, order)
}
func (h *Handler) driverCityOrderSkip(c *gin.Context) {
	orderId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid order_id param")
		return
	}
	err = utils.SkipTaxi(orderId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, "ok")
}
func (h *Handler) driverCityOrderAccept(c *gin.Context) {
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
	var req models.CityOrderRequest
	if err = c.Bind(&req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if req.DriverLastLocation == nil {
		newErrorResponse(c, http.StatusBadRequest, "driver_last_location is required")
		return
	}
	err = h.services.DriverOrders.CityOrderChangeStatus(req, models.CancelOrRateReasons{}, orderId, userId, "driver_accepted")
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	err = utils.AcceptTaxi(orderId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, "ok")
}
func (h *Handler) driverCityOrderArrived(c *gin.Context) {
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
	err = h.services.DriverOrders.CityOrderChangeStatus(models.CityOrderRequest{}, models.CancelOrRateReasons{}, orderId, userId, "driver_arrived")
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, "ok")
}
func (h *Handler) driverCityOrderStart(c *gin.Context) {
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
	err = h.services.DriverOrders.CityOrderChangeStatus(models.CityOrderRequest{}, models.CancelOrRateReasons{}, orderId, userId, "trip_started")
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, "ok")
}
func (h *Handler) driverCityOrderDone(c *gin.Context) {
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
	var req models.CityOrderRequest
	if err = c.Bind(&req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if req.WaitTime == nil || req.WaitTimeAmount == nil || req.RideDistance == nil || req.RideAmount == nil || req.RideTime == nil {
		newErrorResponse(c, http.StatusBadRequest, "fields are required")
		return
	}
	orderAmount := *req.WaitTimeAmount + *req.RideAmount
	commission := (*req.WaitTimeAmount + *req.RideAmount) / 10
	req.Commission = &commission
	req.OrderAmount = &orderAmount
	err = h.services.DriverOrders.CityOrderChangeStatus(req, models.CancelOrRateReasons{}, orderId, userId, "order_completed")
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, "ok")
}
func (h *Handler) driverCityOrderCancel(c *gin.Context) {
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
	err = h.services.DriverOrders.CityOrderChangeStatus(models.CityOrderRequest{}, cancelOrRate, orderId, userId, "driver_cancelled")
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, "ok")
}

func (h *Handler) driverOrdersCalculatePrice(c *gin.Context) {
	var points models.NewPointsRequest
	if err := c.Bind(&points); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if len(points.Points) == 0 {
		newErrorResponse(c, http.StatusBadRequest, "points field is required")
		return
	}
	if points.TariffId == nil {
		newErrorResponse(c, http.StatusBadRequest, "tariff_id field is required")
		return
	}
	list, err := h.services.DriverOrders.CalculateRouteAmount(points.Points, *points.TariffId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, list)
}
