// date: 2019-03-09
package reactor

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
)

type websocketHandler struct {
	upgrader *websocket.Upgrader

	binder *binder

	//作为扩展使用 计算用户token
	calcUserIdFunc func(token string) (userId string, ok bool)
}

type RequestMessage struct {
	//token
	Token string

	//事件类型  sub/unsub
	Event string

	//事件频道
	Channel string
}

func (wh *websocketHandler) ServeHttp(w http.ResponseWriter, r *http.Request) {
	wsCon, err := wh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer wsCon.Close()

	conn := NewConn(wsCon)

	//Accept request
	conn.AfterReadFunc = func(messageType int, r io.Reader) {
		var msg RequestMessage
		decoder := json.NewDecoder(r)
		if err := decoder.Decode(&msg); err != nil {
			return
		}

		//calculate auth key
		authKey := msg.Token
		if wh.calcUserIdFunc != nil {
			key, ok := wh.calcUserIdFunc(authKey)
			if !ok {
				return
			}
			authKey = key
		}
		//bind  parameter
		wh.binder.Bind(authKey, msg.Event, conn)
	}
	conn.BeforeCloseFunc = func() {
		//unBind
		wh.binder.Unbind(conn)
	}
	conn.Listen()
}

func (wh *websocketHandler) closeConns(uuid, event string) {

}
