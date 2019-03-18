// date: 2019-03-15
package quote

type Quote interface {

	//注册
	Register()

	//取消注册
	DeRegister()

	//心跳
	HeartBeat()
}
