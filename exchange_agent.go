// date: 2019-03-05
package main

import (
	"fmt"
	"github.com/Jarvens/Exchange-Agent/grpc"
	"github.com/Jarvens/Exchange-Agent/tcp"
	"net/http"
)

func main() {
	inChan := make(chan byte)

	http.HandleFunc("/", tcp.WebsocketHandler)
	go http.ListenAndServe("0.0.0.0:12345", nil)
	fmt.Println("[Exchange-Agent]-websocket  启动成功")

	go grpc.QuoteServerStart()
	fmt.Println("[Exchange-Agent]-gRPC 启动成功")

	fmt.Println(<-inChan)
}
