package routers

import (
	"github.com/gin-gonic/gin"
	"go-swan/common"
	"net/http"
	"time"
)

func GetSystemTime(c *gin.Context) {
	c.JSON(http.StatusOK, common.CreateSuccessResponse(time.Now().UnixNano()))
}
