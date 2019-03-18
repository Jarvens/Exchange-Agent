// date: 2019-03-15
package cache

import (
	"errors"
	"sync"
)

//定义缓存模型
type memoryCache struct {
	//缓存内容
	ca map[string][]byte

	//读写锁 支持多读单写
	sync.RWMutex

	//状态嵌套  匿名字段
	Stat
}

//实现Cache  get 接口
func (m *memoryCache) Get(k string) ([]byte, error) {
	m.Lock()
	defer m.Unlock()
	v, ok := m.ca[k]
	if !ok {
		return nil, errors.New("NotFound")
	}
	return v, nil
}

//实现Cache set接口
func (m *memoryCache) Set(k string, v []byte) error {
	m.Lock()
	defer m.Unlock()
	tmp, ok := m.ca[k]
	if ok {
		m.del(k, tmp)
	}
	m.ca[k] = v
	m.add(k, v)
	return nil
}

//实现Cache  del接口
func (m *memoryCache) Del(k string) error {
	m.Lock()
	defer m.Unlock()
	tmp, ok := m.ca[k]
	if ok {
		delete(m.ca, k)
		m.del(k, tmp)
		return nil
	}
	return nil
}

func (m *memoryCache) GetStat() Stat {
	return m.Stat
}

//实例化
func newInMemoryCache() *memoryCache {
	return &memoryCache{make(map[string][]byte), sync.RWMutex{}, Stat{}}
}
