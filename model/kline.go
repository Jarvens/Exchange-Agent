// date: 2019-03-05
package model

type Kline struct {
	Id     uint    `json:"id"`
	Symbol string  `json:"symbol"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Volume float64 `json:"volume"`
	Ts     uint    `json:"ts"`
}
