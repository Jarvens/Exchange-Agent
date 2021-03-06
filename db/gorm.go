// date: 2019-03-05
package db

import (
	"github.com/Jarvens/Exchange-Agent/config"
	"github.com/Jarvens/Exchange-Agent/util/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var ORM, Errs = GormInit()

func GormInit() (*gorm.DB, error) {
	db, error := gorm.Open("mysql", config.MysqlDSL())
	db.DB()
	db.DB().Ping()
	db.DB().SetMaxOpenConns(2000)
	db.DB().SetMaxIdleConns(1000)
	db.SingularTable(true)
	if config.Environment == "DEVELOPMENT" {
		db.LogMode(true)
	}
	log.CheckError(error)
	return db, error
}
