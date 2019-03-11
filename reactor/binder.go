// date: 2019-03-09
package reactor

import (
	"errors"
	"sync"
)

type eventConn struct {
	Event   string
	Conn    *Conn
	Channel []string
}

type binder struct {
	sync.RWMutex

	//uuid-连接地址 map
	uuid2EventMap map[string]*eventConn
}

// 绑定 uuid-address-event
func (b *binder) Bind(uuid, event string, conn *Conn) error {

	if uuid == "" {
		return errors.New("uuid不能为空")
	}
	if event == "" {
		return errors.New("event不能为空")
	}
	if conn == nil {
		return errors.New("连接不能为空")
	}

	b.Lock()
	defer b.Unlock()

	//if eConns,ok:=b.uuid2AddressMap[]
	return nil
}

func (b *binder) FindConn(uuid string) (*Conn, bool) {
	if uuid == "" {
		return nil, false
	}
	address, ok := b.uuid2AddressMap[uuid]
	if ok {
		if events, ok := b.address2EventConnMap[address]; ok {
			for i := range *events {
				if (*events)[i].Conn.GetID() == uuid {
					return (*events)[i].Conn, true
				}
			}
		}
	}
	return nil, false
}
