package routers

import (
	"github.com/gin-gonic/gin"
	response2 "go-swan/common/response"
	"net/http"
	"time"
)

func GetSystemTime(c *gin.Context) {
	//c.JSON(http.StatusOK, common.CreateSuccessResponse(time.Now().UnixNano()))
	appG := response2.Gin{c}
	appG.Reponse(http.StatusOK, response2.SUCCESS, time.Now().UnixNano())
}
