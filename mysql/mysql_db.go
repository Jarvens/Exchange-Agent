// date: 2019-03-05
package mysql

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Db *gorm.DB

//使用mysql初始化数据库
//func init()  {
////	Db,_=sql.Open("mysql","root:@tcp(127.0.0.1:3306)/eoe?charset=utf8&parseTime=true")
////	Db.SetMaxOpenConns(2000)
////	Db.SetMaxIdleConns(1000)
////	Db.Ping()
////}

//使用gorm初始化数据库
//使用gorm框架连接数据库需要导入dialects
//dialects库实际就是 针对不同的数据库import了不同的驱动而已
func init() {
	Db, _ = gorm.Open("mysql", "root:root@/eoe?charset=utf8&parseTime=True&loc=Local")
	Db.DB().SetMaxOpenConns(2000)
	Db.DB().SetMaxIdleConns(1000)
	//关闭复数形式 默认表明为 model 复数
	Db.SingularTable(true)
	Db.LogMode(true)
}
