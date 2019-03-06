// date: 2019-03-05
package main

import (
	"fmt"
	"github.com/Jarvens/Exchange-Agent/tcp"
	"net/http"
)

func main() {
	inchan := make(chan byte)
	http.HandleFunc("/", tcp.WebsocketHandler)
	go http.ListenAndServe("0.0.0.0:12345", nil)
	fmt.Println(<-inchan)
}
