package handler

import (
	"abir/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(JSONMiddleware())
	api := router.Group("/api/v1", h.language)
	{
		client := api.Group("/client")
		{
			auth := client.Group("/auth")
			{
				auth.POST("/send-code", h.clientSendCode)
				auth.POST("/sign-in", h.clientSignIn)
			}
			protected := client.Group("", h.clientIdentity)
			{
				authenticate := protected.Group("/auth")
				{
					authenticate.POST("/sign-up", h.clientSignUp)
					authenticate.GET("/me", h.clientGetMe)
					authenticate.POST("/update", h.clientUpdate)
					authenticate.POST("/update-phone/send-code", h.clientUpdatePhoneSendCode)
					authenticate.POST("/update-phone", h.clientUpdatePhone)
				}
				savedAddresses := protected.Group("/saved-addresses")
				{
					savedAddresses.GET("/", h.clientSavedAddressesList)
					savedAddresses.POST("/store", h.clientSavedAddressesStore)
					savedAddresses.POST("/update/:id", h.clientSavedAddressesUpdate)
					savedAddresses.POST("/delete/:id", h.clientSavedAddressesDelete)
				}
				creditCards := protected.Group("/credit-cards")
				{
					creditCards.GET("/", h.clientCreditCardsList)
					creditCards.POST("/store", h.clientCreditCardsStore)
					creditCards.POST("/send-activation-code/:id", h.clientCreditCardsSendActivationCode)
					creditCards.POST("/activate/:id", h.clientCreditCardsActivate)
					creditCards.POST("/delete/:id", h.clientCreditCardsDelete)
				}
				chat := protected.Group("/chat")
				{
					chat.GET("/fetch", h.clientChatFetch)
				}
				orders := protected.Group("/orders")
				{
					activity := orders.Group("/activity")
					{
						activity.GET("/active", h.clientOrdersActivityActive)
						activity.GET("/recently-completed", h.clientOrdersActivityRecentlyCompleted)
						activity.GET("/history", h.clientOrdersActivityHistory)
					}
					ride := orders.Group("/rides")
					{
						ride.GET("/", h.clientOrdersRideList)
						ride.GET("/:id", h.clientOrdersRideSingle)
						ride.POST("/:id/book", h.clientOrdersRideSingleBook)
						ride.GET("/:id/status", h.clientOrdersRideSingleStatus)
						ride.POST("/:id/cancel/:order_id", h.clientOrdersRideSingleCancel)
					}
				}
				city := protected.Group("/city")
				{
					city.POST("/new-order", h.clientCityOrder)
					city.POST("/tariffs", h.clientCityTariffs)
					order := city.Group("/order")
					{
						order.GET("/:id", h.clientCityOrderView)
						order.POST("/:id/cancel", h.clientCityOrderCancel)
						order.POST("/:id/going-out", h.clientCityOrderGoingOut)
						order.POST("/:id/rate", h.clientCityOrderRate)
						order.POST("/:id/change-order-points", h.clientCityOrderChange)
					}
				}
				rentCategories := protected.Group("/rent-categories")
				{
					rentCategories.GET("/", h.rentCategoriesList)
					rentCategories.GET("/:id", h.rentCarsByCategoryId)
					rentCategories.GET("/:id/:car_id", h.rentCarByCategoryIdCarId)
					rentCategories.POST("/:id/:car_id", h.rentCarFromCategoryCreate)
				}
				rentCompanies := protected.Group("/rent-companies")
				{
					rentCompanies.GET("/", h.rentCompaniesList)
					rentCompanies.GET("/:id", h.rentCompanyById)
					rentCompanies.GET("/:id/:car_id", h.rentCarByCompanyIdCarId)
					rentCompanies.POST("/:id/:car_id", h.rentCarFromCompanyCreate)
				}

				rent := protected.Group("/rent")
				{
					myCompanies := rent.Group("/my-companies")
					{
						myCompanies.GET("/", h.myCompaniesList)
						myCompanies.GET("/:id", h.myCompanyById)
						myCompanies.POST("/", h.rentMyCompaniesCreate)
						myCompanies.GET("/:id/my-car-park", h.myCarPark)
						myCompanies.GET("/:id/my-car-park/:car_id", h.myCarByCompanyId)

						announcement := myCompanies.Group("/:id/my-car-park/announcement")
						{
							announcement.GET("/", h.myCarPark)
							announcement.POST("/", h.rentAnnouncementCreate)
							announcement.GET("/:car_id", h.myCarByCompanyId)
							announcement.PUT("/:carId", h.rentAnnouncementUpdate)
							//announcement.DELETE("/:announcement_id", h.rentAnnouncementDelete)
						}
					}
				}
			}
		}
		driver := api.Group("/driver")
		{
			auth := driver.Group("/auth")
			{
				auth.POST("/send-code", h.driverSendCode)
				auth.POST("/sign-in", h.driverSignIn)
			}
			protected := driver.Group("", h.driverIdentity)
			{
				authProtected := protected.Group("/auth")
				{
					authProtected.POST("/sign-up", h.driverSignUp)
					authProtected.GET("/me", h.driverGetMe)
					authProtected.POST("/update", h.driverUpdate)
					authProtected.GET("/verification", h.driverVerification)
					authProtected.POST("/sign-up/send-for-moderating", h.driverSignUpSendForModerating)
					authProtected.POST("/update-phone/send-code", h.driverUpdatePhoneSendCode)
					authProtected.POST("/update-phone", h.driverUpdatePhone)
				}
				carProtected := protected.Group("/car")
				{
					carProtected.POST("/update", h.driverCarUpdate)
					carProtected.GET("/fetch", h.driverCarFetch)
				}
				settingsProtected := protected.Group("/settings")
				{
					settingsProtected.POST("/set-online", h.driverCitySetOnline)
					settingsProtected.GET("/tariffs", h.driverCityTariffs)
					settingsProtected.POST("/tariffs/enable", h.driverCityTariffsEnable)
					settingsProtected.GET("/stats", h.driverStats)
					settingsProtected.GET("/stats/orders", h.driverStatOrders)
				}
				chat := protected.Group("/chat")
				{
					chat.GET("/fetch", h.driverChatFetch)
				}
				orders := protected.Group("/orders")
				{
					orders.POST("/new-order", h.driverOrdersCreateRide)
					ride := orders.Group("/rides")
					{
						order := ride.Group("/order")
						{
							order.GET("/:order_id", h.driverOrdersSingleRideOrderView)
							order.POST("/:order_id/accept", h.driverOrdersSingleRideOrderAccept)
							order.POST("/:order_id/cancel", h.driverOrdersSingleRideOrderCancel)
						}
						ride.GET("/", h.driverOrdersRideList)
						ride.POST("/create", h.driverOrdersCreateRide)
						ride.GET("/active", h.driverOrdersSingleRideActive)
						ride.GET("/:id/view", h.driverOrdersSingleRide)
						ride.POST("/:id/update", h.driverOrdersUpdateRide)
						ride.POST("/:id/start", h.driverOrdersStartRide)
						ride.POST("/:id/complete", h.driverOrdersCompleteRide)
						ride.POST("/:id/cancel", h.driverOrdersCancelRide)
					}
					orders.POST("/calculate-price", h.driverOrdersCalculatePrice)
					city := orders.Group("/city")
					{
						city.GET("/:id", h.driverCityOrderView)
						city.POST("/:id/skip", h.driverCityOrderSkip)
						city.POST("/:id/accept", h.driverCityOrderAccept)
						city.POST("/:id/arrived", h.driverCityOrderArrived)
						city.POST("/:id/start", h.driverCityOrderStart)
						city.POST("/:id/done", h.driverCityOrderDone)
						city.POST("/:id/cancel", h.driverCityOrderCancel)
					}
				}
			}
		}
		utils := api.Group("/utils", h.language)
		{
			utilsProtected := utils.Group("", h.userIdentity)
			{
				utilsProtected.GET("/colors", h.utilsColor)
				utilsProtected.GET("/car-markas", h.utilsCarMarka)
				utilsProtected.GET("/car-markas/:id", h.utilsCarModel)
				utilsProtected.GET("/test", h.test)
				utilsProtected.GET("/regions", h.utilsRegion)
				utilsProtected.GET("/regions/:id", h.utilsDistrict)
				utilsProtected.GET("/client-rate-options", h.utilsClientRateOptions)
				utilsProtected.GET("/driver-cancel-order-options", h.utilsDriverCancelOrderOptions)
				utilsProtected.GET("/client-cancel-order-options", h.utilsClientCancelOrderOptions)
			}
		}
	}
	return router
}
