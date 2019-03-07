// date: 2019-03-07
package common

//quote.kline.1m.btc_usdt
//quote.tick.btc_usdt
//quote.depth.btc_usdt

var SocketMap = make(map[string]map[string][]string)

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
