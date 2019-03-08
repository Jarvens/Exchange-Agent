// date: 2019-03-07
package pb

import (
	"fmt"
	"github.com/Jarvens/Exchange-Agent/common"
	"strings"
	"time"
)

// 请求分发处理器，根据指令不同调用具体处理器，处理后将处理结果通过 RpcResponse1返回到客户端
// 订阅/取消订阅，执行同时需要将 GsMap全局对象加锁，防止脏读数据出现
// 根据指令分类为 订阅(subscribe) 取消订阅(un_subscribe)两个大模块 指令分为: 心跳(ping)
// 最新成交(tick) 深度(depth) K线(kline)
// @param  stream gRPC 流对象
// @param  req 请求封装体
// @param  address 客户端地址
// @return error
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

// 订阅最新成交逻辑处理器，根据address key获取GsMap中是否存在该地址的订阅信息
// 不存在则直接将订阅的channel放入tick模块中
// @param  address 客户端地址
// @param  channel 订阅频道
// @param  stream  gRPC流处理器
// @return error   错误信息
func subTick(address, channel string, stream RpcBidStream1_QuoteBidStreamServer) error {
	//给全局Map添加锁，防止脏读
	common.GsMap.Lock.Lock()
	defer common.GsMap.Lock.Unlock()
	tickMap := common.GsMap.ConnMap[address]
	if tickMap == nil {
		tickMap := make(map[string][]string)
		tickMap[common.Tick] = []string{channel}
		common.GsMap.ConnMap[address] = tickMap

		response := &RpcResponse1{Message: "订阅成功", Code: common.Success, Channel: channel, Timestamp: time.Now().Unix()}
		err := sendMessage(stream, response)
		fmt.Printf("订阅成功: %v\n", common.GsMap.ConnMap)
		return err
	} else {
		tickChan := tickMap[common.Tick]
		exist, _ := common.Contain(channel, tickChan)
		if exist {
			response := &RpcResponse1{Message: "重复订阅", Code: common.Fail, Channel: channel, Timestamp: time.Now().Unix()}
			return sendMessage(stream, response)
		} else {
			newChan := append(tickChan, channel)
			tickMap[common.Tick] = newChan
			common.GsMap.ConnMap[address] = tickMap
			response := &RpcResponse1{Message: "订阅成功", Code: common.Success, Channel: channel, Timestamp: time.Now().Unix()}
			err := sendMessage(stream, response)
			fmt.Printf("订阅成功: %v\n", common.GsMap.ConnMap)
			return err
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
		response := &RpcResponse1{Message: "深度订阅成功", Code: common.Success, Channel: channel, Timestamp: time.Now().Unix()}
		err := sendMessage(stream, response)
		fmt.Printf("深度订阅成功: 客户端地址: %s 订阅指令: %s\n", address, channel)
		return err
	} else {
		depthChan := depthMap[common.Depth]
		exist, _ := common.Contain(channel, depthChan)
		if exist {
			response := &RpcResponse1{Message: "深度订阅重复", Code: common.Fail, Channel: channel, Timestamp: time.Now().Unix()}
			err := sendMessage(stream, response)
			fmt.Printf("深度订阅重复: 客户端地址: %s 订阅指令: %s\n", address, channel)
			return err
		} else {
			newChan := append(depthChan, channel)
			depthMap[common.Depth] = newChan
			common.GsMap.ConnMap[address] = depthMap
			response := &RpcResponse1{Message: "深度订阅成功", Code: common.Success, Channel: channel, Timestamp: time.Now().Unix()}
			err := sendMessage(stream, response)
			fmt.Printf("深度订阅成功: 客户端地址: %s 订阅指令: %s\n", address, channel)
			return err
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
		response := &RpcResponse1{Message: "k线订阅成功", Code: common.Success, Channel: channel, Timestamp: time.Now().Unix()}
		err := sendMessage(stream, response)
		fmt.Printf("K线订阅成功:  客户端地址:%s K线订阅指令: %s\n", address, channel)
		return err
	} else {
		klineChan := klineMap[common.Kline]
		exist, _ := common.Contain(channel, klineChan)
		if exist {
			response := &RpcResponse1{Message: "K线重复订阅", Code: common.Fail, Channel: channel, Timestamp: time.Now().Unix()}
			err := sendMessage(stream, response)
			fmt.Printf("K线重复订阅: 客户端地址: %s 订阅指令: %s\n", address, channel)
			return err
		} else {
			newChan := append(klineChan, channel)
			klineMap[common.Kline] = newChan
			common.GsMap.ConnMap[address] = klineMap
			response := &RpcResponse1{Message: "K线订阅成功", Code: common.Success, Channel: channel, Timestamp: time.Now().Unix()}
			err := sendMessage(stream, response)
			fmt.Printf("K线订阅成功: 客户端地址: %s  订阅指令: %s\n", address, channel)
			return err
		}
	}

	return nil
}

//取消成交订阅
func unSubTick(address, channel string, stream RpcBidStream1_QuoteBidStreamServer) error {
	common.GsMap.Lock.Lock()
	defer common.GsMap.Lock.Unlock()

	tickMap := common.GsMap.ConnMap[address]
	if tickMap == nil {
		response := &RpcResponse1{Message: "数据不存在", Code: common.Fail, Channel: channel, Timestamp: time.Now().Unix()}
		err := sendMessage(stream, response)
		fmt.Printf("数据不存在: 客户端地址: %s  取消指令: %s\n", address, channel)
		return err
	} else {
		tickChan := tickMap[common.Tick]
		if len(tickChan) > 0 {
			common.SliceRemove(tickChan, channel)
		}
		response := &RpcResponse1{Message: "取消订阅成功", Code: common.Success, Channel: channel, Timestamp: time.Now().Unix()}
		err := sendMessage(stream, response)
		fmt.Printf("取消深度订阅成功: 客户端地址: %s 取消指令: %s\n", address, channel)
		return err
	}
	return nil

}

//取消深度订阅
func unSubDepth(address, channel string, stream RpcBidStream1_QuoteBidStreamServer) error {
	common.GsMap.Lock.Lock()
	defer common.GsMap.Lock.Unlock()

	depthMap := common.GsMap.ConnMap[address]
	if depthMap == nil {
		response := &RpcResponse1{Message: "数据不存在", Code: common.Fail, Channel: channel, Timestamp: time.Now().Unix()}
		err := sendMessage(stream, response)
		fmt.Printf("数据不存在: 客户端地址: %s  取消指令: %s\n", address, channel)
		return err
	} else {
		depthChan := depthMap[common.Depth]
		if len(depthChan) > 0 {
			common.SliceRemove(depthChan, channel)
		}
		response := &RpcResponse1{Message: "取消订阅成功", Code: common.Success, Channel: channel, Timestamp: time.Now().Unix()}
		err := sendMessage(stream, response)
		fmt.Printf("取消深度订阅成功: 客户端地址: %s 取消指令: %s\n", address, channel)
		return err
	}
	return nil
}

//取消K线订阅
func unSubKline(address, channel string, stream RpcBidStream1_QuoteBidStreamServer) error {
	common.GsMap.Lock.Lock()
	defer common.GsMap.Lock.Unlock()

	klineMap := common.GsMap.ConnMap[address]
	if klineMap == nil {
		response := &RpcResponse1{Message: "数据不存在", Code: common.Fail, Channel: channel, Timestamp: time.Now().Unix()}
		err := sendMessage(stream, response)
		fmt.Printf("数据不存在: 客户端地址: %s  取消指令: %s\n", address, channel)
		return err
	} else {
		klineChan := klineMap[common.Kline]
		if len(klineChan) > 0 {
			common.SliceRemove(klineChan, channel)
		}
		response := &RpcResponse1{Message: "取消订阅成功", Code: common.Success, Channel: channel, Timestamp: time.Now().Unix()}
		err := sendMessage(stream, response)
		fmt.Printf("K线取消订阅成功: 客户端地址: %s 取消指令: %s\n", address, channel)
		return err
	}
	return nil
}

//发送消息
func sendMessage(stream RpcBidStream1_QuoteBidStreamServer, response *RpcResponse1) error {
	err := stream.Send(response)
	if err != nil {
		fmt.Printf("消息发送失败: %v 指令: %s\n", err, response.Channel)
		return err
	}
	return nil
}
