// date: 2019-03-15
package quote

import (
	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
	"io"
	"sync"
)

//连接模型
type Connection struct {
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

	id string
}

//分配ID
func (c *Connection) GetID() string {
	c.once.Do(func() {
		u, _ := uuid.NewV1()
		c.id = u.String()
	})
	return c.id
}
