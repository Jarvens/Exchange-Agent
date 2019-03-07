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
		chanC := strings.Split(channel, ".")
		switch chanC[1] {
		case common.Tick:
			//channel: quote.tick.btc_usdt
			return subTick(address, channel, stream)
		case common.Depth:
			//channel: quote.depth.btc_usdt
			return subDepth(address, channel, stream)
		case common.Kline:
			//channel: quote.kline.btc_usdt
			return subKline(address, channel, stream)
		default:
			//指令不存在

		}
	case common.UbSubscribe:
		chanC := strings.Split(channel, ".")
		switch chanC[1] {
		case common.Tick:
			return unSubTick(address, channel, stream)
		case common.Depth:
			return unSubDepth(address, channel, stream)
		case common.Kline:
			return unSubKline(address, channel, stream)
		default:
			//指令不存在

		}

	}
	return nil
}

//成交
func subTick(address, channel string, stream RpcBidStream1_QuoteBidStreamServer) error {
	//给全局Map添加锁，防止脏读
	common.GsMap.Lock.Lock()
	defer common.GsMap.Lock.Unlock()
	tickMap := common.GsMap.ConnMap[address]
	if tickMap == nil {
		tickMap := make(map[string][]string)
		tickMap[common.Tick] = []string{channel}
		common.GsMap.ConnMap[address] = tickMap
		err := stream.Send(&RpcResponse1{Message: "订阅成功", Code: common.Success, Channel: channel, Timestamp: time.Now().Unix()})
		if err != nil {
			return stream.Context().Err()
		}
		fmt.Printf("订阅成功: %v\n", common.GsMap.ConnMap)
	} else {
		tickChan := tickMap[common.Tick]
		exist, _ := common.Contain(channel, tickChan)
		if exist {
			err := stream.Send(&RpcResponse1{Message: "重复订阅", Code: common.Fail, Channel: channel, Timestamp: time.Now().Unix()})
			if err != nil {
				return stream.Context().Err()
			}
		} else {
			newChan := append(tickChan, channel)
			tickMap[common.Tick] = newChan
			common.GsMap.ConnMap[address] = tickMap
			err := stream.Send(&RpcResponse1{Message: "订阅成功", Code: common.Success, Channel: channel, Timestamp: time.Now().Unix()})
			if err != nil {
				return stream.Context().Err()
			}
			fmt.Printf("订阅成功: %v\n", common.GsMap.ConnMap)
		}
	}
	return nil
}

//深度
func subDepth(address, channel string, stream RpcBidStream1_QuoteBidStreamServer) error {
	//给全局Map添加锁，防止脏读
	common.GsMap.Lock.Lock()
	defer common.GsMap.Lock.Unlock()
	depthMap := common.GsMap.ConnMap[address]
	if depthMap == nil {
		depthMap = make(map[string][]string)
		depthMap[common.Depth] = []string{channel}
		common.GsMap.ConnMap[address] = depthMap
		err := stream.Send(&RpcResponse1{Message: "深度订阅成功", Code: common.Success, Channel: channel, Timestamp: time.Now().Unix()})
		if err != nil {
			fmt.Printf("深度订阅消息发送失败: %v\n", err)
			return err
		}
		fmt.Printf("深度订阅成功: 客户端地址: %s 订阅指令: %s\n", address, channel)
	} else {
		depthChan := depthMap[common.Depth]
		exist, _ := common.Contain(channel, depthChan)
		if exist {
			err := stream.Send(&RpcResponse1{Message: "深度订阅重复", Code: common.Fail, Channel: channel, Timestamp: time.Now().Unix()})
			if err != nil {
				fmt.Printf("深度订阅重复消息发送失败: %v\n", err)
				return err
			}
			fmt.Printf("深度订阅重复: 客户端地址: %s 订阅指令: %s\n", address, channel)
		} else {
			newChan := append(depthChan, channel)
			depthMap[common.Depth] = newChan
			common.GsMap.ConnMap[address] = depthMap
			err := stream.Send(&RpcResponse1{Message: "深度订阅成功", Code: common.Success, Channel: channel, Timestamp: time.Now().Unix()})
			if err != nil {
				fmt.Printf("深度订阅成功消息发送失败: %v\n", err)
				return err
			}
			fmt.Printf("深度订阅成功: 客户端地址: %s 订阅指令: %s\n", address, channel)
		}
	}
	return nil
}

//k线
func subKline(address, channel string, stream RpcBidStream1_QuoteBidStreamServer) error {
	//给全局Map添加锁，防止脏读
	common.GsMap.Lock.Lock()
	defer common.GsMap.Lock.Unlock()
	klineMap := common.GsMap.ConnMap[address]

	if klineMap == nil {
		klineMap = make(map[string][]string)
		klineMap[common.Kline] = []string{channel}
		common.GsMap.ConnMap[address] = klineMap
		err := stream.Send(&RpcResponse1{Message: "k线订阅成功", Code: common.Success, Channel: channel, Timestamp: time.Now().Unix()})
		if err != nil {
			fmt.Printf("K线订阅消息发送失败: %v\n", err)
			return err
		}
		fmt.Printf("K线订阅成功:  客户端地址:%s K线订阅指令: %s\n", address, channel)
	} else {
		klineChan := klineMap[common.Kline]
		exist, _ := common.Contain(channel, klineChan)
		if exist {
			err := stream.Send(&RpcResponse1{Message: "K线重复订阅", Code: common.Fail, Channel: channel, Timestamp: time.Now().Unix()})
			if err != nil {
				fmt.Printf("K线重复订阅消息发送失败: %v\n", err)
				return err
			}
			fmt.Printf("K线重复订阅: 客户端地址: %s 订阅指令: %s\n", address, channel)
		} else {
			newChan := append(klineChan, channel)
			klineMap[common.Kline] = newChan
			common.GsMap.ConnMap[address] = klineMap
			err := stream.Send(&RpcResponse1{Message: "K线订阅成功", Code: common.Success, Channel: channel, Timestamp: time.Now().Unix()})
			if err != nil {
				fmt.Printf("K线订阅成功消息发送失败: %v\n", err)
				return err
			}
			fmt.Printf("K线订阅成功: 客户端地址: %s  订阅指令: %s\n", address, channel)
		}
	}

	return nil
}

//取消成交订阅
func unSubTick(address, channel string, stream RpcBidStream1_QuoteBidStreamServer) error {
	common.GsMap.Lock.Lock()
	defer common.GsMap.Lock.Unlock()

	return nil

}

//取消深度订阅
func unSubDepth(address, channel string, stream RpcBidStream1_QuoteBidStreamServer) error {
	common.GsMap.Lock.Lock()
	defer common.GsMap.Lock.Unlock()

	depthMap := common.GsMap.ConnMap[address]
	if depthMap == nil {

	}
	return nil
}

//取消K线订阅
func unSubKline(address, channel string, stream RpcBidStream1_QuoteBidStreamServer) error {
	common.GsMap.Lock.Lock()
	defer common.GsMap.Lock.Unlock()

	klineMap := common.GsMap.ConnMap[address]
	if klineMap == nil {
		err := stream.Send(&RpcResponse1{Message: "K先订阅数据不存在", Code: common.Fail, Channel: channel, Timestamp: time.Now().Unix()})
		if err != nil {
			fmt.Printf("K线取消订阅消息发送失败: %v\n", err)
			return err
		}
		fmt.Printf("K线取消订阅数据不存在: ")
	} else {
		klineChan := klineMap[common.Kline]
		if len(klineChan) > 0 {
			common.SliceRemove(klineChan, channel)
		}
		err := stream.Send(&RpcResponse1{Message: "K线取消订阅成功", Code: common.Success, Channel: channel, Timestamp: time.Now().Unix()})
		if err != nil {
			fmt.Printf("K线取消订阅消息发送失败: %v\n", err)
		}
		fmt.Printf("K线取消订阅成功: 客户端地址: %s 取消指令: %s\n", address, channel)
	}
	return nil
}
