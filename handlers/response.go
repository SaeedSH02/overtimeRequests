package handlers

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Reason  string `json:"reason,omitempty"` //if it's empty, don't include it in JSON
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func SendErrorResponse(c *gin.Context, statusCode int, message, reason string){
	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Message: message,
		Reason: reason,
	})
}

func SendSuccessResponse(c *gin.Context, statusCode int, data interface{}){
	c.JSON(statusCode, SuccessResponse{
		Success: true,
		Data: data,
	})
}