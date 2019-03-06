// date: 2019-03-05
package mysql

import (
	"fmt"
	"reflect"
)

//新增
func Create(table interface{}, data interface{}) {
	res := validTable(table)
	if res {
		if err := Db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(m).Error; err != nil {
			fmt.Printf("创建表结构失败: %v\n", err)
		}
		if err := Db.Create(data).Error; err != nil {
			fmt.Printf("写入数据失败: %v\n", err)
		}
	}
}

//更新
func Update(table interface{}, data interface{}) {
	tableName := reflect.TypeOf(table).Name()
	if !Db.HasTable(tableName) {
		fmt.Printf("%s 不存在\n", tableName)
	}
}

//检查表结构是否存在
func validTable(table interface{}) bool {
	tableName := reflect.TypeOf(table).Name()
	if !Db.HasTable(tableName) {
		fmt.Printf("%s 表不存在\n", tableName)
		return false
	}
	return true
}
