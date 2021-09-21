package routers

import (
	"github.com/gin-gonic/gin"
	"go-swan/routers/response"
	"net/http"
	"time"
)

func GetSystemTime(c *gin.Context) {
	appG := response.Gin{c}
	appG.Reponse(http.StatusOK, response.SUCCESS, time.Now().UnixNano())
}
