// date: 2019-03-15
package _interface

type WsServer interface {
	//连接
	Connect()

	//断开连接
	DisConnect()

	//异常
	Exception()
}
