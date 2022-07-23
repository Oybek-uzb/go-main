package handler

import (
	"abir/models"
	"abir/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) clientSignUp(c *gin.Context)  {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var input models.Client

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = input.ValidateCreate()
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}
	err = h.services.Authorization.CreateClient(c.Request.Context(),input,userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	//client, err := h.services.GetClient(userId)
	//if err != nil {
	//	newErrorResponse(c, http.StatusInternalServerError, err.Error())
	//	return
	//}
	newSuccessResponse(c, http.StatusOK, "ok")
}

func (h *Handler) clientUpdate(c *gin.Context)  {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var input models.Client

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = input.ValidateUpdate()
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}
	err = h.services.Authorization.CreateClient(c.Request.Context(),input,userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	//client, err := h.services.GetClient(userId)
	//if err != nil {
	//	newErrorResponse(c, http.StatusInternalServerError, err.Error())
	//	return
	//}
	newSuccessResponse(c, http.StatusOK, "ok")
}
type UpdatePhoneSendCodeInput struct {
	Phone string `json:"phone" form:"phone" binding:"required"`
}
func (h *Handler) clientUpdatePhoneSendCode(c *gin.Context){
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var input UpdatePhoneSendCodeInput
	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	client, err := h.services.GetClient(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "client not found")
		return
	}
	phone := utils.StripString(input.Phone)
	if *client.Phone == phone {
		newErrorResponse(c, http.StatusBadRequest, "you can't change your phone number")
		return
	}
	err = h.services.Authorization.SendActivationCode(userId, phone)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, map[string]string{
		"phone": input.Phone,
	})
}
type UpdatePhoneInput struct {
	Phone string `json:"phone" form:"phone" binding:"required"`
	Code string `json:"code" form:"code" binding:"required"`
}
func (h *Handler) clientUpdatePhone(c *gin.Context){
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var input UpdatePhoneInput
	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	phone := utils.StripString(input.Phone)
	err = h.services.Authorization.ClientUpdatePhone(userId, phone, input.Code)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, "ok")
}

func (h *Handler) clientGetMe(c *gin.Context)  {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	client, err := h.services.GetClient(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "client not found")
		return
	}
	newSuccessResponse(c, http.StatusOK, client)
}
func (h *Handler) driverGetMe(c *gin.Context)  {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	driver, err := h.services.GetDriver(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "driver not found")
		return
	}
	newSuccessResponse(c, http.StatusOK, driver)
}
func (h *Handler) driverVerification(c *gin.Context)  {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	verification, err := h.services.GetDriverVerification(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, verification)
}
type signInInput struct {
	Login string `json:"login" form:"login" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}
func (h *Handler) clientSignIn(c *gin.Context)  {
	var input signInInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	token, err := h.services.Authorization.GenerateToken(input.Login, input.Password, "client")
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	userId, _, err := h.services.ParseToken(token)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	client, err := h.services.GetClient(userId)
	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
		"registered": client.Id != 0,
	})
}
func (h *Handler) driverSignIn(c *gin.Context)  {
	var input signInInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	token, err := h.services.Authorization.GenerateToken(input.Login, input.Password, "driver")
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	userId, _, err := h.services.ParseToken(token)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	driver, err := h.services.GetDriver(userId)
	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
		"registered": driver.Id != 0,
	})
}

type sendCodeInput struct {
	Login string `form:"login" json:"login" binding:"required"`
}
func (h *Handler) clientSendCode(c *gin.Context)  {
	var input sendCodeInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err := h.services.Authorization.ClientSendCode(input.Login)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, "Code successfully sent")
}

func (h *Handler) driverSendCode(c *gin.Context)  {
	var input sendCodeInput
	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err := h.services.Authorization.DriverSendCode(input.Login)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, "Code successfully sent")
}


func (h *Handler) driverSignUp(c *gin.Context)  {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var input models.Driver
	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = input.ValidateCreate()
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}
	err = h.services.Authorization.CreateDriver(c.Request.Context(),input,userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	//client, err := h.services.GetDriver(userId)
	//if err != nil {
	//	newErrorResponse(c, http.StatusInternalServerError, err.Error())
	//	return
	//}
	newSuccessResponse(c, http.StatusOK, "ok")
}

func (h *Handler) driverSignUpSendForModerating(c *gin.Context){
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.services.Authorization.SendForModerating(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, "ok")
}

func (h *Handler) driverUpdate(c *gin.Context)  {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var input models.Driver

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = input.ValidateUpdate()
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}
	err = h.services.Authorization.UpdateDriver(c.Request.Context(),input,userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, "ok")
}
func (h *Handler) driverCarUpdate(c *gin.Context)  {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var input models.DriverCar

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.services.Authorization.UpdateDriverCar(c.Request.Context(),input,userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newSuccessResponse(c, http.StatusOK, "ok")
}
func (h *Handler) driverCarFetch(c *gin.Context)  {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	langId, err := getLangId(c)
	if err != nil {
		return
	}
	client, err := h.services.GetDriverCar(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	carInfo, err := h.services.GetDriverCarInfo(langId, userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	client.CarInfo = &carInfo
	newSuccessResponse(c, http.StatusOK, client)
}