// date: 2019-03-14
package websocket

import "strings"

type Request struct {
	//订阅 | 取消订阅 | 心跳
	Event string

	// 频道
	Channel string
}

//检查事件
func (r *Request) checkEvent() (string, bool) {

	if r.Event == "" {
		return "", false
	}
	return r.Event, true
}

//检查订阅模块
//market.kline.1m.btc_usdt
func (r *Request) checkModule() (string, string, bool) {
	if r.Channel == "" {
		return "", "", false
	}
	if len(strings.Split(r.Channel, ".")) < 2 {
		return "", "", false
	}
	return r.Channel, strings.Split(r.Channel, ".")[1], true
}
