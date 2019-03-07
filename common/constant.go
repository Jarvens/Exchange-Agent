// date: 2019-03-07
package common

import "sync"

//quote.kline.1m.btc_usdt
//quote.tick.btc_usdt
//quote.depth.btc_usdt

type SocketMap struct {
	Lock    sync.Mutex
	ConnMap map[string]map[string][]string
}

var GsMap = SocketMap{Lock: sync.Mutex{}, ConnMap: make(map[string]map[string][]string)}

var KlineChan = "quote.kline.%s.%s"
var TickChan = "quote.tick.%s"
var DepthChan = "quote.depth.%s"

const (
	Ping        = "ping"
	Subscribe   = "subscribe"
	UbSubscribe = "ub_subscribe"
	Tick        = "tick"
	Depth       = "depth"
	Kline       = "kline"
	Success     = 0
	Fail        = 1
)
