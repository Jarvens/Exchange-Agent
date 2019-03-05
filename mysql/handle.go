// date: 2019-03-05
package mysql

import (
	"fmt"
	"reflect"
)

func Create(m interface{}, data interface{}) {
	fmt.Println("打印结构体名称", reflect.TypeOf(m).Name())
	if !Db.HasTable(reflect.TypeOf(m).Name()) {
		if err := Db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(m).Error; err != nil {
			fmt.Printf("创建表结构失败: %v\n", err)
		}
		if err := Db.Create(data).Error; err != nil {
			fmt.Printf("写入数据失败: %v\n", err)
		}

	}
}
