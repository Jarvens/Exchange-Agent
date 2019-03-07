// date: 2019-03-06
package common

import (
	"errors"
	"reflect"
)

// 数组 map 包含判断
func Contain(source interface{}, target interface{}) (bool, error) {
	tVal := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < tVal.Len(); i++ {
			if tVal.Index(i).Interface() == source {
				return true, nil
			}
		}
	case reflect.Map:
		if tVal.MapIndex(reflect.ValueOf(source)).IsValid() {
			return true, nil
		}

	}
	return false, errors.New("不存在")
}

func SliceRemove(source []string, value string) []string {
	target := source[:0]
	for _, item := range source {
		if item != value {
			target = append(target, item)
		}
	}
	return target
}
