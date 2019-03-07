// date: 2019-03-07
package pb

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"io"
	"net"
	"strings"
)

const (
	serverPort = ":3000"
)

type quoteServer struct{}

func (q *quoteServer) QuoteBidStream(stream RpcBidStream1_QuoteBidStreamServer) error {
	ctx := stream.Context()
	address, err := getClientIp(ctx)
	if err != nil {
		fmt.Printf("获取客户端IP错误: %v\n", err)
	}
	fmt.Printf("客户端IP地址: %s\n", address)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("收到客户端主动断开请求")
			return ctx.Err()
		default:
			request, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("关闭客户端")
				return nil
			}
			if err != nil {
				fmt.Println("读取客户端数据流出错")
			}
			fmt.Printf("读取到数据流: %v\n", request)

		}
	}
	return nil
}

//回写数据
func sendData(stream RpcBidStream1_QuoteBidStreamServer, message, channel string, code int32) error {
	err := stream.Send(&RpcResponse1{Code: code, Message: message, Channel: channel})
	if err != nil {
		fmt.Printf("回写数据出错: %v\n", err)
		return err
	}
	return nil
}

//启动服务
func QuoteServerStart() {
	server := grpc.NewServer()
	RegisterRpcBidStream1Server(server, &quoteServer{})
	address, err := net.Listen("tcp", serverPort)
	if err != nil {
		panic(err)
	}
	if err := server.Serve(address); err != nil {
		panic(err)
	}
}

//获取客户端IP信息
func getClientIp(ctx context.Context) (string, error) {
	pr, ok := peer.FromContext(ctx)
	if !ok {
		return "", errors.New("getClientIp from ctx fail")
	}
	if pr.Addr == net.Addr(nil) {
		return "", errors.New("getClientIp  peer.Address is nil")
	}
	addSlice := strings.Split(pr.Addr.String(), ".")
	return addSlice[0], nil

}
