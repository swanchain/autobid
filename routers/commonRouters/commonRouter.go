package commonRouters

import (
	"github.com/gin-gonic/gin"
	"go-swan/common"
	"go-swan/common/constants"
	"net/http"
	"time"
)

func HostManager(router *gin.RouterGroup) {
	router.GET(constants.URL_HOST_GET_HOST_INFO, GetSwanMinerVersion)
	router.GET(constants.URL_HOST_GET_HEALTH_CHECK, GetSystemTime)
}

func GetSwanMinerVersion(c *gin.Context) {
	info := getSwanMinerHostInfo()
	c.JSON(http.StatusOK, common.CreateSuccessResponse(info))
}

func GetSystemTime(c *gin.Context) {
	curTime := time.Now().UnixNano()
	c.JSON(http.StatusOK, common.CreateSuccessResponse(curTime))
}
