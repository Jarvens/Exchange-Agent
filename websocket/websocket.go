// date: 2019-03-06
package websocket

import (
	"encoding/json"
	"errors"
	"github.com/Jarvens/Exchange-Agent/common"
	"github.com/Jarvens/Exchange-Agent/util/log"
	"github.com/gorilla/websocket"
	"net/http"
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
	sync.RWMutex
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
		log.Errorf("处理器upgrade失败: %v", err)
		return
	}
	wsCon.SetCloseHandler(func(code int, text string) error {
		message := websocket.FormatCloseMessage(code, "")
		wsCon.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
		return nil
	})

	if conn, err = NewWebsocket(wsCon); err != nil {
		log.Infof("处理器初始化失败: %v", err)
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
		log.Infof("[Agent]连接关闭: 客户端地址: %s", c.wsConn.RemoteAddr().String())
		return nil, errors.New("连接关闭")
	}
	return
}

func (c *Connection) Write(data []byte) (err error) {
	if err = c.wsConn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Errorf("回写消息失败: %v", err)
		return err
	}
	log.Infof("回写内容: %s", string(data))
	return
}

func (c *Connection) onConnected() {
	log.Infof("[Agent]建立连接: 客户端地址: %s", c.wsConn.RemoteAddr().String())
}

//请求分发
func dispatcher(c *Connection, req *Request) {
	event, eventExist := req.checkEvent()
	if eventExist {
		switch event {
		case common.PING:
			c.wsConn.WriteJSON(Pong())
		case common.SUB:
			channel, module, moduleExist := req.checkModule()
			if moduleExist {
				switch module {
				case common.TICK:
					subTick(c, channel)
				case common.DEPTH:
					subDepth(c, channel)
				case common.KLINE:
					subKline(c, channel)
				}
			} else {
				c.wsConn.WriteJSON(ChannelErr())
			}
		case common.UN_SUB:
			//TODO 取消订阅
		default:
			c.wsConn.WriteJSON(EventErr())
		}
	}
	c.wsConn.WriteJSON(EventErr())
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
		event.e[common.TICK] = []string{channel}
		c.event = event
		Wmap[address] = c
	} else {
		tickMap := con.event.e[common.TICK]
		if tickMap == nil {
			event := &event{make(map[string][]string)}
			event.e[common.TICK] = []string{channel}
			con.event = event
			Wmap[address] = con
		} else {
			exist, _ := common.Contain(channel, tickMap)
			if exist {
				response := SubRepeat(channel)
				err := c.wsConn.WriteJSON(response)
				if err != nil {
					log.Errorf("ws订阅失败: %v", err)
					return err
				}
				log.Infof("ws订阅成功: 客户端地址: %s 订阅指令: %s\n", c.wsConn.RemoteAddr().String(), channel)
			} else {
				tickMap = append(tickMap, channel)
				c.event.e[common.TICK] = tickMap
				Wmap[address] = c
				response := Success(channel)
				err := c.wsConn.WriteJSON(response)
				if err != nil {
					log.Errorf("ws订阅失败: %v", err)
					return err
				}
				log.Infof("ws订阅成功: 客户端地址: %s 订阅指令: %s", c.wsConn.RemoteAddr().String(), channel)
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
		channelSlice, ok := depthMap.event.e[common.DEPTH]
		if ok {
			exist, _ := common.Contain(channel, channelSlice)
			if exist {
				response := SubRepeat(channel)
				err := c.wsConn.WriteJSON(response)
				if err != nil {
					log.Errorf("ws深度订阅重复消息发送失败: %v", err)
				}
				log.Infof("ws深度订阅成功: 客户端地址: %s 订阅指令: %s", address, channel)
			} else {
				channelSlice := append(channelSlice, channel)
				depthMap.event.e[common.DEPTH] = channelSlice
				Wmap[address] = depthMap
				response := Success(channel)
				err := c.wsConn.WriteJSON(response)
				if err != nil {
					log.Errorf("深度订阅消息发送失败: %v", err)
				}
				log.Infof("ws深度订阅成功: 客户端地址:%s 订阅指令: %s", address, channel)
			}
		}
	} else {
		event := &event{make(map[string][]string)}
		event.e[common.DEPTH] = []string{channel}
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
		klineChannel, ok := klineMap.event.e[common.KLINE]
		if ok {
			response := SubRepeat(channel)
			err := c.wsConn.WriteJSON(response)
			if err != nil {
				log.Errorf("K线重复订阅消息发送失败: %v", err)
				return err
			}
			log.Infof("K线重复订阅: 客户端地址: %s 订阅指令: %s", address, channel)
			return nil
		} else {
			klineChannel = append(klineChannel, channel)
			klineMap.event.e[common.KLINE] = klineChannel
			Wmap[address] = klineMap
		}
	} else {
		event := &event{make(map[string][]string)}
		event.e[common.KLINE] = []string{channel}
		c.event = event
		Wmap[address] = c
		response := Success(channel)
		err := c.wsConn.WriteJSON(response)
		if err != nil {
			log.Errorf("K线订阅成功消息发送失败: %v", err)
			return err
		}
		log.Infof("ws K线订阅成功: 客户端地址: %s 订阅指令: %s", address, channel)
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
