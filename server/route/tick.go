// date: 2019-03-13
package route

import (
	"github.com/Jarvens/Exchange-Agent/util/config"
	"github.com/gin-gonic/gin"
)

func Tick(parentRoute *gin.RouterGroup) {
	route := parentRoute
	route.GET("/log/access", func(context *gin.Context) {
		context.File(config.AccessLogFilePath + config.AccessLogFileExtension)
	})
}
