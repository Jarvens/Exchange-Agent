// date: 2019-03-13
package service

import (
	"github.com/Jarvens/Exchange-Agent/db"
	"github.com/Jarvens/Exchange-Agent/model"
	"github.com/Jarvens/Exchange-Agent/util/log"
)

func FindKline() {

}

//K线入库失败需要将错误K线数据写入处理容器，并且需要提供工具类操作未写入成功的数据
//

//分钟线
func CreateKline1m(kline model.Kline1m) {
	if err := db.ORM.Create(&kline).Error; err != nil {
		log.Infof("分钟线入库失败: %v", err)
	}
}

//5分钟线
func CreateKline5m(kline model.Kline5m) {
	if err := db.ORM.Create(&kline).Error; err != nil {
		log.Infof("5分钟线入库失败: %v", err)
	}
}

//15分钟线
func CreateKline15m(kline model.Kline15m) {
	if err := db.ORM.Create(&kline).Error; err != nil {
		log.Infof("15分钟线入库失败: %v", err)
	}
}

//30分钟线
func CreateKline30m(kline model.Kline30m) {
	if err := db.ORM.Create(&kline).Error; err != nil {
		log.Infof("30分钟线入库失败: %v", err)
	}
}

//小时线分钟线
func CreateKline60m(kline model.Kline60m) {
	if err := db.ORM.Create(&kline).Error; err != nil {
		log.Infof("小时线入库失败: %v", err)
	}
}

//4小时线
func CreateKline240m(kline model.Kline240m) {
	if err := db.ORM.Create(&kline).Error; err != nil {
		log.Infof("4小时线入库失败: %v", err)
	}
}

//1天线
func CreateKline1440m(kline model.Kline1440m) {
	if err := db.ORM.Create(&kline).Error; err != nil {
		log.Infof("1天线入库失败: %v", err)
	}
}

//5天线
func CreateKline7200m(kline model.Kline7200m) {
	if err := db.ORM.Create(&kline).Error; err != nil {
		log.Infof("5天线入库失败: %v", err)
	}
}

//周线
func CreateKline10080m(kline model.Kline10080m) {
	if err := db.ORM.Create(&kline).Error; err != nil {
		log.Infof("周线入库失败: %v", err)
	}
}

//月线
func CreateKline8640m(kline model.Kline8640m) {
	if err := db.ORM.Create(&kline).Error; err != nil {
		log.Infof("月线入库失败: %v", err)
	}
}
