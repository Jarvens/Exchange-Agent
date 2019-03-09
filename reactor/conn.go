// date: 2019-03-09
package reactor

import (
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
	"io"
	"sync"
	"time"
)

type Conn struct {
	Conn            *websocket.Conn
	AfterReadFunc   func(messageType int, r io.Reader)
	BeforeCloseFunc func()
	once            sync.Once
	id              string
	stopCh          chan struct{}
}

func (c *Conn) Write(p []byte) (n int, err error) {
	select {
	case <-c.stopCh:
		return 0, errors.New("连接已被关闭,不能继续写入数据")
	default:
		err := c.Conn.WriteMessage(websocket.TextMessage, p)
		if err != nil {
			return 0, nil
		}
		return len(p), nil
	}

	return 0, nil
}

//生成UUID  NewV1 根据时间戳 & 机器 MAC地址生成
func (c *Conn) GetID() string {
	c.once.Do(func() {
		uid, _ := uuid.NewV1()
		c.id = uid.String()
	})
	return c.id
}

//初始化 Conn 结构体
func NewConn(conn *websocket.Conn) *Conn {
	return &Conn{Conn: conn, stopCh: make(chan struct{})}
}

//关闭当前连接
func (c *Conn) Close() error {
	select {
	case <-c.stopCh:
		return errors.New("连接已经关闭")
	default:
		c.Conn.Close()
		close(c.stopCh)
		return nil
	}
}

func (c *Conn) Listen() {
	c.Conn.SetCloseHandler(func(code int, text string) error {
		if c.BeforeCloseFunc != nil {
			c.BeforeCloseFunc()
		}

		if err := c.Close(); err != nil {
			fmt.Printf("关闭连接发生错误: %v\n", err)
		}
		message := websocket.FormatCloseMessage(code, "")
		c.Conn.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
		return nil
	})

ReadLoop:
	for {
		select {
		case <-c.stopCh:
			break ReadLoop
		default:
			messageType, r, err := c.Conn.NextReader()
			if err != nil {
				//TODO 需要处理错误信息
				break ReadLoop
			}

			if c.AfterReadFunc != nil {
				c.AfterReadFunc(messageType, r)
			}
		}
	}
}
