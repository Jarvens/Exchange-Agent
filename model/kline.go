// date: 2019-03-05
package model

type kline struct {
	Id     uint    `json:"id"`
	Symbol string  `json:"symbol"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Volume float64 `json:"volume"`
	Ts     uint    `json:"ts"`
}

//分钟线
type Kline1m struct {
	kline
}

//5分钟K线
type Kline5m struct {
	kline
}

//15分钟K线
type Kline15m struct {
	kline
}

//30分钟K线
type Kline30m struct {
	kline
}

//60分钟K线
type Kline60m struct {
	kline
}

//4小时K线
type Kline240m struct {
	kline
}

//1天K线
type Kline1440m struct {
	kline
}

//5天K线
type Kline7200m struct {
	kline
}

//1周K线
type Kline10080m struct {
	kline
}

//1月K线
type Kline8640m struct {
	kline
}
