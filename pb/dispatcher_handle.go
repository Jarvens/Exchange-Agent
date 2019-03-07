// date: 2019-03-07
package pb

import (
	"fmt"
	"github.com/Jarvens/Exchange-Agent/common"
	"strings"
	"time"
)

//请求分发处理器
func Dispatcher(stream RpcBidStream1_QuoteBidStreamServer, req *RpcRequest1, address string) error {
	event := req.Event
	channel := req.Channel
	switch event {
	case common.Ping:
		if err := stream.Send(&RpcResponse1{Code: 0, Message: "pong", Timestamp: time.Now().Unix()}); err != nil {
			return err
		}
	case common.Subscribe:
		chans := strings.Split(channel, ".")
		switch chans[1] {
		case common.Tick:
			//channel: quote.tick.btc_usdt
			//判断重复订阅
			tickMap := common.SocketMap[address]
			if tickMap == nil {
				tickMap := make(map[string][]string)
				tickMap[common.Tick] = []string{channel}
				common.SocketMap[address] = tickMap
				err := stream.Send(&RpcResponse1{Message: "订阅成功", Code: common.Success, Channel: channel})
				if err != nil {
					return stream.Context().Err()
				}
				fmt.Printf("订阅成功: %v\n", common.SocketMap)
				break
			} else {
				tickChan := tickMap[common.Tick]
				exist, _ := common.Contain(channel, tickChan)
				if exist {
					err := stream.Send(&RpcResponse1{Message: "重复订阅", Code: common.Fail, Channel: channel})
					if err != nil {
						return stream.Context().Err()
					}
					fmt.Printf("重复订阅: %s\n", channel)
				} else {

					tickMap[address] = append(tickChan, channel)
					err := stream.Send(&RpcResponse1{Message: "订阅成功", Code: common.Success, Channel: channel})
					if err != nil {
						return stream.Context().Err()
					}
					fmt.Printf("订阅成功: %v\n", common.SocketMap)
				}
			}
		case common.Depth:
		case common.Kline:
		default:

		}
		//common.SocketMap[address]=
	case common.UbSubscribe:
	}
	return nil
}
