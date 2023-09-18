package dialect

import (
	"reflect"
	"sync"
)

// Dialect 封装底层数据库的差异
type Dialect interface {
	DataType(value reflect.Value) string
	TableExist(tableName string) (string, []interface{})
}

var dialectMap map[string]Dialect
var once sync.Once

// RegisterDialect 注册数据库类型
func RegisterDialect(name string, dialect Dialect) {
	once.Do(func() {
		dialectMap = make(map[string]Dialect)
	})

	dialectMap[name] = dialect
}

// GetDialect 根据数据库类型获取Dialect
func GetDialect(name string) (Dialect, bool) {
	if dialectMap == nil {
		return nil, false
	}

	res, ok := dialectMap[name]
	return res, ok
}
