// date: 2019-03-13
package config

import "github.com/Jarvens/Exchange-Agent/config"

const (
	redisAddrDev = "127.0.0.1"
	redisPortDev = "6379"
	redisAddrPro = "127.0.0.1"
	redisPortPro = "6379"
	Capacity     = 1
	MaxCap       = 2
)

//构造redis地址
func RedisAddr() string {
	var redisAddr string
	switch config.Environment {
	case "DEVELOPMENT":
		redisAddr = redisAddrDev + ":" + redisPortDev
	case "TEST":
	//TODO 测试配置
	default:
		redisAddr = redisAddrPro + ":" + redisPortPro
	}
	return redisAddr
}
