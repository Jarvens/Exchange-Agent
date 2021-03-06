// date: 2019-03-07
package grpc

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/Jarvens/Exchange-Agent/util/log"
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
		log.Infof("获取客户端IP错误: %v\n", err)
	}
	log.Infof("客户端IP地址Hash值: %x\n", sha256.Sum256([]byte(address)))
	for {
		select {
		case <-ctx.Done():
			log.Infof("收到客户端主动断开请求")
			return ctx.Err()
		default:
			request, err := stream.Recv()
			if err == io.EOF {
				log.Infof("关闭客户端")
				return nil
			}

			if err != nil {
				log.Infof("读取客户端数据流出错: %v\n", err)
				return nil
			}
			err = Dispatcher(stream, request, address)
			if err != nil {
				fmt.Println(err)
				return ctx.Err()
			}
		}
	}
	return nil
}

//回写数据
func sendData(stream RpcBidStream1_QuoteBidStreamServer, message, channel string, code int32) error {
	err := stream.Send(&RpcResponse1{Code: code, Message: message, Channel: channel})
	if err != nil {
		log.Infof("回写数据出错: %v\n", err)
		return err
	}
	return nil
}

//启动服务
func QuoteServerStart() {
	server := grpc.NewServer()
	RegisterRpcBidStream1Server(server, &quoteServer{})
	address, err := net.Listen("websocket", serverPort)
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
		return "", errors.New("获取客户端IP失败")
	}
	if pr.Addr == net.Addr(nil) {
		return "", errors.New("获取客户端IP失败")
	}
	addSlice := strings.Split(pr.Addr.String(), ".")
	return addSlice[0], nil

}
