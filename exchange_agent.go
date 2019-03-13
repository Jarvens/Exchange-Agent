// date: 2019-03-05
package main

import (
	"fmt"
	"github.com/Jarvens/Exchange-Agent/grpc"
	"github.com/Jarvens/Exchange-Agent/server"
	"github.com/Jarvens/Exchange-Agent/tcp"
	_ "github.com/Jarvens/Exchange-Agent/util/log"
	"net/http"
)

func main() {
	inChan := make(chan byte)

	http.HandleFunc("/", tcp.WebsocketHandler)
	go http.ListenAndServe("0.0.0.0:12345", nil)
	go grpc.QuoteServerStart()
	server.Run()

	fmt.Println(<-inChan)
}
