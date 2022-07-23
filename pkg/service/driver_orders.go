package service

import (
	"abir/models"
	"abir/pkg/repository"
)

type DriverOrdersService struct {
	repo repository.DriverOrders
}

func NewDriverOrdersService(repo repository.DriverOrders) *DriverOrdersService {
	return &DriverOrdersService{repo: repo}
}

func (s *DriverOrdersService) CreateRide(ride models.Ride, userId int) (int, error){
	return s.repo.CreateRide(ride, userId)
}

func (s *DriverOrdersService) RideList(userId int) ([]models.Ride, error){
	return s.repo.RideList(userId)
}

func (s *DriverOrdersService) UpdateRide(ride models.Ride, id, userId int) error{
	return s.repo.UpdateRide(ride, id, userId)
}

func (s *DriverOrdersService) ChangeRideStatus(id, userId int, status string) error{
	return s.repo.ChangeRideStatus(id, userId, status)
}

func (s *DriverOrdersService) RideSingle(id, userId int) (models.Ride, error) {
	return s.repo.RideSingle(id, userId)
}
func (s *DriverOrdersService) RideSingleActive(userId int) (models.Ride, error){
	return s.repo.RideSingleActive(userId)
}

func (s *DriverOrdersService) RideSingleNotifications(id, userId int) ([]models.RideNotification, error){
	return s.repo.RideSingleNotifications(id, userId)
}
func (s *DriverOrdersService) RideSingleOrderList(id int) ([]models.InterregionalOrder, error){
	return s.repo.RideSingleOrderList(id)
}
func (s *DriverOrdersService) RideSingleOrderView(orderId int) (models.InterregionalOrder, error){
	return s.repo.RideSingleOrderView(orderId)
}
func (s *DriverOrdersService) RideSingleOrderAccept(driverId, orderId int) error{
	return s.repo.RideSingleOrderAccept(driverId, orderId)
}
func (s *DriverOrdersService) RideSingleOrderCancel(driverId, orderId int) error{
	return s.repo.RideSingleOrderCancel(driverId, orderId)
}

func (s *DriverOrdersService) ChatFetch(userId, rideId,orderId int) ([]models.ChatMessages, error){
	return s.repo.ChatFetch(userId, rideId,orderId)
}
