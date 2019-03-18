// date: 2019-03-15
package cache

import (
	"github.com/Jarvens/Exchange-Agent/itf"
	"log"
)

func New(typ string) itf.Cache {
	var c itf.Cache
	if typ == "inMemory" {
		c = newInMemoryCache()
	}
	if c == nil {
		panic("未知缓存类型:" + typ)
	}
	log.Println(typ, "缓存对象实例化失败")
	return c
}
