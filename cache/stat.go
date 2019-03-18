// date: 2019-03-15
package cache

type Stat struct {
	//缓存总量
	Count int64

	//当前缓存 key大小
	KeySize int64

	//当前缓存value大小
	ValueSize int64
}

//写入缓存时 同时要修缓存状态信息
func (s *Stat) add(k string, v []byte) {
	s.Count += 1
	s.KeySize += int64(len(k))
	s.ValueSize += int64(len(v))
}

//删除缓存时同时需要修改缓存状态信息
func (s *Stat) del(k string, v []byte) {
	s.Count -= 1
	s.KeySize -= int64(len(k))
	s.ValueSize -= int64(len(v))
}
