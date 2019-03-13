// date: 2019-03-05
package main

import (
	"fmt"
	"github.com/Jarvens/Exchange-Agent/config"
	"github.com/Jarvens/Exchange-Agent/grpc"
	"github.com/Jarvens/Exchange-Agent/server"
	"github.com/Jarvens/Exchange-Agent/tcp"
	"github.com/Jarvens/Exchange-Agent/util/log"
	"net/http"
)

func main() {
	inChan := make(chan byte)

	http.HandleFunc("/", tcp.WebsocketHandler)
	go http.ListenAndServe("0.0.0.0:12345", nil)
	log.Debug("[Exchange-Agent]-websocket  启动成功")

	go grpc.QuoteServerStart()
	log.Debug("[Exchange-Agent]-gRPC 启动成功")
	server.Run()

	fmt.Println(<-inChan)
}

func init() {
	log.Init(config.Environment)
}
