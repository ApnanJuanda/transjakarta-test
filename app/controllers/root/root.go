package root

import (
	"github.com/ApnanJuanda/transjakarta/lib/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(context *gin.Context) {
	response.Json(context, http.StatusOK, nil)
}
