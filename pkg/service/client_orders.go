package service

import (
	"abir/models"
	"abir/pkg/repository"
)

type ClientOrdersService struct {
	repo repository.ClientOrders
}

func NewClientOrdersService(repo repository.ClientOrders) *ClientOrdersService {
	return &ClientOrdersService{repo: repo}
}

func (s *ClientOrdersService) RideList(ride models.Ride, langId, page int) ([]models.ClientRideList, models.Pagination, error) {
	return s.repo.RideList(ride, langId, page)
}
func (s *ClientOrdersService) RideSingle(langId, id, userId int) (models.ClientRideList, error) {
	return s.repo.RideSingle(langId, id, userId)
}
func (s *ClientOrdersService) RideSingleBook(bookRide models.Ride, rideId, userId int) (int, error){
	return s.repo.RideSingleBook(bookRide, rideId, userId)
}

func (s *ClientOrdersService) RideSingleStatus(rideId, userId int) (models.InterregionalOrder, error){
	return s.repo.RideSingleStatus(rideId, userId)
}

func (s *ClientOrdersService) Activity(userId, page int, activityType string) ([]models.Activity, models.Pagination, error){
	return s.repo.Activity(userId, page, activityType)
}

func (s *ClientOrdersService) ChatFetch(userId, rideId,orderId int) ([]models.ChatMessages, error){
	return s.repo.ChatFetch(userId, rideId,orderId)
}

func (s *ClientOrdersService) RideSingleCancel(cancelRide models.CanceledOrders, rideId, orderId, userId int) error{
	return s.repo.RideSingleCancel(cancelRide, rideId, orderId, userId)
}
