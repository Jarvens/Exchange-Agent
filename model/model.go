// date: 2019-03-05
package model

import "time"

type Like struct {
	Id        int    `gorm:"primary_key"`
	Ip        string `gorm:"type:varchar(20);not null;index:ip_idx"`
	Ua        string `gorm:"type:varchar(256);not null;"`
	Title     string `gorm:"type:varchar(128);not null;index:title_idx"`
	Hash      uint64 `gorm:unique_index:hash_idx;`
	CreatedAt time.Time
}

type Test struct {
}
