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
	"strings"
)

type ClientOrdersService struct {
	repo repository.ClientOrders
	ch   *amqp.Channel
}

func NewClientOrdersService(repo repository.ClientOrders, ch *amqp.Channel) *ClientOrdersService {
	return &ClientOrdersService{repo: repo, ch: ch}
}

func (s *ClientOrdersService) RideList(ride models.Ride, langId, page int) ([]models.ClientRideList, models.Pagination, error) {
	return s.repo.RideList(ride, langId, page)
}
func (s *ClientOrdersService) RideSingle(langId, id, userId int) (models.ClientRideList, error) {
	return s.repo.RideSingle(langId, id, userId)
}
func (s *ClientOrdersService) RideSingleBook(bookRide models.Ride, rideId, userId int) (int, error) {
	return s.repo.RideSingleBook(bookRide, rideId, userId)
}

func (s *ClientOrdersService) RideSingleStatus(rideId, userId int) (models.InterregionalOrder, error) {
	return s.repo.RideSingleStatus(rideId, userId)
}

func (s *ClientOrdersService) Activity(userId, page int, activityType, orderType string) ([]models.Activity, models.Pagination, error) {
	return s.repo.Activity(userId, page, activityType, orderType)
}

func (s *ClientOrdersService) ChatFetch(userId, rideId, orderId int) ([]models.ChatMessages, error) {
	return s.repo.ChatFetch(userId, rideId, orderId)
}

func (s *ClientOrdersService) RideSingleCancel(cancelRide models.CancelOrRateReasons, rideId, orderId, userId int) error {
	return s.repo.RideSingleCancel(cancelRide, rideId, orderId, userId)
}

func (s *ClientOrdersService) CityTariffs(points [][2]float64, langId int) ([]models.CityTariffs, error) {
	districtId, err := utils.GetMyDistrictId(points[0][1], points[0][0])
	if err != nil {
		logrus.Errorf(err.Error())
		return []models.CityTariffs{}, errors.New("error while getting district id")
	}
	inside, outside, err := utils.CalculateRoute(points)
	if err != nil {
		logrus.Errorf(err.Error())
		return []models.CityTariffs{}, errors.New("error while calculating")
	}
	lists, err := s.repo.CityTariffs(districtId, langId)
	for i, list := range lists {
		startPrice := 0
		pricePerKm := 0
		pricePerKmOuter := 0
		ifNil := 0
		if list.StartPrice != nil {
			startPrice = *list.StartPrice
		} else {
			lists[i].StartPrice = &ifNil
		}
		if list.PricePerKm != nil {
			pricePerKm = *list.PricePerKm
		} else {
			lists[i].PricePerKm = &ifNil
		}
		if list.PricePerKmOuter != nil {
			pricePerKmOuter = *list.PricePerKmOuter
		} else {
			lists[i].PricePerKmOuter = &ifNil
		}
		price := float64(startPrice) + (float64(inside)/1000)*float64(pricePerKm) + (float64(outside)/1000)*float64(pricePerKmOuter)
		price = math.Round(price/100) * 100
		priceInt := int(price)
		lists[i].Price = &priceInt

		if list.Icon != nil {
			lists[i].Icon = utils.GetFileUrl(strings.Split(*list.Icon, "/"))
		}
		if list.Image != nil {
			lists[i].Image = utils.GetFileUrl(strings.Split(*list.Image, "/"))
		}
	}
	return lists, err
}

func (s *ClientOrdersService) CityNewOrder(order models.CityOrder, userId int) (int, error) {
	return s.repo.CityNewOrder(order, userId)
}
func (s *ClientOrdersService) CityOrderView(orderId, userId int) (models.CityOrder, error) {
	return s.repo.CityOrderView(orderId, userId)
}

func (s *ClientOrdersService) CityOrderChangeStatus(cancelOrRate models.CancelOrRateReasons, orderId, userId int, status string) error {
	driverId, err := s.repo.CityOrderChangeStatus(cancelOrRate, orderId, userId, status)
	if err == nil {
		if status == "client_going_out" {
			info := models.DriverOrderSocket{
				Id:       orderId,
				DriverId: driverId,
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

func (s *ClientOrdersService) CityOrderChange(points string, orderId int) error {
	driverId, err := s.repo.CityOrderChange(points, orderId)
	if err != nil {
		return err
	}

	if err == nil {
		info := models.DriverOrderSocket{
			Id:       orderId,
			DriverId: driverId,
			Status:   "order_changed",
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

	return nil
}
