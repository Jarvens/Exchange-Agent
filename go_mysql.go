// date: 2019-03-05
package main

import (
	"fmt"
	"github.com/Jarvens/Exchange-Agent/model"
	"github.com/Jarvens/Exchange-Agent/mysql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func main() {
	//创建
	like := &model.Like{Ip: "127.0.0.1", Ua: "kjsdywrjn", Title: "点赞", Hash: 04732545456230, CreatedAt: time.Now()}
	mysql.Create(model.Like{}, like)
	//删除
	//like1:=&model.Like{Id:27175}
	//mysql.Db.Delete(&model.Like{},like1)
	//分页查询
	//var count uint64
	//likeList:=make([]*model.Like,0)
	//if err:=mysql.Db.Offset(1).Limit(50).Find(&likeList).Count(&count).Error;err!=nil{
	//	fmt.Printf("分页查询错误: %v\n",err)
	//}
	likes := make([]*model.Like, 0)
	var count uint64

	if err := mysql.Db.Find(&likes).Count(&count).Error; err != nil {
		fmt.Printf("查询结果集失败: %v\n", err)
	}
	fmt.Printf("打印查询结果:%v\n", likes)
	fmt.Printf("打印查询结果:%d\n", count)

	//fmt.Printf("查询结果集: %v count值: %d",likeList,count)
}
