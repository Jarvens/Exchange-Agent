// date: 2019-03-15
package quote

import (
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"sync"
	"time"
)

type ServerN struct {
	//websocket连接信息
	conn *websocket.Conn

	//读取消息之后处理逻辑
	AfterReadFunc func(messageType int, r io.Reader)

	//关闭连接之前处理逻辑
	BeforeCloseFunc func()

	//用于生成UUID
	once sync.Once

	//关闭信号接收通道
	stopCh chan struct{}

	//返回客户端id 唯一标识
	id string
}

type serverHandler struct {
	*ServerN
}

var upgrade = websocket.Upgrader{EnableCompression: true, CheckOrigin: func(r *http.Request) bool {
	return true
}, HandshakeTimeout: time.Second * 3,
	ReadBufferSize:  1024,
	WriteBufferSize: 1024}

//实现ServeHTTP接口
func (s *serverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	//升级协议
	//将http协议升级为upgrade协议，长连接
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("协议升级失败: ", err)
		return
	}

	//设置关闭事件处理器
	conn.SetCloseHandler(func(code int, text string) error {
		//判断关闭处理器是否为空，如不为空则先执行关闭逻辑
		if s.BeforeCloseFunc != nil {

			//TODO 执行关闭逻辑，此处逻辑是否正确 有待查证

		}

		msg := websocket.FormatCloseMessage(code, "")
		err := conn.WriteControl(websocket.TextMessage, msg, time.Now().Add(time.Second))
		if err != nil {
			return err
		}
		return nil
	})

}

func (s *ServerN) newServerHandler() http.Handler {
	return &serverHandler{s}
}
