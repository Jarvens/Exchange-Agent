// date: 2019-03-09
package reactor

import "github.com/gorilla/websocket"

type websocketHandler struct {
	upgrader *websocket.Upgrader

	binder *binder

	//作为扩展使用 计算用户token
	calcUserIdFunc func(token string) (userId string, ok bool)
}

type RequestMessage struct {
	//token
	Token string

	//事件类型  subscribe / un_subscribe
	Event string

	//事件频道
	Channel string
}
