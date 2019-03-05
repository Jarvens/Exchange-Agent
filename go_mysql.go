// date: 2019-03-05
package main

import (
	"github.com/Jarvens/Exchange-Agent/model"
	"github.com/Jarvens/Exchange-Agent/mysql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func main() {
	//if !mysql.Db.HasTable(&model.Like{}){
	//	if err:=mysql.Db.Set("gorm:table_options","ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&model.Like{}).Error;err!=nil{
	//		panic(err)
	//	}
	//}
	//like:=&model.Like{Ip:"127.0.0.1",Ua:"kjsdywrjn",Title:"点赞",Hash:04732545456230,CreatedAt:time.Now()}
	//if err:=mysql.Db.Create(like).Error;err!=nil{
	//	fmt.Println("写入数据库发生错误")
	//}
	like := &model.Like{Ip: "127.0.0.1", Ua: "kjsdywrjn", Title: "点赞", Hash: 04732545456230, CreatedAt: time.Now()}
	mysql.Create(model.Like{}, like)
	//mysql.Create(model.Test{},&model.Test{})
}
