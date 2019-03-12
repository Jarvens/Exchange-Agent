// date: 2019-03-05
package db

import (
	"github.com/Jarvens/Exchange-Agent/model"
	"testing"
	"time"
)

func BenchmarkCreate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		like := &model.Like{Ip: "127.0.0.1", Ua: "kjsdywrjn", Title: "点赞", Hash: 04732545456230, CreatedAt: time.Now()}
		Create(model.Like{}, like)
	}
}
