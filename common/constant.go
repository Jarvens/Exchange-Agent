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
	Ping        = "ping"
	Pong        = "pong"
	Subscribe   = "subscribe"
	UnSubscribe = "ub_subscribe"
	Tick        = "tick"
	Depth       = "depth"
	Kline       = "kline"
	Success     = 0
	Fail        = 1
)
