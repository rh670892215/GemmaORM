package session

import (
	"GemmaORM/log"
	"reflect"
)

const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

// CallMethod 根据method调用session的tableSchema实现的指定的方法，value可以传入实现了method方法的struct，优先使用value结构体实现的方法
func (s *Session) CallMethod(method string, value interface{}) {
	fm := reflect.ValueOf(s.GetRefTable().Model).MethodByName(method)
	if value != nil {
		fm = reflect.ValueOf(value).MethodByName(method)
	}

	if fm.IsValid() {
		params := []reflect.Value{reflect.ValueOf(s)}
		if res := fm.Call(params); len(res) > 0 {
			if err, ok := res[0].Interface().(error); ok {
				log.Error(err)
			}
		}
	}

}
