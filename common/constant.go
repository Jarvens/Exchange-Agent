// date: 2019-03-07
package common

import "sync"

//quote.kline.1m.btc_usdt
//quote.tick.btc_usdt
//quote.depth.btc_usdt

type SocketMap struct {
	sync.Mutex
	ConnMap map[string]map[string][]string
	//缺少连接属性
}

var (
	Smap      = &SocketMap{sync.Mutex{}, make(map[string]map[string][]string)}
	KlineChan = "quote.kline.%s.%s"
	TickChan  = "quote.tick.%s"
	DepthChan = "quote.depth.%s"
)

const (
	PING    = "ping"
	PONG    = "pong"
	SUB     = "sub"
	UN_SUB  = "un_sub"
	TICK    = "tick"
	DEPTH   = "depth"
	KLINE   = "kline"
	SUCCESS = 0
	FAIL    = 1
)
