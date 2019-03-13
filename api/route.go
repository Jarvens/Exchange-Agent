// date: 2019-03-13
package api

import (
	"github.com/Jarvens/Exchange-Agent/api/v1"
	"github.com/Jarvens/Exchange-Agent/config"
	"github.com/gin-gonic/gin"
)

func RouteAPI(parentRoute *gin.Engine) {
	route := parentRoute.Group(config.ApiURL)
	{
		v1.Tick(route)
	}
}
