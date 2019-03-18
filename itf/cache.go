// date: 2019-03-15
package itf

import "github.com/Jarvens/Exchange-Agent/cache"

type Cache interface {

	//查询缓存值 key
	Get(k string) ([]byte, error)

	//设置缓存  key-value
	Set(k string, v []byte) error

	//删除缓存 key
	Del(k string) error

	GetStat() cache.Stat
}
