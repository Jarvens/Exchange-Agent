// date: 2019-03-06
package tcp

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Jarvens/Exchange-Agent/common"
	"github.com/gorilla/websocket"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"
)

var upgrade = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
}, EnableCompression: true}

// 使用fileDescription 文件描述符的优点在于不受限制
// 使用fileDescription 需要每次进行参数校验判断fd是否为当前连接fd
// 使用地址+端口的方式   受限制
var Wmap = make(map[string]*Connection)

type Connection struct {
	wsConn    *websocket.Conn
	inChan    chan []byte
	outChan   chan []byte
	closeChan chan byte
	sync.Mutex
	isClosed bool
	event    map[string][]string
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
	//fd := makeFd(ws)
	//fmt.Printf("文件描述符: %d\n", fd)
	conn := &Connection{wsConn: ws, inChan: make(chan []byte, 1024), outChan: make(chan []byte, 1024), closeChan: make(chan byte, 1)}
	//连接进来不需要将当前连接加入链接储存器，由后续订阅操作添加
	//Wmap[fd] = conn
	conn.onConnected()
	go conn.LoopRead()
	go conn.LoopWrite()
	return conn, nil
}

func (c *Connection) Close() {
	c.wsConn.Close()
	c.Lock()
	defer c.Unlock()
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

//请求分发
func dispatcher(c *Connection, req *Request) {
	event := req.Event
	channel := req.Channel
	switch event {
	case common.Ping:
		c.wsConn.WriteJSON(Pong(channel))
	case common.Subscribe:
		channelStr := strings.Split(channel, ".")
		switch channelStr[1] {
		case common.Tick:
		case common.Depth:
		case common.Kline:

		}

	case common.UnSubscribe:

	}
}

type Response struct {
	Code      int
	Message   string
	Timestamp int64
	Channel   string
}

type Request struct {
	Event   string
	Channel string
}

func Success(channel string) *Response {
	return &Response{Code: common.Success, Message: "成功", Channel: channel, Timestamp: time.Now().Unix()}
}

func Fail(channel string) *Response {
	return &Response{Code: common.Fail, Message: "失败", Channel: channel, Timestamp: time.Now().Unix()}
}

func SubRepeat(channel string) *Response {
	return &Response{Message: "重复订阅", Code: common.Fail, Channel: channel, Timestamp: time.Now().Unix()}
}

func Pong(channel string) *Response {
	return &Response{Code: common.Success, Message: common.Pong, Channel: channel, Timestamp: time.Now().Unix()}
}

func toJson(data interface{}) []byte {
	bytes, err := json.Marshal(&data)
	if err != nil {
		fmt.Printf("JSON转换错误: %v\n", err)
		return nil
	}
	return bytes
}

//订阅最新成交
func subTick(c *Connection, channel string) error {
	//加锁 防止脏读数据
	c.Lock()
	defer c.Unlock()
	address := c.wsConn.RemoteAddr().String()
	con := Wmap[address]
	if con == nil {
		tickMap := make(map[string][]string)
		tickMap[common.Tick] = []string{channel}
		c.event = tickMap
		Wmap[address] = c
	} else {
		tickMap := con.event[common.Tick]
		if tickMap == nil {
			tickMap := make(map[string][]string)
			tickMap[common.Tick] = []string{channel}
			con.event = tickMap
			Wmap[address] = con
		} else {
			exist, _ := common.Contain(channel, tickMap)
			if exist {
				response := SubRepeat(channel)
				err := c.wsConn.WriteJSON(response)
				if err != nil {
					fmt.Printf("ws订阅失败: %v\n", err)
					return err
				}
				fmt.Printf("ws订阅成功: 客户端地址: %s 订阅指令: %s\n", c.wsConn.RemoteAddr().String(), channel)
			} else {
				tickMap = append(tickMap, channel)
				c.event[common.Tick] = tickMap
				Wmap[address] = c
				response := Success(channel)
				err := c.wsConn.WriteJSON(response)
				if err != nil {
					fmt.Printf("ws订阅失败: %v\n", err)
					return err
				}
				fmt.Printf("ws订阅成功: 客户端地址: %s 订阅指令: %s\n", c.wsConn.RemoteAddr().String(), channel)
			}
		}
	}
	return nil
}

//订阅深度
func subDepth(c *Connection, channel string) {
	c.Lock()
	defer c.Unlock()
	address := c.wsConn.RemoteAddr().String()
	depthMap, ok := Wmap[address]
	if ok {
		channelSlice, ok := depthMap.event[common.Depth]
		if ok {
			exist, _ := common.Contain(channel, channelSlice)
			if exist {
				response := SubRepeat(channel)
				err := c.wsConn.WriteJSON(response)
				if err != nil {
					fmt.Printf("ws深度订阅重复消息发送失败: %v\n", err)
				}
				fmt.Printf("ws深度订阅成功: 客户端地址: %s 订阅指令: %s\n", address, channel)
			} else {
				channelSlice := append(channelSlice, channel)
				depthMap.event[common.Depth] = channelSlice
				Wmap[address] = depthMap
				response := Success(channel)
				err := c.wsConn.WriteJSON(response)
				if err != nil {
					fmt.Printf("深度订阅消息发送失败: %v\n", err)
				}
				fmt.Printf("ws深度订阅成功: 客户端地址:%s 订阅指令: %s\n", address, channel)
			}
		}
	} else {
		depthMap := make(map[string][]string)
		depthMap[common.Depth] = []string{channel}
		c.event = depthMap
		Wmap[address] = c
	}
}

//订阅K线
func subKline(c *Connection, channel string) error {
	c.Lock()
	defer c.Unlock()
	address := c.wsConn.RemoteAddr().String()
	klineMap, ok := Wmap[address]
	if ok {
		klineChannel, ok := klineMap.event[common.Kline]
		if ok {
			response := SubRepeat(channel)
			err := c.wsConn.WriteJSON(response)
			if err != nil {
				fmt.Printf("K线重复订阅消息发送失败: %v\n", err)
				return err
			}
			fmt.Printf("K线重复订阅: 客户端地址: %s 订阅指令: %s\n", address, channel)
			return nil
		} else {
			klineChannel = append(klineChannel, channel)
			klineMap.event[common.Kline] = klineChannel
			Wmap[address] = klineMap
		}
	} else {
		klineMap := make(map[string][]string)
		klineMap[common.Kline] = []string{channel}
		c.event = klineMap
		Wmap[address] = c
		response := Success(channel)
		err := c.wsConn.WriteJSON(response)
		if err != nil {
			fmt.Printf("K线订阅成功消息发送失败: %v\n", err)
			return err
		}
		fmt.Printf("ws K线订阅成功: 客户端地址: %s 订阅指令: %s\n", address, channel)
		return nil
	}
	return nil
}

//取消成交订阅
func unSubTick(c *Connection, channel string) {
	c.Lock()
	defer c.Unlock()
}

//取消深度订阅
func unSubDepth(c *Connection, channel string) {
	c.Lock()
	defer c.Unlock()
}

//取消K线订阅
func unSubKline(c *Connection, channel string) {
	c.Lock()
	defer c.Unlock()

}
