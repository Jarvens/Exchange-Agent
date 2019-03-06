// date: 2019-03-06
package tcp

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"reflect"
	"sync"
)

var upgrade = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
}, EnableCompression: true}

var Cmap = make(map[int]*Connection)

type Connection struct {
	wsConn    *websocket.Conn
	inChan    chan []byte
	outChan   chan []byte
	closeChan chan byte
	mutes     sync.Mutex
	isClosed  bool
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	var (
		wsCon *websocket.Conn
		err   error
		conn  *Connection
	)

	if wsCon, err = upgrade.Upgrade(w, r, nil); err != nil {
		fmt.Printf("处理器upgrade失败: %v\n", err)
		return
	}

	if conn, err = NewWebsocket(wsCon); err != nil {
		fmt.Printf("处理器初始化失败: %v\n", err)
		conn.Close()
	}

	for {
		if _, err = conn.Read(); err != nil {
			goto ERR
		}
	}

ERR:
	conn.Close()
}

//生成文件描述符
func makeFd(conn *websocket.Conn) int {
	connVal := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn").Elem()
	tcpConn := reflect.Indirect(connVal).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	return int(pfdVal.FieldByName("Sysfd").Int())
}

func (c *Connection) LoopRead() {
	var (
		data []byte
		err  error
	)
	for {
		if _, data, err = c.wsConn.ReadMessage(); err != nil {
			fmt.Printf("读取消息失败,关闭连接: %v\n", err)
			goto ERR
		}
		fmt.Printf("读取消息: %s\n", string(data))
		select {
		case c.inChan <- data:
		case <-c.closeChan:
			break
		}
	}
ERR:
	c.Close()
}

func (c *Connection) LoopWrite() {
	var (
		data []byte
		err  error
	)
	for {
		select {
		case data = <-c.outChan:
		case <-c.closeChan:
			goto ERR
		}
		if err = c.wsConn.WriteMessage(websocket.TextMessage, data); err != nil {
			goto ERR
		}
	}
ERR:
	c.Close()
}

func NewWebsocket(ws *websocket.Conn) (c *Connection, err error) {
	fd := makeFd(ws)
	fmt.Printf("文件描述符: %d\n", fd)
	conn := &Connection{wsConn: ws, inChan: make(chan []byte, 1024), outChan: make(chan []byte, 1024), closeChan: make(chan byte, 1)}
	Cmap[fd] = conn
	conn.onConnected()
	go conn.LoopRead()
	go conn.LoopWrite()
	return conn, nil
}

func (c *Connection) Close() {
	c.wsConn.Close()
	c.mutes.Lock()
	defer c.mutes.Unlock()
	if !c.isClosed {
		close(c.closeChan)
		c.isClosed = true
	}
}

func (c *Connection) Read() (data []byte, err error) {
	select {
	case data = <-c.inChan:
		c.Write([]byte("test"))
	case <-c.closeChan:
		fmt.Println("连接关闭")
		return nil, errors.New("连接关闭")
	}
	return
}

func (c *Connection) Write(data []byte) (err error) {
	if err = c.wsConn.WriteMessage(websocket.TextMessage, data); err != nil {
		fmt.Printf("回写消息失败: %v\n", err)
		return err
	}
	fmt.Printf("回写内容: %s\n", string(data))
	return
}

func (c *Connection) onConnected() {
	fmt.Printf("[Exchange-Agent]%s 加入会话\n", c.wsConn.RemoteAddr().String())
}
