package utils

import (
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/dto"
	"github.com/gin-gonic/gin"
)

type Reponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

func ResponseSuccess(c *gin.Context, code int, message string, data any) {
	response := Reponse{
		Status:  true,
		Message: message,
		Data:    data,
	}
	c.JSON(code, response)

}

func ResponseFailed(c *gin.Context, code int, message string, errors any) {
	response := Reponse{
		Status:  false,
		Message: message,
		Errors:  errors,
	}
	c.JSON(code, response)
}

func ResponsePagination(c *gin.Context, code int, message string, data any, pagination dto.Pagination) {
	response := struct {
		Status     bool           `json:"status"`
		Message    string         `json:"message"`
		Data       any            `json:"data,omitempty"`
		Pagination dto.Pagination `json:"pagination"`
	}{
		Status:     true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	}

	c.JSON(code, response)
}
