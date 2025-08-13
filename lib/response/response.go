package response

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(code, Response{Data: nil, Message: message})
}

func Json(c *gin.Context, code int, data interface{}) {
	c.JSON(code, Response{Data: data, Message: "OK"})
}

type ResponsePagination struct {
	Data interface{} `json:"data"`
	Meta Meta        `json:"meta"`
}

type Meta struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

func JsonPagination(c *gin.Context, code int, data interface{}, page, limit int, total int64) {
	c.JSON(code, Response{
		Data: ResponsePagination{
			Data: data,
			Meta: Meta{
				Page:  page,
				Limit: limit,
				Total: total,
			},
		},
		Message: "OK",
	})
}
