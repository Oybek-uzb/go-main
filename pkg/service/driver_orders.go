package service

import (
	"abir/models"
	"abir/pkg/repository"
	"abir/pkg/utils"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"math"
	"strconv"
	"strings"
)

type DriverOrdersService struct {
	repo repository.DriverOrders
	ch   *amqp.Channel
}

func NewDriverOrdersService(repo repository.DriverOrders, ch *amqp.Channel) *DriverOrdersService {
	return &DriverOrdersService{repo: repo, ch: ch}
}

func (s *DriverOrdersService) CreateRide(ride models.Ride, userId int) (int, error) {
	return s.repo.CreateRide(ride, userId)
}

func (s *DriverOrdersService) RideList(userId int) ([]models.Ride, error) {
	return s.repo.RideList(userId)
}

func (s *DriverOrdersService) UpdateRide(ride models.Ride, id, userId int) error {
	return s.repo.UpdateRide(ride, id, userId)
}

func (s *DriverOrdersService) ChangeRideStatus(id, userId int, status string) error {
	return s.repo.ChangeRideStatus(id, userId, status)
}

func (s *DriverOrdersService) RideSingle(id, userId int) (models.Ride, error) {
	return s.repo.RideSingle(id, userId)
}
func (s *DriverOrdersService) RideSingleActive(userId int) (models.Ride, error) {
	return s.repo.RideSingleActive(userId)
}

func (s *DriverOrdersService) RideSingleNotifications(id, userId int) ([]models.RideNotification, error) {
	return s.repo.RideSingleNotifications(id, userId)
}
func (s *DriverOrdersService) RideSingleOrderList(id int) ([]models.InterregionalOrder, error) {
	return s.repo.RideSingleOrderList(id)
}
func (s *DriverOrdersService) RideSingleOrderView(orderId int) (models.InterregionalOrder, error) {
	return s.repo.RideSingleOrderView(orderId)
}
func (s *DriverOrdersService) RideSingleOrderAccept(driverId, orderId int) error {
	return s.repo.RideSingleOrderAccept(driverId, orderId)
}
func (s *DriverOrdersService) RideSingleOrderCancel(driverId, orderId int) error {
	return s.repo.RideSingleOrderCancel(driverId, orderId)
}

func (s *DriverOrdersService) ChatFetch(userId, rideId, orderId int) ([]models.ChatMessages, error) {
	return s.repo.ChatFetch(userId, rideId, orderId)
}
func (s *DriverOrdersService) CityOrderView(orderId, userId int) (models.CityOrder, error) {
	return s.repo.CityOrderView(orderId, userId)
}
func (s *DriverOrdersService) CityOrderChangeStatus(req models.CityOrderRequest, cancelOrRate models.CancelOrRateReasons, orderId, userId int, status string) error {
	clientId, err := s.repo.CityOrderChangeStatus(req, cancelOrRate, orderId, userId, status)
	if err == nil {
		if status == "driver_accepted" {
			info := models.ClientOrderSocket{
				Id:       orderId,
				ClientId: clientId,
				Location: req.DriverLastLocation,
				Status:   status,
			}
			infoJson, _ := json.Marshal(info)
			s.ch.Publish(
				"",
				"socket-service",
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        infoJson,
				},
			)
		} else if status == "order_completed" {
			info := models.ClientOrderSocket{
				Id:          orderId,
				ClientId:    clientId,
				OrderAmount: req.OrderAmount,
				Status:      status,
			}
			infoJson, _ := json.Marshal(info)
			s.ch.Publish(
				"",
				"socket-service",
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        infoJson,
				},
			)
		} else {
			info := models.ClientOrderSocket{
				Id:       orderId,
				ClientId: clientId,
				Status:   status,
			}
			infoJson, _ := json.Marshal(info)
			s.ch.Publish(
				"",
				"socket-service",
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        infoJson,
				},
			)
		}
	}
	return err
}
func (s *DriverOrdersService) CityTariffInfo(points string, tariffId int) (models.TariffInfo, error) {
	pointsArr := strings.Split(points, ",")
	lat, err := strconv.ParseFloat(pointsArr[0], 32)
	if err != nil {
		return models.TariffInfo{}, err
	}
	lng, err := strconv.ParseFloat(pointsArr[1], 32)
	if err != nil {
		return models.TariffInfo{}, err
	}
	districtId, err := utils.GetMyDistrictId(lat, lng)
	if err != nil {
		return models.TariffInfo{}, err
	}
	info, err := s.repo.CityTariffInfo(districtId, tariffId)
	ifNil := 0
	if info.StartPrice == nil {
		info.StartPrice = &ifNil
	}
	if info.PricePerKm == nil {
		info.PricePerKm = &ifNil
	}
	if info.PricePerKmOuter == nil {
		info.PricePerKmOuter = &ifNil
	}
	if info.ACPrice == nil {
		info.ACPrice = &ifNil
	}
	if info.Expectation == nil {
		info.Expectation = &ifNil
	}
	return info, err
}

func (s *DriverOrdersService) CalculateRouteAmount(points [][2]float64, tariffId int) (map[string]interface{}, error) {
	districtId, err := utils.GetMyDistrictId(points[0][1], points[0][0])
	if err != nil {
		logrus.Errorf(err.Error())
		return map[string]interface{}{}, errors.New("error while getting district id")
	}
	inside, outside, err := utils.CalculateRoute(points)
	if err != nil {
		logrus.Errorf(err.Error())
		return map[string]interface{}{}, errors.New("error while calculating")
	}
	list, err := s.repo.CityTariff(districtId, tariffId)
	startPrice := 0
	pricePerKm := 0
	pricePerKmOuter := 0
	ifNil := 0
	if list.StartPrice != nil {
		startPrice = *list.StartPrice
	} else {
		list.StartPrice = &ifNil
	}
	if list.PricePerKm != nil {
		pricePerKm = *list.PricePerKm
	} else {
		list.PricePerKm = &ifNil
	}
	if list.PricePerKmOuter != nil {
		pricePerKmOuter = *list.PricePerKmOuter
	} else {
		list.PricePerKmOuter = &ifNil
	}
	if list.ACPrice == nil {
		list.ACPrice = &ifNil
	}
	price := float64(startPrice) + (float64(inside)/1000)*float64(pricePerKm) + (float64(outside)/1000)*float64(pricePerKmOuter)
	price = math.Round(price/100) * 100
	priceInt := int(price)
	priceWithAcInt := int(price) + *list.ACPrice
	list.Price = &priceInt
	list.ACPrice = &priceWithAcInt
	return map[string]interface{}{
		"price":         *list.Price,
		"price_with_ac": *list.ACPrice,
	}, err
}
