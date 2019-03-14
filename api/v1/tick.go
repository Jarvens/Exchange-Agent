// date: 2019-03-13
package v1

import (
	"github.com/Jarvens/Exchange-Agent/api/response"
	"github.com/gin-gonic/gin"
)

func Tick(parentRoute *gin.RouterGroup) {
	//服务分组
	route := parentRoute.Group("/tick", gin.BasicAuth(gin.Accounts{"Admin": "Admin"}))
	route.GET("/all", allTick)
}

func allTick(c *gin.Context) {
	//TODO service 调用
	messageTypes := &response.MessageTypes{
		OK:                  "registration.done",
		Unauthorized:        "login.error.fail",
		NotFound:            "registration.error.fail",
		InternalServerError: "registration.error.fail",
	}
	messages := &response.Messages{OK: "User is registered successfully."}
	response.JSON(c, 200, messageTypes, messages, nil)
}
