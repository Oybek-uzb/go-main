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

func (h *Handler) InitRoutes() *gin.Engine{
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
				}
				carProtected := protected.Group("/car")
				{
					carProtected.POST("/update", h.driverCarUpdate)
					carProtected.GET("/fetch", h.driverCarFetch)
				}
				settingsProtected := protected.Group("/settings")
				{
					settingsProtected.GET("/tariffs", h.driverCityTariffs)
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
				utilsProtected.GET("/driver-cancel-order-options", h.utilsDriverCancelOrderOptions)
				utilsProtected.GET("/client-cancel-order-options", h.utilsClientCancelOrderOptions)
			}
		}
	}
	return router
}