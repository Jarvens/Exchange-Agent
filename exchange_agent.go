// date: 2019-03-05
package main

import (
	"fmt"
	"github.com/Jarvens/Exchange-Agent/grpc"
	"github.com/Jarvens/Exchange-Agent/server"
	_ "github.com/Jarvens/Exchange-Agent/util/log"
	"github.com/Jarvens/Exchange-Agent/websocket"
	"net/http"
)

func main() {
	inChan := make(chan byte)

	http.HandleFunc("/", websocket.WebsocketHandler)
	go http.ListenAndServe("0.0.0.0:12345", nil)
	go grpc.QuoteServerStart()
	server.Run()

	fmt.Println(<-inChan)
}
