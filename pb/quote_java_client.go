// date: 2019-03-07
package pb

import (
	"bufio"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"os"
)

func QuoteBidStreamClient() {
	conn, err := grpc.Dial("localhost:3000", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("连接错误: %v\n", err)
		return
	}
	defer conn.Close()
	client := NewRpcBidStream1Client(conn)
	ctx := context.Background()
	stream, err := client.QuoteBidStream(ctx)
	if err != nil {
		fmt.Printf("创建数据流失败: %v\n", err)
	}
	go func() {
		fmt.Printf("请输入消息 ...\n")
		input := bufio.NewReader(os.Stdin)
		for {
			line, _ := input.ReadString('\n')
			fmt.Printf("命令行输入: %v\n", line)
			if err := stream.Send(&RpcRequest1{Event: "quote.kline.1m.btc_usdt", Channel: "subscribe"}); err != nil {
				return
			}
		}
	}()

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("收到客户端断开信号")
			break
		}
		if err != nil {
			fmt.Printf("客户端接收数据出错: %v\n", err)
		}
		fmt.Printf("打印数据: message: %s code: %d", res.Message, res.Code)
	}
}
