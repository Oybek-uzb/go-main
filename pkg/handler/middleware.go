package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	languageHeader = "Lang"
	userCtx = "userId"
	userLangCtx = "langId"
)
var langs = map[string]int{
	"uz": 1,
	"ru": 2,
}

func (h *Handler) userIdentity(c *gin.Context)  {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer"{
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	if len(headerParts[1]) == 0 {
		newErrorResponse(c, http.StatusUnauthorized, "token is empty")
		return
	}

	userId, _, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.Set(userCtx, userId)
}

func (h *Handler) driverIdentity(c *gin.Context)  {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer"{
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	if len(headerParts[1]) == 0 {
		newErrorResponse(c, http.StatusUnauthorized, "token is empty")
		return
	}

	userId, userType, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	if userType != "driver" {
		newErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	c.Set(userCtx, userId)
}

func (h *Handler) clientIdentity(c *gin.Context)  {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer"{
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	if len(headerParts[1]) == 0 {
		newErrorResponse(c, http.StatusUnauthorized, "token is empty")
		return
	}

	userId, userType, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	if userType != "client" {
		newErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	c.Set(userCtx, userId)
}


func (h *Handler) language(c *gin.Context)  {
	header := c.GetHeader(languageHeader)
	if header == "" {
		c.Set(userLangCtx, langs["ru"])
		return
	}
	if header != "uz" && header != "ru" {
		c.Set(userLangCtx, langs["ru"])
		return
	}
	c.Set(userLangCtx, langs[header])
}

func JSONMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	}
}
func getLangId(c *gin.Context) (int, error) {
	id, ok := c.Get(userLangCtx)
	if !ok {
		return langs["ru"], nil
	}
	idInt, ok := id.(int)
	if !ok {
		return langs["ru"], nil
	}
	return idInt, nil
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		return 0, errors.New("user id is invalid type")
	}
	return idInt, nil
}
