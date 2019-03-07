// date: 2019-03-05
package model

import "time"

type Like struct {
	Id        int    `gorm:"primary_key"`
	Ip        string `gorm:"type:varchar(20);not null;index:ip_idx"`
	Ua        string `gorm:"type:varchar(256);not null;"`
	Title     string `gorm:"type:varchar(128);not null;index:title_idx"`
	Hash      uint64
	CreatedAt time.Time
}

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
