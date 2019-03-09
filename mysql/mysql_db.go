// date: 2019-03-05
package mysql

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Db *gorm.DB

func init() {
	Db, _ = gorm.Open("mysql", "root:root@/eoe?charset=utf8&parseTime=True&loc=Local")
	Db.DB().SetMaxOpenConns(2000)
	Db.DB().SetMaxIdleConns(1000)
	Db.SingularTable(true)
	Db.LogMode(true)
}
