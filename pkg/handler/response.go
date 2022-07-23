package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Content *string `json:"content"`
}
type successResponse struct {
	Success bool `json:"success"`
	Message *string `json:"message"`
	Content interface{} `json:"content"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{false,message, nil})
}

func newSuccessResponse(c *gin.Context, statusCode int, content interface{}) {
	c.AbortWithStatusJSON(statusCode, successResponse{true, nil,content})
}


