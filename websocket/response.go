// date: 2019-03-14
package websocket

import (
	"github.com/Jarvens/Exchange-Agent/common"
	"time"
)

type Response struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
	Channel   string `json:"channel"`
}

func Success(channel string) *Response {
	return &Response{Code: common.SUCCESS, Message: "成功", Channel: channel, Timestamp: Now()}
}

func Fail(channel string) *Response {
	return &Response{Code: common.FAIL, Message: "失败", Channel: channel, Timestamp: Now()}
}

func SubRepeat(channel string) *Response {
	return &Response{Message: "重复订阅", Code: common.FAIL, Channel: channel, Timestamp: Now()}
}

func Pong() *Response {
	return &Response{Code: common.SUCCESS, Message: common.PONG, Timestamp: Now()}
}

func EventErr() *Response {
	return &Response{Code: common.FAIL, Message: "事件异常", Timestamp: Now()}
}

func ChannelErr() *Response {
	return &Response{Code: common.FAIL, Message: "频道异常", Timestamp: Now()}
}

func Now() int64 {
	return time.Now().Unix()
}
