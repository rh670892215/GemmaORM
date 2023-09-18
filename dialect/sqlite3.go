package dialect

import (
	"reflect"
	"time"
)

// Sqlite3 sqlite3数据库
type Sqlite3 struct{}

func init() {
	RegisterDialect("sqlite3", &Sqlite3{})
}

// DataType 实现DataType接口，golang类型映射成sqlite3中的类型
func (s *Sqlite3) DataType(value reflect.Value) string {
	switch value.Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		return "integer"
	case reflect.Int64, reflect.Uint64:
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "real"
	case reflect.String:
		return "text"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Struct:
		if _, ok := value.Interface().(time.Time); ok {
			return "datetime"
		}
	default:
	}
	return ""
}

// TableExist 输出判断表是否存在的sql
func (s *Sqlite3) TableExist(tableName string) (string, []interface{}) {
	vars := []interface{}{tableName}
	return "SELECT name FROM sqlite_master WHERE type='table' and name = ?", vars
}
