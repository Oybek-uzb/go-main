package service

import (
	"abir/models"
	"abir/pkg/repository"
	"abir/pkg/storage"
	"abir/pkg/utils"
	"context"

	"github.com/go-redis/redis"
	"github.com/streadway/amqp"
)

type Authorization interface {
	CreateClient(ctx context.Context, user models.Client, userId int) error
	GetClient(userId int) (models.Client, error)
	GenerateToken(username, password, userType string) (string, error)
	ClientSendCode(login string) error
	ParseToken(token string) (int, string, error)
	ClientSendActivationCode(userId int, phone string) error
	DriverSendActivationCode(userId int, phone string) error
	ClientUpdatePhone(userId int, phone, code string) error
	PutFirebaseToken(userId int, firebaseToken string) error

	DriverSendCode(login string) error
	CreateDriver(ctx context.Context, user models.Driver, userId int) error
	UpdateDriver(ctx context.Context, user models.Driver, userId int) error
	SendForModerating(userId int) error
	UpdateDriverCar(ctx context.Context, car models.DriverCar, userId int) error
	GetDriver(userId int) (models.Driver, error)
	GetTodayInfo(userId int, date models.CurrentDate) (models.TodayInfo, error)
	GetDriverId(userId int) (int, error)
	GetDriverVerification(userId int) ([]models.DriverVerification, error)
	GetDriverCar(userId int) (models.DriverCar, error)
	GetDriverCarInfo(langId, userId int) (models.DriverCarInfo, error)
	GetDriverInfo(langId, driverId int) (models.Driver, models.DriverCar, models.DriverCarInfo, error)
	DriverUpdatePhone(userId int, phone, code string) error
}

type Utils interface {
	GetColors(langId int) ([]models.Color, error)
	GetCarMarkas() ([]models.CarMarka, error)
	GetCarModels(id int) ([]models.CarModel, error)
	Test(ctx context.Context) error
	GetRegions(langId int) ([]models.Region, error)
	GetDistricts(langId int, id int) ([]models.District, error)
	GetDistrictsArr(langId int) (map[int]string, error)
	GetTariffs(langId int) (map[string]string, error)
	DriverCancelOrderOptions(langId int) ([]models.DriverCancelOrderOptions, error)
	ClientCancelOrderOptions(langId int, optionType string) ([]models.ClientCancelOrderOptions, error)
	ClientRateOptions(langId, rate int, optionType string) ([]models.ClientRateOptions, error)
}

type SavedAddresses interface {
	Get(userId int) ([]models.SavedAddresses, error)
	Store(address models.SavedAddresses, userId int) error
	Update(address models.SavedAddresses, addressId, userId int) error
	Delete(addressId, userId int) error
}

type CreditCards interface {
	Get(userId int) ([]models.CreditCards, error)
	Store(creditCard models.CreditCards, userId int) (int, error)
	SendActivationCode(creditCardId, userId int) (string, error)
	Activate(creditCardId, userId int, code string) error
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
	ChatFetch(userId, rideId, orderId int) ([]models.ChatMessages, error)
	CityOrderView(orderId, userId int) (models.CityOrder, error)
	CityOrderChangeStatus(req models.CityOrderRequest, cancelOrRate models.CancelOrRateReasons, orderId, userId int, status string) error
	CalculateRouteAmount(points [][2]float64, tariffId int) (map[string]interface{}, error)
	CityTariffInfo(points string, tariffId int) (models.TariffInfo, error)
}
type ClientOrders interface {
	RideList(ride models.Ride, langId, page int) ([]models.ClientRideList, models.Pagination, error)
	RideSingle(langId, id, userId int) (models.ClientRideList, error)
	RideSingleBook(bookRide models.Ride, rideId, userId int) (int, error)
	RideSingleStatus(rideId, userId int) (models.InterregionalOrder, error)
	Activity(userId int, page int, activityType, orderType string) ([]models.Activity, models.Pagination, error)
	ChatFetch(userId, rideId, orderId int) ([]models.ChatMessages, error)
	RideSingleCancel(cancelRide models.CancelOrRateReasons, rideId, orderId, userId int) error
	CityTariffs(points [][2]float64, langId int) ([]models.CityTariffs, error)
	CityNewOrder(order models.CityOrder, userId int) (int, error)
	CityOrderView(orderId, userId int) (models.CityOrder, error)
	CityOrderChangeStatus(cancelOrRate models.CancelOrRateReasons, orderId, userId int, status string) error
	CityOrderChange(points string, orderId int) error
}

type DriverSettings interface {
	GetTariffs(userId, langId int) ([]models.DriverTariffs, error)
	TariffsEnable(userId, tariffId int, isActive bool) error
	SetOnline(userId int, isActive int) error
}

type RentCars interface {
	GetCategoriesList(langId int) ([]models.CarCategory, error)
	GetCarsByCategoryId(categoryId int) ([]models.CarByCategoryId, error)
	GetCarByCategoryIdCarId(categoryId, carId, langId int) (models.Car, error)
	PostRentCarByCarId(userId, carId int, rentCarDetails models.RentCarDetails) (int, error)
	GetCompaniesList() ([]models.CarCompany, error)
	GetCarsByCompanyId(companyId int) (models.CarCompanyDetails, error)
	GetCarByCompanyIdCarId(companyId, carId, langId int) (models.Car, error)
	GetMyCompaniesList(userId int) ([]models.CarCompany, error)
	GetMyCompanyById(userId, companyId int) (models.CarCompany, error)
	GetMyCarParkByCompanyId(userId, companyId int, inDiscount bool) ([]models.Car, error)
	PostMyCompany(ctx context.Context, userId int, company models.RentMyCompanyCreate) (int, error)
	PostMyCar(ctx context.Context, userId, companyId int, car models.CarCreate) (int, error)
	PutMyCar(ctx context.Context, userId, carId, companyId int, car models.CarCreate) (int, error)
	DeleteMyCar(userId, carId, myCompanyId int) (int, error)

	GetCategoriesForEvents(langId int) ([]models.CarCategory, error)
	GetMyCarsForEvents(userId, langId int) ([]models.MyCarForEvents, error)
	PostMyCarForEvents(ctx context.Context, userId int, car models.CarCreate) (int, error)
}

type Service struct {
	Authorization
	Utils
	SavedAddresses
	CreditCards
	DriverOrders
	ClientOrders
	DriverSettings
	RentCars
}

func NewService(repos *repository.Repository, client *redis.Client, storage storage.Storage, ch *amqp.Channel, fcmClient *utils.FCMClient) *Service {
	return &Service{
		Authorization:  NewAuthService(repos.Authorization, client, storage, fcmClient),
		Utils:          NewUtilsService(repos.Utils, fcmClient),
		SavedAddresses: NewSavedAddressesService(repos.SavedAddresses, fcmClient),
		CreditCards:    NewCreditCardsService(repos.CreditCards, client, fcmClient),
		DriverOrders:   NewDriverOrdersService(repos.DriverOrders, ch, fcmClient),
		ClientOrders:   NewClientOrdersService(repos.ClientOrders, ch, fcmClient),
		DriverSettings: NewDriverSettingsService(repos.DriverSettings, fcmClient),
		RentCars:       NewClientRentService(repos.RentCars, ch, storage, fcmClient),
	}
}
