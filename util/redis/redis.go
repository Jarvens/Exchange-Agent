// date: 2019-03-13
package redis

import (
	"context"
	"github.com/Jarvens/Exchange-Agent/util/config"
	"github.com/garyburd/redigo/redis"
	"github.com/youtube/vitess/go/pools"
	"time"
)

var Pool, Resource, InitErr = RedisInit()

type ResourceConn struct {
	redis.Conn
}

func (r ResourceConn) Close() {
	r.Conn.Close()
}

func RedisInit() (*pools.ResourcePool, ResourceConn, error) {
	var resourceConn ResourceConn
	pool := pools.NewResourcePool(func() (pools.Resource, error) {
		c, err := redis.Dial("websocket", config.RedisAddr())
		return ResourceConn{c}, err
	}, config.Capacity, config.MaxCap, time.Minute)
	ctx := context.TODO()
	r, err := pool.Get(ctx)
	if err != nil {
		return pool, resourceConn, err
	}
	defer pool.Put(r)
	return pool, resourceConn, nil
}

func (*ResourceConn) Append(key string, value string) error {
	if InitErr != nil {
		return InitErr
	}
	_, err := Resource.Do("SET", key, value)
	return err
}

func (*ResourceConn) Get(key string) (string, error) {
	var value string
	var err error
	if InitErr != nil {
		return value, InitErr
	}
	value, err = redis.String(Resource.Do("GET", key))
	return value, err
}

func (*ResourceConn) Del(key string) error {
	if InitErr != nil {
		return InitErr
	}
	_, err := Resource.Do("DEL", key)
	return err
}
