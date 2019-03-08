// date: 2019-03-08
package service

import "github.com/Jarvens/Exchange-Agent/mysql"

type Kline struct {
	Id     int     //主键
	Symbol string  //交易对   BTC_USDT
	Open   float64 //开盘价
	Close  float64 //收盘价
	High   float64 //最高价
	Low    float64 //最低价
	Volume float64 //24小时成交量
	Ts     int64   //成交时间戳
}

// 查询24H成交总量
func Queqy24hTick() float64 {
	mysql.Db.Table("kline1m").Select("sum(volume)").Limit(1440)
	return 0
}
