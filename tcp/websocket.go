// date: 2019-03-06
package tcp

import (
	"encoding/json"
	"errors"
	"github.com/Jarvens/Exchange-Agent/common"
	"github.com/Jarvens/Exchange-Agent/util/log"
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

var Wmap = make(map[string]*Connection)

type event struct {
	e map[string][]string
}

type Connection struct {
	wsConn    *websocket.Conn
	inChan    chan []byte
	outChan   chan []byte
	closeChan chan byte
	sync.Mutex
	isClosed bool
	*event
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {

	var (
		wsCon *websocket.Conn
		err   error
		conn  *Connection
	)
	if wsCon, err = upgrade.Upgrade(w, r, nil); err != nil {
		log.Errorf("处理器upgrade失败: %v\n", err)
		return
	}
	wsCon.SetCloseHandler(func(code int, text string) error {
		message := websocket.FormatCloseMessage(code, "")
		wsCon.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
		return nil
	})

	if conn, err = NewWebsocket(wsCon); err != nil {
		log.Infof("处理器初始化失败: %v\n", err)
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
			goto ERR
		}
		request := Request{}
		json.Unmarshal(data, &request)
		dispatcher(c, &request)
		select {
		case c.inChan <- data:
		case <-c.closeChan:
			log.Infof("关闭消息接收")
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
	conn := &Connection{wsConn: ws, inChan: make(chan []byte, 1024), outChan: make(chan []byte, 1024), closeChan: make(chan byte, 1)}
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
	case <-c.closeChan:
		log.Infof("[Exchange-Agent]连接关闭: 客户端地址: %s\n", c.wsConn.RemoteAddr().String())
		return nil, errors.New("连接关闭")
	}
	return
}

func (c *Connection) Write(data []byte) (err error) {
	if err = c.wsConn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Errorf("回写消息失败: %v\n", err)
		return err
	}
	log.Infof("回写内容: %s\n", string(data))
	return
}

func (c *Connection) onConnected() {
	log.Infof("[Exchange-Agent]%s 加入会话\n", c.wsConn.RemoteAddr().String())
}

//请求分发
func dispatcher(c *Connection, req *Request) {
	event := req.Event
	channel := req.Channel
	if event == "" || channel == "" {
		log.Errorf("参数错误: event: %s channel: %s", event, channel)
		c.wsConn.WriteJSON(Fail(channel))
		return
	}
	switch event {
	case common.Ping:
		c.wsConn.WriteJSON(Pong(channel))
	case common.Subscribe:
		channelStr := strings.Split(channel, ".")
		if len(channelStr) <= 1 {
			log.Infof("数据错误: %s", channel)
			response := Fail(channel)
			data, _ := json.Marshal(response)
			c.Write(data)
			break
		}
		switch channelStr[1] {
		case common.Tick:
			subTick(c, channel)
		case common.Depth:
			subDepth(c, channel)
		case common.Kline:
			subKline(c, channel)

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
		log.Errorf("JSON转换错误: %v\n", err)
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
		event := &event{make(map[string][]string)}
		event.e[common.Tick] = []string{channel}
		c.event = event
		Wmap[address] = c
	} else {
		tickMap := con.event.e[common.Tick]
		if tickMap == nil {
			event := &event{make(map[string][]string)}
			event.e[common.Tick] = []string{channel}
			con.event = event
			Wmap[address] = con
		} else {
			exist, _ := common.Contain(channel, tickMap)
			if exist {
				response := SubRepeat(channel)
				err := c.wsConn.WriteJSON(response)
				if err != nil {
					log.Errorf("ws订阅失败: %v\n", err)
					return err
				}
				log.Infof("ws订阅成功: 客户端地址: %s 订阅指令: %s\n", c.wsConn.RemoteAddr().String(), channel)
			} else {
				tickMap = append(tickMap, channel)
				c.event.e[common.Tick] = tickMap
				Wmap[address] = c
				response := Success(channel)
				err := c.wsConn.WriteJSON(response)
				if err != nil {
					log.Errorf("ws订阅失败: %v\n", err)
					return err
				}
				log.Infof("ws订阅成功: 客户端地址: %s 订阅指令: %s\n", c.wsConn.RemoteAddr().String(), channel)
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
		channelSlice, ok := depthMap.event.e[common.Depth]
		if ok {
			exist, _ := common.Contain(channel, channelSlice)
			if exist {
				response := SubRepeat(channel)
				err := c.wsConn.WriteJSON(response)
				if err != nil {
					log.Errorf("ws深度订阅重复消息发送失败: %v\n", err)
				}
				log.Infof("ws深度订阅成功: 客户端地址: %s 订阅指令: %s\n", address, channel)
			} else {
				channelSlice := append(channelSlice, channel)
				depthMap.event.e[common.Depth] = channelSlice
				Wmap[address] = depthMap
				response := Success(channel)
				err := c.wsConn.WriteJSON(response)
				if err != nil {
					log.Errorf("深度订阅消息发送失败: %v\n", err)
				}
				log.Infof("ws深度订阅成功: 客户端地址:%s 订阅指令: %s\n", address, channel)
			}
		}
	} else {
		event := &event{make(map[string][]string)}
		event.e[common.Depth] = []string{channel}
		c.event = event
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
		klineChannel, ok := klineMap.event.e[common.Kline]
		if ok {
			response := SubRepeat(channel)
			err := c.wsConn.WriteJSON(response)
			if err != nil {
				log.Errorf("K线重复订阅消息发送失败: %v\n", err)
				return err
			}
			log.Infof("K线重复订阅: 客户端地址: %s 订阅指令: %s\n", address, channel)
			return nil
		} else {
			klineChannel = append(klineChannel, channel)
			klineMap.event.e[common.Kline] = klineChannel
			Wmap[address] = klineMap
		}
	} else {
		event := &event{make(map[string][]string)}
		event.e[common.Kline] = []string{channel}
		c.event = event
		Wmap[address] = c
		response := Success(channel)
		err := c.wsConn.WriteJSON(response)
		if err != nil {
			log.Errorf("K线订阅成功消息发送失败: %v\n", err)
			return err
		}
		log.Infof("ws K线订阅成功: 客户端地址: %s 订阅指令: %s\n", address, channel)
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
