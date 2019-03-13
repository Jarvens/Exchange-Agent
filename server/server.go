// date: 2019-03-13
package server

import (
	"github.com/Jarvens/Exchange-Agent/api"
	"github.com/Jarvens/Exchange-Agent/config"
	"github.com/Jarvens/Exchange-Agent/util/log"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			// c.Abort(200)
			c.Abort()
			return
		}
		c.Next()
	}
}

func Run() {
	r := gin.New()

	switch config.Environment {
	case "DEVELOPMENT":
		r.Use(log.AccessLogger())
	case "TEST":
		r.Use(log.AccessLogger())
	case "PRODUCTION":
		r.Use(log.AccessLogger())
	}
	r.Use(gin.Recovery())
	r.Use(CORSMiddleware())
	api.RouteAPI(r)
	r.Run(":3001")
}
