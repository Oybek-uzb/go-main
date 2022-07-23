package repository

import (
	"abir/models"
	"abir/pkg/repository/postgres"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateClient(user models.Client, userId int) error
	GetClient(userId int) (models.Client, error)
	CreateOrUpdateClient(user models.User) (int, error)
	GetUser(login, password, userType string) (models.User, error)
	ClientSendCode(login, password string) error
	ClientCheckPhone(phone string) error
	ClientUpdatePhone(userId int, phone string) error

	CreateDriver(user models.Driver, userId int) error
	UpdateDriver(user models.Driver, userId int) error
	SendForModerating(userId int) error
	UpdateDriverCar(car models.DriverCar, userId int) error
	DriverSendCode(login, password string) error
	CreateOrUpdateDriver(user models.User) (int, error)
	GetDriver(userId int) (models.Driver, error)
	GetDriverVerification(userId int) ([]models.DriverVerification, error)
	GetDriverCar(userId int) (models.DriverCar, error)
	GetDriverCarInfo(langId, userId int) (models.DriverCarInfo, error)
	GetDriverInfo(langId, driverId int) (models.Driver, models.DriverCar, models.DriverCarInfo, error)
}

type Utils interface {
	GetColors(langId int) ([]models.Color, error)
	GetCarMarkas() ([]models.CarMarka, error)
	GetCarModels(id int) ([]models.CarModel, error)
	GetRegions(langId int) ([]models.Region, error)
	GetDistricts(langId int, id int) ([]models.District, error)
	GetDistrictsArr(langId int) (map[int]string, error)
	GetTariffs(langId int) (map[int]string, error)
	DriverCancelOrderOptions(langId int) ([]models.DriverCancelOrderOptions, error)
	ClientCancelOrderOptions(langId int, optionType string) ([]models.ClientCancelOrderOptions, error)
}

type SavedAddresses interface {
	Get(userId int) ([]models.SavedAddresses, error)
	Store(address models.SavedAddresses, userId int) error
	Update(address models.SavedAddresses, addressId, userId int) error
	Delete(addressId, userId int) error
}

type CreditCards interface {
	Get(userId int) ([]models.CreditCards, error)
	GetSingleCard(creditCardId, userId int) (models.CreditCards, error)
	Store(creditCard models.CreditCards, userId int) (int, error)
	Activate(creditCardId, userId int) error
	Delete(creditCardId, userId int) error
}

type DriverOrders interface {
	CreateRide(ride models.Ride, userId int) (int, error)
	UpdateRide(ride models.Ride, id, userId int) error
	ChangeRideStatus(id, userId int, status string) error
	RideList(userId int) ([]models.Ride, error)
	RideSingle(id, userId int) (models.Ride, error)
	RideSingleActive(userId int) (models.Ride, error)
	RideSingleNotifications(id, userId int) ([]models.RideNotification, error)
	RideSingleOrderList(id int) ([]models.InterregionalOrder, error)
	RideSingleOrderView(orderId int) (models.InterregionalOrder, error)
	RideSingleOrderAccept(driverId, orderId int) error
	RideSingleOrderCancel(driverId, orderId int) error
	ChatFetch(userId, rideId,orderId int) ([]models.ChatMessages, error)
}

type ClientOrders interface {
	RideList(ride models.Ride, langId, page int) ([]models.ClientRideList, models.Pagination, error)
	RideSingle(langId, id, userId int) (models.ClientRideList, error)
	RideSingleBook(bookRide models.Ride, rideId, userId int) (int, error)
	RideSingleStatus(rideId, userId int) (models.InterregionalOrder, error)
	Activity(userId int, page int, activityType string) ([]models.Activity, models.Pagination, error)
	ChatFetch(userId, rideId,orderId int) ([]models.ChatMessages, error)
	RideSingleCancel(cancelRide models.CanceledOrders, rideId, orderId, userId int) error
}

type DriverSettings interface {
	GetTariffs(userId, langId int) ([]models.DriverTariffs, error)
}

type Repository struct {
	Authorization
	Utils
	SavedAddresses
	CreditCards
	DriverOrders
	ClientOrders
	DriverSettings
}


func NewRepository(dashboard *sqlx.DB, public *sqlx.DB) *Repository {
	return &Repository{
		Authorization: postgres.NewAuthPostgres(public, dashboard),
		Utils: postgres.NewUtilsPostgres(public, dashboard),
		SavedAddresses: postgres.NewSavedAddressesPostgres(public),
		CreditCards: postgres.NewCreditCardsPostgres(public),
		DriverOrders: postgres.NewDriverOrdersPostgres(public),
		ClientOrders: postgres.NewClientOrdersPostgres(public),
		DriverSettings: postgres.NewDriverSettingsPostgres(public, dashboard),
	}
}
