// date: 2019-03-06
package sys

import (
	"fmt"
	"github.com/gorilla/websocket"
	"golang.org/x/sys/unix"
	"net/http"
	"reflect"
	"sync"
	"time"
)

var epoll *Epoll

var upgrade = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
}, EnableCompression: true, HandshakeTimeout: 5 * time.Second}

type Epoll struct {
	fd       int
	connMap  map[int]*websocket.Conn
	lock     *sync.RWMutex
	inChan   chan []byte
	outChan  chan []byte
	isClosed bool
	event    map[string]interface{}
}

func NewEpoll() (*Epoll, error) {
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	return &Epoll{fd: fd, lock: &sync.RWMutex{}, connMap: make(map[int]*websocket.Conn)}, nil
}

func (e *Epoll) Register(conn *websocket.Conn) error {
	fd := makeFd(conn)
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Events: unix.POLLIN | unix.POLLHUP, Fd: int32(fd)})
	if err != nil {
		fmt.Printf("Epoll 注册失败: %v\n", err)
		return err
	}
	e.lock.Lock()
	defer e.lock.Unlock()
	e.connMap[fd] = conn
	return nil
}

func (e *Epoll) UnRegister(conn *websocket.Conn) error {
	fd := makeFd(conn)
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_DEL, fd, nil)
	if err != nil {
		fmt.Printf("Epoll 取消注册失败: %v\n", err)
		return err
	}
	e.lock.Lock()
	defer e.lock.Unlock()
	delete(e.connMap, fd)
	return nil
}

func makeFd(conn *websocket.Conn) int {
	connVal := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn").Elem()
	tcpConn := reflect.Indirect(connVal).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	return int(pfdVal.FieldByName("Sysfd").Int())
}

func (e *Epoll) Wait() ([]*websocket.Conn, error) {
	events := make([]unix.EpollEvent, 100)
	n, err := unix.EpollWait(e.fd, events, 100)
	if err != nil {
		return nil, err
	}
	e.lock.Lock()
	defer e.lock.Unlock()
	var conns []*websocket.Conn
	for i := 0; i < n; i++ {
		conn := e.connMap[int(events[i].Fd)]
		conns = append(conns, conn)
	}
	return conns, nil
}

func WsHandle(rw http.ResponseWriter, rq *http.Request) {
	conn, err := upgrade.Upgrade(rw, rq, nil)
	if err != nil {
		fmt.Printf("处理器upgrade错误: %v\n", err)
		return
	}
	if err = epoll.Register(conn); err != nil {
		fmt.Printf("注册Epoll失败: %v\n", err)
		conn.Close()
	}

}

func Start() {
	for {
		cons, err := epoll.Wait()
		if err != nil {
			fmt.Printf("Epoll-wait 错误: %v\n", err)
			continue
		}
		for _, conn := range cons {
			if conn == nil {
				break
			}
			_, msg, err := conn.ReadMessage()
			if err != nil {
				if err := epoll.UnRegister(conn); err != nil {
					fmt.Printf("Epoll-取消注册 错误: %v\n", err)
				}
			} else {
				fmt.Printf("读取客户端消息: %s\n", msg)
			}
		}
	}
}
